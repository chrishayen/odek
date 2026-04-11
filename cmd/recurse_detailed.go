package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/567-labs/instructor-go/pkg/instructor/core"
	instructor_openai "github.com/567-labs/instructor-go/pkg/instructor/providers/openai"
	"github.com/sashabaranov/go-openai"
)

// Configuration for the OpenAI-compatible API
const (
	BASE_URL       = "http://localhost:1234/v1"
	MODEL_NAME     = "gpt-4o"
	MAX_EXPANSIONS = 3 // Maximum number of expansion iterations
)

// The Rune-Based System Prompt
const SYSTEM_PROMPT = `
**Role & Objective**
You are an architectural decomposition engine. Your goal is to help users break down software requirements into small, efficient chunks called **"Runes."** You approach every request as if building reusable, production-grade software components rather than one-off scripts.

**Core Concept: The Rune**
A Rune is the smallest unit of work. Each Rune must contain exactly **one behavior** and include the following five elements:

1.  **Description:** A clear explanation of what the rune does.
2.  **Function Signature:** Using the specific type system defined below.
3.  **Positive Tests:** Scenarios where the function succeeds.
4.  **Negative Tests:** Scenarios where the function fails or handles edge cases.
5.  **Assumptions:** Explicitly list any assumptions made in lieu of user specification (e.g., "Assumes input is UTF-8 encoded").

**Architecture & Composition**
Runes are arranged in **composition trees** representing complexity from most to least. The hierarchy uses dot notation (parent.child.grandchild).

*Example Tree Structure:*
one.two.three
one.two.four
*Meaning: one is composed of two; two is composed of three and four.*

**Package Strategy: Reusability First**
You must distinguish between requirement-specific code and reusable utilities.

1.  **Top-Level Requirement Package:** Contains the specific logic for the user's immediate request (e.g., hello_world, my_http_server).
2.  **std (Standard Library):** A cumulative package for reusable components. When analyzing requirements, identify generic behaviors that could be reused in future projects and place them here.

*Example:* If a user asks for an HTTP server, do not build a one-off server inside the project folder. Instead, create std.httpd and import it into the top-level requirement package.

**Type System & Signatures**
When generating function signatures, strictly use the following type system:

*   **Integers:** i8, i16, i32, i64 (Signed) | u8, u16, u32, u64 (Unsigned)
*   **Floating Point:** f32, f64
*   **Primitives:** string, bool, bytes
*   **Collections:** list[T], map[K, V]
*   **Nullable:** optional[T]
*   **Fallible:** result[T, E]
*   **Void:** void
*   **Nested Types Allowed:** e.g., result[list[i32], string]

**Examples of Decomposition**

*Example 1: Hello World*
User Input: "decompose 'hello world'"

Output Structure:
hello_world # Top level package (Requirement Specific)
└── say_hello() -> void

std # Top level accumulative std package (Reusable)
└── io.write_output(string) -> result[void, string]

*Example 2: HTTP Server*
User Input: "Build an HTTP server"

Output Structure:
my_http_server # Top level package
├── start_server(port: u16) -> result[void, string]
└── handle_request(req: bytes) -> list[u8]

std # Reusable components
└── httpd.server(port: u16) -> result[void, string]

**Instructions for Output**
When the user provides a requirement:
1.  Analyze the request to identify distinct behaviors.
2.  Determine which behaviors are generic enough to belong in std and which are specific to the project root.
3.  Construct the composition tree using dot notation.
4.  For each Rune identified, provide the Description, Signature, Tests (Positive/Negative), and Assumptions.

**CRITICAL OUTPUT FORMAT**: You MUST return ONLY valid JSON matching this exact schema - no markdown, no explanations:
{
  "project_package": {
    "name": "string",
    "runes": {
      "rune_name": {
        "description": "string",
        "function_signature": "string",
        "positive_tests": ["string"],
        "negative_tests": ["string"],
        "assumptions": ["string"]
      }
    }
  },
  "std_package": {
    "name": "std",
    "runes": {}
  }
}
`

