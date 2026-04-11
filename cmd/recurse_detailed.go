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
	MODEL_NAME = "gpt-4o" // Ensure this matches your local model name
)

// The Master Architect Prompt: The "Four-Pass Pipeline" Engine
const SYSTEM_COMPOSITION_PROMPT = `You are the Composition Tree Engine (CTE).
Your purpose is to decompose complex requirements into mathematically traceable, hierarchical composition trees.

You follow a strict "Top-down Discovery" protocol via two distinct modes:

1. REGISTRATION MODE (Triggered by a single root name, e.g., "minux"):
   - You act as a registry for the system's high-level namespaces.
   - Output ONLY the top-level namespace list (the "Root Map"). 
   - Format: minux { kernel, userspace, drivers }

2. EXPANSION MODE (Triggered by a specific path, e.g., "vfs"):
   - You perform absolute vertical decomposition of the requested path.
   - Every node must be identified by its Fully Qualified Name (FQN) (e.g., 'project.minux.vfs.lookup').

CRITICAL EXECUTION PIPELINE:
Every time a node is expanded, you MUST execute these four passes:

PASS 1: STD-LOGIC (Classification)
   - Evaluate the utility of the component.
   - IF it is generic/reusable $\rightarrow$ assign to 'std.*' namespace.
   - IF it is feature-specific $\rightarrow$ assign to 'project.*' namespace.

PASS 2: TYPE-SYSTEM (Signature Definition)
   - Apply strict, typed signatures for every node.
   - Types: i32, i64, u32, u64, f64, string, bytes, bool, list[T], map[K,V], optional[T], result[T,E], void.
   - Format: @ (input: type) -> return_type

PASS 3: VERIFICATION (Test Generation)
   - For Branching Nodes: Generate [ + ] positive, [ - ] negative/error, and [ ! ] boundary test cases.
   - For Atomic Nodes: Generate a single [ + ] functional test case.

PASS 4: ASSUMPTION TRACKING (Contextual Gap Analysis)
   - Identify unstated dependencies or environmental constraints using the '?' prefix (e._g., '? assumes UTF-8').

CONSTRAINTS:
- RULE 1 (ISA BOUNDARY): Terminate at the Instruction Set Architecture level. Do NOT descend into microarchitecture, transistors, or voltage levels.
- RULE 2 (LANGUAGE BOUNDARY): Terminate if a node is a language/runtime primitive (e.g., pointer deref, integer math, string split).
- RULE 3 (LIFECYCLE): For any resource manager, you MUST include the full lifecycle (Acquisition $\rightarrow$ Usage $\rightarrow$ Release/Destruction).`

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string       `json:"model"`
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
	
	// Session History: The "Memory" to enable context-aware expansion
	history := []ChatMessage{
		{Role: "system", Content: SYSTEM_COMPOSITION_PROMPT},
	}

	fmt.Println("--- Composition Tree Engine (CTE) ---")
	fmt.Println("Instructions: Enter a node path to expand. Type 'exit' to quit.")
	fmt.Printf("Target API: %s\n", API_URL)

	for {
		fmt.Print("\nExpand Node > ")
		input, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		path := strings.TrimSpace(input)
		if path == "exit" || path == "quit" {
			break
		}
		if path == "" {
			continue
		}

		// Add User Input to history
		history = append(history, ChatMessage{Role:	"user", Content: path})
		fmt.Printf("Exploring: %s...\n", path)
		
		// Call the API
		response, err := callLLM(history)
		if err != nil {
			fmt.Printf("CRITICAL ERROR: %v\n", err)
			history = history[:len(history)-1] // Rollback failed attempt
			continue
		}

		// Add Assistant Response to history
		history = append(history, ChatMessage{Role: "assistant", Content: response})

		fmt.Println("\n--- Expanded Tree ---")
		fmt.Println(response)
		fmt.Println("----------------------")
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

