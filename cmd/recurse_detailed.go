package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// Configuration for the OpenAI-compatible API
const (
	API_URL    = "http://localhost:1234/v1/chat/completions"
	MODEL_NAME = "gpt-4o"
)

// The Master Architect Prompt
const SYSTEM_COMPOSITION_PROMPT = `You are the Composition Tree Engine (CTE).
Your purpose is to decompose requirements into mathematically traceable, hierarchical composition trees.

You follow a strict "Top-向 (Top-Down) Discovery" protocol with two distinct modes:

1. REGISTRATION MODE (Triggered by a single root name, e.g., "minux"):
   - You act as a registry for the system's high-level namespaces.
   - Output ONLY the top-level namespace list (the "Root Map"). 
   - Format: minux { kernel, userspace, drivers }

2. EXPANSION MODE (Triggered by a specific path, e.g., "vfs"):
   - You perform absolute vertical decomposition of the requested path.
   - Every node must have: @ (input) -> return_type AND # responsibility_tag.
   - Use 'std.*' and 'project.*' notation for all components.

CRITICAL CONSTRAINTS:
- RULE 1: THE ISA BOUNDARY: Terminate decomposition at the Instruction Set Architecture (ISA) level. Do NOT descend into microarchitecture or transistors.
- RULE 2: THE LANGUAGE PRIMITIVE BOUNDARY: Terminate if a node is a language primitive (e.g., pointer deref, integer math, string split).
- RULE 3: MANDATORY FQN NOTATION: Every single component in every expanded tree MUST be identified by its Fully Qualified Name (FQN). Never use relative names like 'lookup'; always use the complete path from the root.`

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

type ChatResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Session History: The "Memory" for structural context
	history := []ChatMessage{
		{Role: "system", Content: SYSTEM_COMPOSITION_PROMPT},
	}

	fmt.Println("=== Requirement Decomposition Engine ===")

	// PHASE 1: Get initial requirement from user
	fmt.Print("\nEnter your requirement: ")
	requirement, _ := reader.ReadString('\n')
	requirement = strings.TrimSpace(requirement)

	if requirement == "" {
		fmt.Println("No requirement provided. Exiting.")
		return
	}

	// Send initial requirement to LLM for decomposition
	history = append(history, ChatMessage{Role: "user", Content: requirement})
	fmt.Printf("\nDecomposing: %s...\n", requirement)

	response, err := callLLM(history)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	history = append(history, ChatMessage{Role: "assistant", Content: response})

	fmt.Println("\n--- Initial Decomposition ---")
	fmt.Println(response)
	fmt.Println("-----------------------------\n")

	// PHASE 2: Expansion loop for drilling down into nodes
	fmt.Println("Enter a node path to expand (or 'exit' to quit):")

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		path := strings.TrimSpace(input)

		if path == "exit" || path == "quit" {
			break
		}
		if path == "" {
			continue
		}

		history = append(history, ChatMessage{Role: "user", Content: path})
		fmt.Printf("Expanding: %s...\n", path)

		response, err := callLLM(history)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			history = history[:len(history)-1]
			continue
		}

		history = append(history, ChatMessage{Role: "assistant", Content: response})

		fmt.Println("\n--- Expanded Tree ---")
		fmt.Println(response)
		fmt.Println("----------------------\n")
	}
}

func callLLM(messages []ChatMessage) (string, error) {
	payload := ChatRequest{
		Model:    MODEL_NAME,
		Messages: messages,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", API_URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Server returned %d. Raw Body: %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("JSON Parse Error: %v | Raw Body: %s", err, string(body))
	}

	if chatResp.Error != nil && chatResp.Error.Message != "" {
		return "", fmt.Errorf("LLM Provider Error: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices returned by model. Raw Body: %s", string(body))
	}

	return chatResp.Choices[0].Message.Content, nil
}