// Rune represents the smallest unit of work - exactly one behavior
type Rune struct {
	Description   string   `json:"description" jsonschema:"description=Clear explanation of what the rune does"`
	FunctionSig   string   `json:"function_signature" jsonschema:"description=Function signature using the defined type system (e.g., 'add(a: i32, b: i32) -> result[i32, string]')"`
	PositiveTests []string `json:"positive_tests" jsonschema:"description=Scenarios where the function succeeds"`
	NegativeTests []string `json:"negative_tests" jsonschema:"description=Scenarios where the function fails or handles edge cases"`
	Assumptions   []string `json:"assumptions" jsonschema:"description=Explicit assumptions made in lieu of user specification"`
}

// PackageNode represents a package in the composition tree
type PackageNode struct {
	Name     string          `json:"name" jsonschema:"description=Package name (e.g., 'hello_world' or 'std')"`
	Runes    map[string]Rune `json:"runes" jsonschema:"description=Map of rune names to their definitions"`
	Children []PackageNode   `json:"children,omitempty" jsonschema:"description=Nested sub-packages"`
}

// DecompositionResponse is the structured output from the LLM
type DecompositionResponse struct {
	ProjectPackage PackageNode  `json:"project_package" jsonschema:"description=Top-level requirement-specific package containing runes specific to this request"`
	StdPackage     *PackageNode `json:"std_package,omitempty" jsonschema:"description=Cumulative standard library package for reusable components that could be used in future projects"`
}

// Client holds the instructor-go client
type Client struct {
	instructorClient *instructor_openai.InstructorOpenAI
}

func main() {
	ctx := context.Background()

	// Initialize instructor-go client with custom endpoint and debug logging
	client := initInstructorClient(ctx)

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== Rune-Based Architectural Decomposition Engine ===")
	fmt.Println("Decomposes requirements into reusable 'Runes' (smallest units of work)")
	fmt.Printf("Max expansions: %d\n\n", MAX_EXPANSIONS)

	// PHASE 1: Get initial requirement from user
	fmt.Print("Enter your requirement: ")
	requirement, _ := reader.ReadString('\n')
	requirement = strings.TrimSpace(requirement)

	if requirement == "" {
		fmt.Println("No requirement provided. Exiting.")
		return
	}

	// Build conversation with system prompt
	conversation := buildConversation(requirement)

	// Call LLM for structured decomposition
	fmt.Printf("\nDecomposing: %s...\n", requirement)

	response, err := client.Decompose(ctx, conversation)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	// Print structured response
	printDecomposition(response, 0)

	// PHASE 2: Expansion loop for drilling down into specific runes
	fmt.Println("\nEnter a rune path to expand details (or 'exit' to quit):")

	expansionCount := 0
	for {
		if expansionCount >= MAX_EXPANSIONS {
			fmt.Printf("\n⚠️  Maximum expansions (%d) reached. Exiting.\n", MAX_EXPANSIONS)
			break
		}

		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		path := strings.TrimSpace(input)

		if path == "exit" || path == "quit" {
			break
		}
		if path == "" {
			continue
		}

		expansionCount++
		fmt.Printf("\n[Expansion %d/%d] Analyzing rune: %s...\n", expansionCount, MAX_EXPANSIONS, path)

		// Expand specific rune - explicitly request structured JSON output
		extendedReq := fmt.Sprintf(`Decompose rune "%s" into sub-runes if applicable. Return ONLY valid JSON matching this structure:
{
  "project_package": {"name": "string", "runes": {"rune_name": {"description": "...", "function_signature": "...", "positive_tests": [], "negative_tests": [], "assumptions": []}}},
  "std_package": {"name": "std", "runes": {...}}
}
Do not return markdown, explanations, or plain text - only raw JSON.`, path)
		conversation = buildConversation(extendedReq)

		expandedResponse, err := client.Decompose(ctx, conversation)
		if err != nil {
			fmt.Printf("ERROR expanding rune: %v\n", err)
			continue
		}

		printDecomposition(expandedResponse, expansionCount)
	}
}

