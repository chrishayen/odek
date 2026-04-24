package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	openai "shotgun.dev/odek/openai"
)

// getCurrentTimeTool is a simple tool the model can call to demonstrate
// the tool-calling round-trip.
var getCurrentTimeTool = openai.Tool{
	Type: openai.ToolTypeFunction,
	Function: &openai.FunctionDefinition{
		Name:        "get_current_time",
		Description: "Get the current date and time in a specified timezone",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"timezone": map[string]any{
					"type":        "string",
					"description": "IANA timezone name, e.g. America/New_York, UTC, Asia/Tokyo. Defaults to UTC if omitted.",
				},
			},
		},
	},
}

func main() {
	promptFlag := flag.String("prompt", "", "Prompt to send to the model (reads from stdin if omitted)")
	useTool := flag.Bool("tool", false, "Demonstrate tool calling with a get_current_time tool")
	modelFlag := flag.String("model", openai.DefaultModel, "Model name to use")
	flag.Parse()

	baseURL := os.Getenv("GPT_OSS_URL")
	token := os.Getenv("GPT_OSS_TOKEN")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	client, err := openai.NewClient(baseURL, token)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	prompt := *promptFlag
	if prompt == "" {
		fmt.Fprint(os.Stderr, "Enter prompt: ")
		reader := bufio.NewReader(os.Stdin)
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read prompt: %v", err)
		}
		prompt = strings.TrimSpace(line)
	}
	if prompt == "" {
		log.Fatal("Prompt is required")
	}

	ctx := context.Background()

	if *useTool {
		runToolDemo(ctx, client, *modelFlag, prompt)
	} else {
		runChat(ctx, client, *modelFlag, prompt)
	}
}

func runChat(ctx context.Context, client *openai.Client, model, prompt string) {
	resp, err := client.Chat(ctx, &openai.ChatCompletionRequest{
		Model:    model,
		Messages: []openai.ChatMessage{{Role: openai.RoleUser, Content: prompt}},
	})
	if err != nil {
		log.Fatalf("Chat failed: %v", err)
	}

	if len(resp.Choices) == 0 {
		log.Fatal("No choices in response")
	}

	msg := resp.Choices[0].Message
	if msg.Content != "" {
		fmt.Printf("%s\n", msg.Content)
	}
	for _, tc := range msg.ToolCalls {
		fmt.Printf("[tool call] %s(%s)\n", tc.Function.Name, tc.Function.Arguments)
	}
}

func runToolDemo(ctx context.Context, client *openai.Client, model, prompt string) {
	messages := []openai.ChatMessage{
		{Role: openai.RoleSystem, Content: "You are a helpful assistant. Use the provided tools when appropriate."},
		{Role: openai.RoleUser, Content: prompt},
	}

	tools := []openai.Tool{getCurrentTimeTool}

	handler := func(_ context.Context, call openai.ToolCall) (string, bool, error) {
		switch call.Function.Name {
		case "get_current_time":
			var args struct {
				Timezone string `json:"timezone"`
			}
			if err := json.Unmarshal([]byte(call.Function.Arguments), &args); err != nil {
				return fmt.Sprintf("error parsing arguments: %v", err), false, nil
			}
			tz := args.Timezone
			if tz == "" {
				tz = "UTC"
			}
			loc, err := time.LoadLocation(tz)
			if err != nil {
				return fmt.Sprintf("error loading timezone %q: %v", tz, err), false, nil
			}
			now := time.Now().In(loc)
			result := fmt.Sprintf("The current time in %s is %s.", tz, now.Format(time.RFC1123))
			fmt.Fprintf(os.Stderr, "[tool executed] get_current_time(%s) -> %s\n", tz, now.Format(time.RFC1123))
			return result, false, nil
		default:
			return fmt.Sprintf("unknown tool: %s", call.Function.Name), false, nil
		}
	}

	final, _, err := client.AskToolLoop(ctx, messages, tools, handler, 5, nil)
	if err != nil {
		log.Fatalf("Tool loop failed: %v", err)
	}

	if final.Content != "" {
		fmt.Printf("%s\n", final.Content)
	}
	for _, tc := range final.ToolCalls {
		fmt.Printf("[tool call] %s(%s)\n", tc.Function.Name, tc.Function.Arguments)
	}
}
