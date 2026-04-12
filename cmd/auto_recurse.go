package main

import (
	"bufio"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

const (
	BASE_URL        = "http://localhost:8080/v1"
	MODEL_NAME      = "Carnice-27b-Q4_K_M.gguf"
	MAX_DEPTH       = 3
	MAX_TOTAL_RUNES = 100
)

//go:embed decompose.md
var SYSTEM_PROMPT string

type Rune = struct {
	Description   string   `json:"description"`
	FunctionSig   string   `json:"function_signature"`
	PositiveTests []string `json:"positive_tests"`
	NegativeTests []string `json:"negative_tests"`
	Assumptions   []string `json:"assumptions"`
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
	ProjectPackage wirePackage  `json:"project_package"`
	StdPackage     *wirePackage `json:"std_package,omitempty"`
}

type ClarificationRequest struct {
	Message string `json:"message"`
}

type Client struct {
	openai *openai.Client
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

var (
	client        *Client
	decomposeTool openai.Tool
)

func init() {
	config := openai.DefaultConfig("default")
	config.BaseURL = BASE_URL
	client = &Client{openai: openai.NewClientWithConfig(config)}

	runeSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"description":        map[string]any{"type": "string"},
			"function_signature": map[string]any{"type": "string"},
			"positive_tests":     map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
			"negative_tests":     map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
			"assumptions":        map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
		},
		"required": []string{"description", "function_signature", "positive_tests", "negative_tests", "assumptions"},
	}
	packageSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name":  map[string]any{"type": "string"},
			"runes": map[string]any{"type": "object", "additionalProperties": runeSchema},
		},
		"required": []string{"name", "runes"},
	}
	decomposeTool = openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "decompose",
			Description: "Submit a rune decomposition. Provide a project_package, and optionally a std_package of reusable utilities.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"project_package": packageSchema,
					"std_package":     packageSchema,
				},
				"required": []string{"project_package"},
			},
		},
	}
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

	messages := newConversation(req)
	fmt.Printf("\nDecomposing: %s...\n", req)

	root, queue, err := initialDecompose(ctx, &messages)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}
	if root == nil {
		return
	}

	for i := range queue {
		queue[i].ParentDecomposition = root
	}

	expandRecursively(ctx, &messages, root, queue)
}

func printBanner() {
	fmt.Println("=== Auto-Recursive Rune Decomposition Engine ===")
	fmt.Printf("Max depth: %d, Max total runes: %d\n\n", MAX_DEPTH, MAX_TOTAL_RUNES)
}

func newConversation(req string) []openai.ChatCompletionMessage {
	return []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: strings.TrimSpace(SYSTEM_PROMPT)},
		{Role: openai.ChatMessageRoleUser, Content: "decompose: " + req},
	}
}

// initialDecompose runs the first decomposition pass. Returns (nil, nil, nil)
// when the model asks for clarification, so the caller should treat a nil root
// as a clean exit.
func initialDecompose(ctx context.Context, messages *[]openai.ChatCompletionMessage) (*AutoDecomposition, []RuneExpansionInfo, error) {
	response, err := client.Decompose(ctx, messages)
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
		return root, collectRunesForExpansion(root.Response), nil

	default:
		return nil, nil, fmt.Errorf("unknown action type: %T", action)
	}
}