// initInstructorClient creates and configures the instructor-go client
func initInstructorClient(ctx context.Context) *Client {
	// Configure OpenAI client with custom base URL
	config := openai.DefaultConfig("ollama") // Placeholder API key
	config.BaseURL = BASE_URL

	instructorClient := instructor_openai.FromOpenAI(
		openai.NewClientWithConfig(config),
		core.WithMode(core.ModeToolCall), // Highest accuracy - uses function calling
		core.WithMaxRetries(3),           // Automatic retries on validation failure
		core.WithLogging("debug"),        // Debug mode enabled
	)

	return &Client{instructorClient: instructorClient}
}

// buildConversation creates the conversation with system prompt and user message
func buildConversation(userMessage string) *core.Conversation {
	conv := core.NewConversation(strings.TrimSpace(SYSTEM_PROMPT))
	conv.AddUserMessage(userMessage)
	return conv
}

// Decompose calls the LLM and returns structured decomposition
func (c *Client) Decompose(ctx context.Context, conversation *core.Conversation) (*DecompositionResponse, error) {
	var response DecompositionResponse

	_, err := c.instructorClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    MODEL_NAME,
			Messages: instructor_openai.ConversationToMessages(conversation),
		},
		&response, // Pass pointer for structured extraction
	)

	if err != nil {
		return nil, fmt.Errorf("structured extraction failed: %w", err)
	}

	return &response, nil
}

// printDecomposition formats and prints the decomposition response
func printDecomposition(resp *DecompositionResponse, expansionLevel int) {
	separator := strings.Repeat("=", 60)

	if expansionLevel > 0 {
		fmt.Println()
		fmt.Println("EXPANDED DECOMPOSITION")
		fmt.Printf("%s\n", separator)
	}

	// Print Project Package
	fmt.Printf("\n📦 PROJECT PACKAGE: %s\n", resp.ProjectPackage.Name)
	fmt.Println(strings.Repeat("-", 40))
	printRunes(resp.ProjectPackage.Runes, "  ")

	// Print Std Package if present
	if resp.StdPackage != nil {
		fmt.Printf("\n📚 STD PACKAGE: %s\n", resp.StdPackage.Name)
		fmt.Println(strings.Repeat("-", 40))
		printRunes(resp.StdPackage.Runes, "  ")
	}

	fmt.Println(separator)
}

// printRunes recursively prints runes with indentation
func printRunes(runes map[string]Rune, indent string) {
	if len(runes) == 0 {
		fmt.Printf("%s(No runes defined)\n", indent)
		return
	}

	for name, rune := range runes {
		fmt.Printf("%s🔹 %s\n", indent, name)
		fmt.Printf("%s   Description: %s\n", indent, rune.Description)
		fmt.Printf("%s   Signature:   %s\n", indent, rune.FunctionSig)

		if len(rune.PositiveTests) > 0 {
			fmt.Printf("%s   Positive Tests:\n", indent)
			for _, test := range rune.PositiveTests {
				fmt.Printf("%s     ✓ %s\n", indent, test)
			}
		}

		if len(rune.NegativeTests) > 0 {
			fmt.Printf("%s   Negative Tests:\n", indent)
			for _, test := range rune.NegativeTests {
				fmt.Printf("%s     ✗ %s\n", indent, test)
			}
		}

		if len(rune.Assumptions) > 0 {
			fmt.Printf("%s   Assumptions:\n", indent)
			for _, assumption := range rune.Assumptions {
				fmt.Printf("%s     • %s\n", indent, assumption)
			}
		}

		fmt.Println()
	}
}
