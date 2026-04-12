package main

import (
	"bufio"
	"context"
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/567-labs/instructor-go/pkg/instructor/core"
	instructor_openai "github.com/567-labs/instructor-go/pkg/instructor/providers/openai"
	"github.com/sashabaranov/go-openai"
)

const (
	BASE_URL        = "http://localhost:1234/v1"
	MODEL_NAME      = "gpt-4o"
	MAX_DEPTH       = 3
	MAX_TOTAL_RUNES = 100
)

//go:embed decompose.md
var SYSTEM_PROMPT string

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
	Type           string       `json:"type" jsonschema:"const=decompose"`
	ProjectPackage PackageNode  `json:"project_package"`
	StdPackage     *PackageNode `json:"std_package,omitempty"`
}

type ClarificationRequest struct {
	Type    string `json:"type" jsonschema:"const=clarify"`
	Message string `json:"message" jsonschema:"description=The explanation if the requirement is too vague"`
}

type Client struct {
	instructorClient *instructor_openai.InstructorOpenAI
	conversation     *core.Conversation
}

// func (c *Client) Decompose(ctx context.Context, conv *core.Conversation) (any, error) {
// 	var action any

// 	err := c.instructorClient.CreateChatCompletionUnion(
// 		ctx,
// 		openai.ChatCompletionRequest{
// 			Model:    MODEL_NAME,
// 			Messages: instructor_openai.ConversationToMessages(conv),
// 		},
// 		core.UnionOptions{
// 			Discriminator: "type",
// 			Variants:      []any{DecompositionResponse{}, ClarificationRequest{}},
// 		},
// 		&action,
// 	)

// 	if err != nil {
// 		return nil, fmt.Errorf("structured extraction failed: %w", err)
// 	}

// 	c.conversation = conv
// 	return action, nil
// }

type RuneExpansionInfo struct {
	FullPath            string
	Depth               int
	ParentDecomposition *AutoDecomposition
}

type AutoDecomposition struct {
	Path       string
	Depth      int
	Response   *DecompositionResponse
	ParentPath string
	ChildPaths []string
}

var client *Client

func init() {
	config := openai.DefaultConfig("default")
	config.BaseURL = BASE_URL

	c := &Client{}
	c.instructorClient = instructor_openai.FromOpenAI(
		openai.NewClientWithConfig(config),
		core.WithMode(core.ModeToolCall),
		core.WithMaxRetries(3),
		// core.WithLogging("debug"),
	)

	client = c
}

func getRequirement(reader *bufio.Reader) (string, error) {
	fmt.Print("Enter your requirement: ")
	requirement, _ := reader.ReadString('\n')
	requirement = strings.TrimSpace(requirement)
	return requirement, nil
}

