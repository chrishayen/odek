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

const (
	BASE_URL      = "http://localhost:1234/v1"
	MODEL_NAME    = "gpt-4o"
	MAX_DEPTH     = 3
	MAX_TOTAL_RUNES = 100
)

const SYSTEM_PROMPT = `**Role & Objective**
You are an architectural decomposition engine. Your goal is to help users break down software requirements into small, efficient chunks called **"Runes."** 

**Core Concept: The Rune**
A Rune is the smallest unit of work. Each Rune must contain exactly **one behavior** and include:

1.  **Description:** A clear explanation of what the rune does.
2.  **Function Signature:** Using the specific type system defined below.
3.  **Positive Tests:** Scenarios where the function succeeds.
4.  **Negative Tests:** Scenarios where the function fails or handles edge cases.
5.  **Assumptions:** Explicitly list any assumptions made.

**Type System & Signatures**
*   **Integers:** i8, i16, i32, i64 (Signed) | u8, u16, u32, u64 (Unsigned)
*   **Floating Point:** f32, f64
*   **Primitives:** string, bool, bytes
*   **Collections:** list[T], map[K, V]
*   **Nullable:** optional[T]
*   **Fallible:** result[T, E]
*   **Void:** void

**CRITICAL OUTPUT FORMAT**: You MUST return ONLY valid JSON matching this exact schema - no markdown:
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

type Rune struct {
	Description   string   `json:"description" jsonschema:"description=Clear explanation of what the rune does"`
	FunctionSig   string   `json:"function_signature" jsonschema:"description=Function signature using the defined type system"`
	PositiveTests []string `json:"positive_tests" jsonschema:"description=Scenarios where the function succeeds"`
	NegativeTests []string `json:"negative_tests" jsonschema:"description=Scenarios where the function fails"`
	Assumptions   []string `json:"assumptions" jsonschema:"description=Explicit assumptions made"`
}

type PackageNode struct {
	Name     string          `json:"name"`
	Runes    map[string]Rune `json:"runes"`
	Children []PackageNode   `json:"children,omitempty"`
}

type DecompositionResponse struct {
	ProjectPackage PackageNode  `json:"project_package"`
	StdPackage     *PackageNode `json:"std_package,omitempty"`
}

type Client struct {
	instructorClient *instructor_openai.InstructorOpenAI
	conversation     *core.Conversation
}

func main() {
	ctx := context.Background()
	client := initInstructorClient(ctx)

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== Auto-Recursive Rune Decomposition Engine ===")
	fmt.Printf("Max depth: %d, Max total runes: %d\n\n", MAX_DEPTH, MAX_TOTAL_RUNES)

	fmt.Print("Enter your requirement: ")
	requirement, _ := reader.ReadString('\n')
	requirement = strings.TrimSpace(requirement)

	if requirement == "" {
		fmt.Println("No requirement provided. Exiting.")
		return
	}

	conversation := core.NewConversation(strings.TrimSpace(SYSTEM_PROMPT))
	conversation.AddUserMessage(requirement)

	fmt.Printf("\nDecomposing: %s...\n", requirement)

	response, err := client.Decompose(ctx, conversation)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	rootDecomposition := &AutoDecomposition{
		Path:         "root",
		Depth:        0,
		Response:     response,
		ParentPath:   "",
		ChildPaths:   make([]string, 0),
	}

	expansionQueue := collectRunesForExpansion(response)
	for _, runeInfo := range expansionQueue {
		runeInfo.ParentDecomposition = rootDecomposition
	}

	allDecompositions := []*AutoDecomposition{rootDecomposition}
	totalRunesCount := countTotalRunes(response)
	visitedRunePaths := make(map[string]bool)
	visitedRunePaths["root"] = true

	fmt.Printf("\n🔄 Starting auto-recursion (depth 0: %d runes)\n", len(expansionQueue))

	for len(expansionQueue) > 0 {
		if totalRunesCount >= MAX_TOTAL_RUNES {
			fmt.Printf("\n⚠️  Max total runes (%d) reached. Stopping expansion.\n", MAX_TOTAL_RUNES)
			break
		}

		runeInfo := expansionQueue[0]
		expansionQueue = expansionQueue[1:]

		if visitedRunePaths[runeInfo.FullPath] || runeInfo.Depth >= MAX_DEPTH {
			continue
		}

		visitedRunePaths[runeInfo.FullPath] = true

		fmt.Printf("   ├─ Expanding [%d/%d] %s...\n", runeInfo.Depth+1, MAX_DEPTH, runeInfo.FullPath)

		extendedReq := fmt.Sprintf(`Based on the previous decomposition, now decompose rune "%s" into sub-runes if applicable. Return ONLY valid JSON matching this structure:
{
  "project_package": {"name": "string", "runes": {"rune_name": {"description": "...", "function_signature": "...", "positive_tests": [], "negative_tests": [], "assumptions": []}}},
  "std_package": {"name": "std", "runes": {...}}
}
Do not return markdown, only raw JSON.`, runeInfo.FullPath)

		conversation.AddUserMessage(extendedReq)

		expandedResponse, err := client.Decompose(ctx, conversation)
		if err != nil {
			fmt.Printf("      ⚠️  ERROR: %v\n", err)
			continue
		}

		childDecomposition := &AutoDecomposition{
			Path:                runeInfo.FullPath,
			Depth:               runeInfo.Depth + 1,
			Response:            expandedResponse,
			ParentPath:          runeInfo.ParentDecomposition.Path,
			ChildPaths:          make([]string, 0),
		}
		allDecompositions = append(allDecompositions, childDecomposition)

		if runeInfo.ParentDecomposition != nil {
			runeInfo.ParentDecomposition.ChildPaths = append(runeInfo.ParentDecomposition.ChildPaths, runeInfo.FullPath)
		}

		newRunes := collectRunesForExpansion(expandedResponse)
		for _, newRune := range newRunes {
			newRune.Depth = runeInfo.Depth + 1
			newRune.ParentDecomposition = childDecomposition
		}
		expansionQueue = append(expansionQueue, newRunes...)

		totalRunesCount += countTotalRunes(expandedResponse)
	}

	fmt.Printf("\n")
	printCompleteTree(allDecompositions, "root", 0, true)

	separator := strings.Repeat("=", 70)
	fmt.Printf("\n%s\n", separator)
	fmt.Printf("📊 SUMMARY: %d decompositions, %d runes discovered at depth %d\n", len(allDecompositions), totalRunesCount, MAX_DEPTH)
	fmt.Printf("%s\n", separator)
}

type RuneExpansionInfo struct {
	FullPath            string
	Depth               int
	ParentDecomposition *AutoDecomposition
}

type AutoDecomposition struct {
	Path         string
	Depth        int
	Response     *DecompositionResponse
	ParentPath   string
	ChildPaths   []string
}

func initInstructorClient(ctx context.Context) *Client {
	config := openai.DefaultConfig("ollama")
	config.BaseURL = BASE_URL

	instructorClient := instructor_openai.FromOpenAI(
		openai.NewClientWithConfig(config),
		core.WithMode(core.ModeToolCall),
		core.WithMaxRetries(3),
		core.WithLogging("debug"),
	)

	return &Client{instructorClient: instructorClient, conversation: nil}
}

func (c *Client) Decompose(ctx context.Context, conv *core.Conversation) (*DecompositionResponse, error) {
	var response DecompositionResponse

	_, err := c.instructorClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    MODEL_NAME,
			Messages: instructor_openai.ConversationToMessages(conv),
		},
		&response,
	)

	if err != nil {
		return nil, fmt.Errorf("structured extraction failed: %w", err)
	}

	c.conversation = conv
	return &response, nil
}

func collectRunesForExpansion(resp *DecompositionResponse) []RuneExpansionInfo {
	var runes []RuneExpansionInfo

	if resp.ProjectPackage.Runes != nil {
		for name := range resp.ProjectPackage.Runes {
			path := fmt.Sprintf("%s.%s", resp.ProjectPackage.Name, name)
			runes = append(runes, RuneExpansionInfo{FullPath: path, Depth: 1})
		}
	}

	if resp.StdPackage != nil && resp.StdPackage.Runes != nil {
		for name := range resp.StdPackage.Runes {
			path := fmt.Sprintf("%s.%s", resp.StdPackage.Name, name)
			runes = append(runes, RuneExpansionInfo{FullPath: path, Depth: 1})
		}
	}

	return runes
}

func countTotalRunes(resp *DecompositionResponse) int {
	count := 0
	if resp.ProjectPackage.Runes != nil {
		count += len(resp.ProjectPackage.Runes)
	}
	if resp.StdPackage != nil && resp.StdPackage.Runes != nil {
		count += len(resp.StdPackage.Runes)
	}
	return count
}

func printCompleteTree(allDecompositions []*AutoDecomposition, path string, depth int, isRoot bool) {
	var decomposition *AutoDecomposition
	for _, d := range allDecompositions {
		if d.Path == path {
			decomposition = d
			break
		}
	}

	if decomposition == nil || decomposition.Response == nil {
		return
	}

	resp := decomposition.Response

	if isRoot {
		fmt.Printf("🌳 ROOT DECOMPOSITION: %s\n", path)
	} else {
		indent := strings.Repeat("   ", depth)
		fmt.Printf("%s🔸 EXPANDED: %s\n", indent, path)
	}

	if resp.ProjectPackage.Runes != nil && len(resp.ProjectPackage.Runes) > 0 {
		pkgHeader := fmt.Sprintf("   📦 %s", resp.ProjectPackage.Name)
		if !isRoot {
			pkgHeader = strings.Repeat("   ", depth) + pkgHeader
		}
		fmt.Printf("%s\n", pkgHeader)
		printRunesIndented(resp.ProjectPackage.Runes, depth+1)
	}

	if resp.StdPackage != nil && resp.StdPackage.Runes != nil && len(resp.StdPackage.Runes) > 0 {
		pkgHeader := fmt.Sprintf("   📚 %s", resp.StdPackage.Name)
		if !isRoot {
			pkgHeader = strings.Repeat("   ", depth) + pkgHeader
		}
		fmt.Printf("%s\n", pkgHeader)
		printRunesIndented(resp.StdPackage.Runes, depth+1)
	}

	for _, childPath := range decomposition.ChildPaths {
		printCompleteTree(allDecompositions, childPath, depth+1, false)
	}
}

func printRunesIndented(runes map[string]Rune, indentLevel int) {
	if len(runes) == 0 {
		return
	}

	indent := strings.Repeat("   ", indentLevel)

	for name, rune := range runes {
		fmt.Printf("%s├─ %s\n", indent, name)
		if rune.Description != "" {
			descIndent := strings.Repeat("   ", indentLevel+1)
			wrappedDesc := wrapText(rune.Description, 70-len(descIndent))
			fmt.Printf("%s│  └─ %s\n", descIndent, wrappedDesc)
		}
		if rune.FunctionSig != "" {
			sigIndent := strings.Repeat("   ", indentLevel+1)
			fmt.Printf("%s│     sig: %s\n", sigIndent, rune.FunctionSig)
		}
	}
}

func wrapText(text string, maxWidth int) string {
	if len(text) <= maxWidth {
		return text
	}
	return text[:maxWidth-3] + "..."
}
