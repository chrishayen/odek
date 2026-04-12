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
	BASE_URL        = "http://localhost:8080/v1"
	MODEL_NAME      = "gemma-4-26B-A4B-it-Q4_K_M.gguf"
	MAX_DEPTH       = 3
	MAX_TOTAL_RUNES = 100
)

//go:embed decompose.md
var SYSTEM_PROMPT string

type Rune = struct {
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

type wirePackage = struct {
	Name  string          `json:"name"`
	Runes map[string]Rune `json:"runes"`
}

type DecompositionResponse struct {
	Type           string       `json:"type" jsonschema:"const=decompose"`
	ProjectPackage wirePackage  `json:"project_package"`
	StdPackage     *wirePackage `json:"std_package,omitempty"`
}

type ClarificationRequest struct {
	Type    string `json:"type" jsonschema:"const=clarify"`
	Message string `json:"message" jsonschema:"description=The explanation if the requirement is too vague"`
}

type Client struct {
	instructorClient *instructor_openai.InstructorOpenAI
	conversation     *core.Conversation
}

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
	printBanner()

	req, err := getRequirement(bufio.NewReader(os.Stdin))
	if err != nil || req == "" {
		fmt.Println("No requirement provided. Exiting.")
		os.Exit(1)
	}

	conversation := newConversation(req)
	fmt.Printf("\nDecomposing: %s...\n", req)

	root, queue, err := initialDecompose(ctx, conversation)
	if err != nil {
		fmt.Printf("ERROR: %v, %v\n", err, conversation.GetMessages())
		return
	}
	if root == nil {
		return
	}

	fmt.Printf("queue %v", queue)

	// expandRecursively(ctx, conversation, root, queue)
}

func printBanner() {
	fmt.Println("=== Auto-Recursive Rune Decomposition Engine ===")
	fmt.Printf("Max depth: %d, Max total runes: %d\n\n", MAX_DEPTH, MAX_TOTAL_RUNES)
}

func newConversation(req string) *core.Conversation {
	conv := core.NewConversation(strings.TrimSpace(SYSTEM_PROMPT))
	conv.AddUserMessage("decompose: " + req)
	return conv
}

// initialDecompose runs the first decomposition pass. Returns (nil, nil, nil)
// when the model asks for clarification, so the caller should treat a nil root
// as a clean exit.
func initialDecompose(ctx context.Context, conv *core.Conversation) (*AutoDecomposition, []RuneExpansionInfo, error) {
	response, err := client.Decompose(ctx, conv)
	if err != nil {
		return nil, nil, err
	}

	switch action := response.(type) {
	case ClarificationRequest:
		fmt.Printf("\n⚠️  CLARIFICATION NEEDED: %s\n", action.Message)
		return nil, nil, nil

	case DecompositionResponse:
		root := &AutoDecomposition{
			Path:       "root",
			Depth:      0,
			Response:   &action,
			ParentPath: "",
			ChildPaths: make([]string, 0),
		}
		fmt.Printf("decomp response %+v", root.Response)
		return root, collectRunesForExpansion(root.Response), nil

	default:
		return nil, nil, fmt.Errorf("unknown action type: %T", action)
	}
}

// expandRecursively drains the expansion queue, decomposing each rune up to
// MAX_DEPTH and stitching child decompositions back into the tree. WIP.
func expandRecursively(ctx context.Context, conv *core.Conversation, root *AutoDecomposition, queue []RuneExpansionInfo) {
	// 	for i := range queue {
	// 		queue[i].ParentDecomposition = root
	// 	}

	// 	allDecompositions := []*AutoDecomposition{root}
	// 	totalRunesCount := countTotalRunes(root.Response)
	// 	visitedRunePaths := make(map[string]bool)
	// 	visitedRunePaths["root"] = true

	// 	fmt.Printf("\n🔄 Starting auto-recursion (depth 0: %d runes)\n", len(queue))

	// 	for len(queue) > 0 {
	// 		if totalRunesCount >= MAX_TOTAL_RUNES {
	// 			fmt.Printf("\n⚠️  Max total runes (%d) reached. Stopping expansion.\n", MAX_TOTAL_RUNES)
	// 			break
	// 		}

	// 		runeInfo := queue[0]
	// 		queue = queue[1:]

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

	// 		conv.AddUserMessage(extendedReq)

	// 		expandedResponse, err := client.Decompose(ctx, conv)
	// 		if err != nil {
	// 			fmt.Printf("      ⚠️  ERROR: %v\n", err)
	// 			continue
	// 		}

	// 		if expandedResponse == nil || expandedResponse.ProjectPackage.Name == "" {
	// 			fmt.Printf("      ⚠️  ERROR: invalid response received, skipping expansion.\n")
	// 			continue
	// 		}

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
	// 		queue = append(queue, newRunes...)

	// 		totalRunesCount += countTotalRunes(expandedResponse)
	// 	}

	// 	fmt.Printf("\n")
	// 	printCompleteTree(allDecompositions, "root", 0, true)

	// 	separator := strings.Repeat("=", 70)
	// 	fmt.Printf("\n%s\n", separator)
	// 	fmt.Printf("📊 SUMMARY: %d decompositions, %d runes discovered at depth %d\n", len(allDecompositions), totalRunesCount, MAX_DEPTH)
	// 	fmt.Printf("%s\n", separator)
}

func (c *Client) Decompose(ctx context.Context, conv *core.Conversation) (any, error) {

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

	c.conversation = conv
	for _, item := range result {
		if req, ok := item.(ClarificationRequest); ok {
			return req, err
		}

		if req, ok := item.(DecompositionResponse); ok {
			fmt.Printf("%+v", req.StdPackage)
			return req, err
		}
	}

	return nil, err

}

func collectRunesForExpansion(resp *DecompositionResponse) []RuneExpansionInfo {
	var runes []RuneExpansionInfo

	fmt.Printf("runes %v", resp.ProjectPackage.Runes)

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