func main() {
	ctx := context.Background()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== Auto-Recursive Rune Decomposition Engine ===")
	fmt.Printf("Max depth: %d, Max total runes: %d\n\n", MAX_DEPTH, MAX_TOTAL_RUNES)

	var req string
	var err error

	if req, err = getRequirement(reader); err != nil {
		fmt.Println("No requirement provided. Exiting.")
		os.Exit(1)
	}

	conversation := core.NewConversation(strings.TrimSpace(SYSTEM_PROMPT))
	conversation.AddUserMessage("decompose: " + req)

	fmt.Printf("\nDecomposing: %s...\n", req)

	response, err := client.Decompose(ctx, conversation)

	if err != nil {
		fmt.Printf("ERROR: %v, %v\n", err, conversation.GetMessages())
		return
	}

	switch action := response.(type) {

	case ClarificationRequest:
		fmt.Printf("\n⚠️  CLARIFICATION NEEDED: %s\n", action.Message)
		return

	case DecompositionResponse:
		rootDecomposition := &AutoDecomposition{
			Path:       "root",
			Depth:      0,
			Response:   &action,
			ParentPath: "",
			ChildPaths: make([]string, 0),
		}

		// fmt.Printf("%v\n", rootDecomposition.Response)
		fmt.Printf("fuck it %v", rootDecomposition.Response)

	// case openai.ChatCompletionResponse:
	// 	println("clarification request")

	default:
		fmt.Printf("Received unknown action type: %T\n", action)
		return
	}

	// 	expansionQueue := collectRunesForExpansion(response)
	// 	for i := range expansionQueue {
	// 		expansionQueue[i].ParentDecomposition = rootDecomposition
	// 	}

	// 	allDecompositions := []*AutoDecomposition{rootDecomposition}
	// 	totalRunesCount := countTotalRunes(response)
	// 	visitedRunePaths := make(map[string]bool)
	// 	visitedRunePaths["root"] = true

	// 	fmt.Printf("\n🔄 Starting auto-recursion (depth 0: %d runes)\n", len(expansionQueue))

	// 	for len(expansionQueue) > 0 {
	// 		if totalRunesCount >= MAX_TOTAL_RUNES {
	// 			fmt.Printf("\n⚠️  Max total runes (%d) reached. Stopping expansion.\n", MAX_TOTAL_RUNES)
	// 			break
	// 		}

	// 		runeInfo := expansionQueue[0]
	// 		expansionQueue = expansionQueue[1:]

	// 		if visitedRunePaths[runeInfo.FullPath] || runeInfo.Depth >= MAX_DEPTH {
	// 			continue
	// 		}

	// 		visitedRunePaths[runeInfo.FullPath] = true

	// 		fmt.Printf("   ├─ Expanding [%d/%d] %s...\n", runeInfo.Depth+1, MAX_DEPTH, runeInfo.FullPath)

	// 		extendedReq := fmt.Sprintf(`Based on the previous decomposition, now decompose rune "%s" into sub-runes if applicable. Return ONLY valid JSON matching this structure:
	// {
	//   "project_package": {"name": "string", "runes": {"rune_name": {"description": "...", "function_signature": "...", "positive_tests": [], "negative_tests": [], "assumptions": []}}},
	//   "std_package": {"name": "std", "runes": {...}}
	// }
	// Do not return markdown, only raw JSON.`, runeInfo.FullPath)

	// 		conversation.AddUserMessage(extendedReq)

	// 		expandedResponse, err := client.Decompose(ctx, conversation)
	// 		if err != nil {
	// 			fmt.Printf("      ⚠️  ERROR: %v\n", err)
	// 			continue
	// 		}

	// 		if expandedResponse == nil || expandedResponse.ProjectPackage.Name == "" {
	// 			fmt.Printf("      ⚠️  ERROR: invalid response received, skipping expansion.\n")
	// 			continue
	// 		}

	// 		// Check if there are any runes to expand - skip empty responses (leaf nodes)
	// 		newRunes := collectRunesForExpansion(expandedResponse)
	// 		if len(newRunes) == 0 {
	// 			fmt.Printf("      ✓ No sub-runes found (leaf node).\n")
	// 			continue
	// 		}

	// 		childDecomposition := &AutoDecomposition{
	// 			Path:       runeInfo.FullPath,
	// 			Depth:      runeInfo.Depth + 1,
	// 			Response:   expandedResponse,
	// 			ParentPath: "",
	// 			ChildPaths: make([]string, 0),
	// 		}
	// 		allDecompositions = append(allDecompositions, childDecomposition)

	// 		if runeInfo.ParentDecomposition != nil {
	// 			runeInfo.ParentDecomposition.ChildPaths = append(runeInfo.ParentDecomposition.ChildPaths, runeInfo.FullPath)
	// 			childDecomposition.ParentPath = runeInfo.ParentDecomposition.Path
	// 		}

	// 		for i := range newRunes {
	// 			newRunes[i].Depth = runeInfo.Depth + 1
	// 			newRunes[i].ParentDecomposition = childDecomposition
	// 		}
	// 		expansionQueue = append(expansionQueue, newRunes...)

	// 		totalRunesCount += countTotalRunes(expandedResponse)
	// 	}

	// 	fmt.Printf("\n")
	// 	printCompleteTree(allDecompositions, "root", 0, true)

	// separator := strings.Repeat("=", 70)
	// fmt.Printf("\n%s\n", separator)
	// fmt.Printf("📊 SUMMARY: %d decompositions, %d runes discovered at depth %d\n", len(allDecompositions), totalRunesCount, MAX_DEPTH)
	// fmt.Printf("%s\n", separator)
}

func (c *Client) Decompose(ctx context.Context, conv *core.Conversation) (any, error) {
	// var action any

	result, _, err := c.instructorClient.CreateChatCompletionUnion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    MODEL_NAME,
			Messages: instructor_openai.ConversationToMessages(conv),
		},
		core.UnionOptions{
			Discriminator: "type",
			Variants:      []any{DecompositionResponse{}, ClarificationRequest{}},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("structured extraction failed: %w", err)
	}

	// fmt.Printf("res %v %v %i", res, resp, len(res))
	// fmt.Printf("res %v", res)

	c.conversation = conv
	for _, item := range result {
		if req, ok := item.(ClarificationRequest); ok {
			return req, err
		}

		if req, ok := item.(DecompositionResponse); ok {
			return req, err
		}
	}

	return nil, err

}

func collectRunesForExpansion(resp *DecompositionResponse) []RuneExpansionInfo {
	var runes []RuneExpansionInfo

	if resp == nil || resp.ProjectPackage.Name == "" {
		return runes
	}

	// Only expand if there are actually runes in the project package
	if len(resp.ProjectPackage.Runes) > 0 {
		for name := range resp.ProjectPackage.Runes {
			path := fmt.Sprintf("%s.%s", resp.ProjectPackage.Name, name)
			runes = append(runes, RuneExpansionInfo{FullPath: path, Depth: 1})
		}
	}

	// Only expand if there are actually runes in the std package
	if resp.StdPackage != nil && len(resp.StdPackage.Runes) > 0 {
		for name := range resp.StdPackage.Runes {
			path := fmt.Sprintf("%s.%s", resp.StdPackage.Name, name)
			runes = append(runes, RuneExpansionInfo{FullPath: path, Depth: 1})
		}
	}

	return runes
}

func countTotalRunes(resp *DecompositionResponse) int {
	if resp == nil {
		return 0
	}
	count := 0
	if len(resp.ProjectPackage.Runes) > 0 {
		count += len(resp.ProjectPackage.Runes)
	}
	if resp.StdPackage != nil && len(resp.StdPackage.Runes) > 0 {
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

	if len(resp.ProjectPackage.Runes) > 0 {
		pkgHeader := fmt.Sprintf("   📦 %s", resp.ProjectPackage.Name)
		if !isRoot {
			pkgHeader = strings.Repeat("   ", depth) + pkgHeader
		}
		fmt.Printf("%s\n", pkgHeader)
		printRunesIndented(resp.ProjectPackage.Runes, depth+1)
	}

	if resp.StdPackage != nil && len(resp.StdPackage.Runes) > 0 {
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