// expandRecursively drains the expansion queue, decomposing each rune up to
// MAX_DEPTH and stitching child decompositions back into the tree.
func expandRecursively(ctx context.Context, messages *[]openai.ChatCompletionMessage, root *AutoDecomposition, queue []RuneExpansionInfo) {
	for i := range queue {
		queue[i].ParentDecomposition = root
	}

	allDecompositions := []*AutoDecomposition{root}
	totalRunesCount := countTotalRunes(root.Response)
	visitedRunePaths := make(map[string]bool)
	visitedRunePaths["root"] = true

	fmt.Printf("\n🔄 Starting auto-recursion (depth 0: %d runes)\n", len(queue))

	for len(queue) > 0 {
		if totalRunesCount >= MAX_TOTAL_RUNES {
			fmt.Printf("\n⚠️  Max total runes (%d) reached. Stopping expansion.\n", MAX_TOTAL_RUNES)
			break
		}

		runeInfo := queue[0]
		queue = queue[1:]

		if visitedRunePaths[runeInfo.FullPath] || runeInfo.Depth >= MAX_DEPTH {
			continue
		}

		visitedRunePaths[runeInfo.FullPath] = true

		fmt.Printf("   ├─ Expanding [%d/%d] %s...\n", runeInfo.Depth+1, MAX_DEPTH, runeInfo.FullPath)

		extendedReq := fmt.Sprintf(`Decompose rune "%s" into sub-runes if applicable. Submit the decomposition by calling the decompose tool, following the same rules as the initial decomposition. If the rune is already atomic and has no meaningful sub-runes, call decompose with an empty runes map.`, runeInfo.FullPath)

		*messages = append(*messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: extendedReq,
		})

		expandedResponse, err := client.Decompose(ctx, messages)
		if err != nil {
			fmt.Printf("      ⚠️  ERROR: %v\n", err)
			continue
		}

		respVal, ok := expandedResponse.(DecompositionResponse)
		if !ok || respVal.ProjectPackage.Name == "" {
			fmt.Printf("      ⚠️  ERROR: invalid response received, skipping expansion.\n")
			continue
		}
		resp := &respVal

		newRunes := collectRunesForExpansion(resp)
		if len(newRunes) == 0 {
			fmt.Printf("      ✓ No sub-runes found (leaf node).\n")
			continue
		}

		childDecomposition := &AutoDecomposition{
			Path:       runeInfo.FullPath,
			Depth:      runeInfo.Depth + 1,
			Response:   resp,
			ParentPath: "",
			ChildPaths: make([]string, 0),
		}
		allDecompositions = append(allDecompositions, childDecomposition)

		if runeInfo.ParentDecomposition != nil {
			runeInfo.ParentDecomposition.ChildPaths = append(runeInfo.ParentDecomposition.ChildPaths, runeInfo.FullPath)
			childDecomposition.ParentPath = runeInfo.ParentDecomposition.Path
		}

		for i := range newRunes {
			newRunes[i].Depth = runeInfo.Depth + 1
			newRunes[i].ParentDecomposition = childDecomposition
		}
		queue = append(queue, newRunes...)

		totalRunesCount += countTotalRunes(resp)
	}

	fmt.Printf("\n")
	printCompleteTree(allDecompositions, "root", 0, true)

	separator := strings.Repeat("=", 70)
	fmt.Printf("\n%s\n", separator)
	fmt.Printf("📊 SUMMARY: %d decompositions, %d runes discovered at depth %d\n", len(allDecompositions), totalRunesCount, MAX_DEPTH)
	fmt.Printf("%s\n", separator)
}

func (c *Client) Decompose(ctx context.Context, messages *[]openai.ChatCompletionMessage) (any, error) {
	resp, err := c.openai.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:      MODEL_NAME,
		Messages:   *messages,
		Tools:      []openai.Tool{decomposeTool},
		ToolChoice: "auto",
	})
	if err != nil {
		return nil, fmt.Errorf("chat completion failed: %w", err)
	}
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	msg := resp.Choices[0].Message
	*messages = append(*messages, msg)

	if len(msg.ToolCalls) > 0 {
		call := msg.ToolCalls[0]
		if call.Function.Name != "decompose" {
			return nil, fmt.Errorf("unexpected tool call: %s", call.Function.Name)
		}
		var decomp DecompositionResponse
		if err := json.Unmarshal([]byte(call.Function.Arguments), &decomp); err != nil {
			return nil, fmt.Errorf("parsing decompose arguments: %w (raw: %s)", err, call.Function.Arguments)
		}
		*messages = append(*messages, openai.ChatCompletionMessage{
			Role:       openai.ChatMessageRoleTool,
			ToolCallID: call.ID,
			Content:    "decomposition recorded",
		})
		return decomp, nil
	}

	if strings.TrimSpace(msg.Content) != "" {
		return ClarificationRequest{Message: msg.Content}, nil
	}

	return nil, fmt.Errorf("model returned neither a tool call nor content")
}

func collectRunesForExpansion(resp *DecompositionResponse) []RuneExpansionInfo {
	var runes []RuneExpansionInfo

	if resp == nil || resp.ProjectPackage.Name == "" {
		return runes
	}

	if len(resp.ProjectPackage.Runes) > 0 {
		for name := range resp.ProjectPackage.Runes {
			path := name
			if !strings.HasPrefix(name, resp.ProjectPackage.Name+".") {
				path = fmt.Sprintf("%s.%s", resp.ProjectPackage.Name, name)
			}
			runes = append(runes, RuneExpansionInfo{FullPath: path, Depth: 1})
		}
	}

	if resp.StdPackage != nil && len(resp.StdPackage.Runes) > 0 {
		for name := range resp.StdPackage.Runes {
			path := name
			if !strings.HasPrefix(name, resp.StdPackage.Name+".") {
				path = fmt.Sprintf("%s.%s", resp.StdPackage.Name, name)
			}
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
