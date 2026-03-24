package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/chrishayen/valkyrie/internal/analyzer"
	runepkg "github.com/chrishayen/valkyrie/internal/rune"
	"github.com/chrishayen/valkyrie/internal/runner"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start MCP server (stdio transport)",
	RunE: func(cmd *cobra.Command, args []string) error {
		s := server.NewMCPServer(
			"valkyrie",
			"0.1.0",
			server.WithToolCapabilities(true),
		)

		s.AddTool(mcp.NewTool("runes_list",
			mcp.WithDescription("List all runes in the registry"),
		), handleRunesList)

		s.AddTool(mcp.NewTool("runes_create",
			mcp.WithDescription("Create a new rune"),
			mcp.WithString("name", mcp.Description("Rune name (slug, optionally prefixed with feature e.g. auth/validate-email)"), mcp.Required()),
			mcp.WithString("description", mcp.Description("What the function does"), mcp.Required()),
			mcp.WithString("signature", mcp.Description("Function signature, e.g. (email: string) -> result[bool, string]"), mcp.Required()),
			mcp.WithString("behavior", mcp.Description("Inputs, outputs, edge cases")),
		), handleRunesCreate)

		s.AddTool(mcp.NewTool("runes_get",
			mcp.WithDescription("Get a rune by name"),
			mcp.WithString("name", mcp.Description("Rune name"), mcp.Required()),
		), handleRunesGet)

		s.AddTool(mcp.NewTool("runes_update",
			mcp.WithDescription("Update a rune's description, signature, or version"),
			mcp.WithString("name", mcp.Description("Rune name"), mcp.Required()),
			mcp.WithString("description", mcp.Description("New description")),
			mcp.WithString("signature", mcp.Description("New function signature")),
			mcp.WithString("version", mcp.Description("New version")),
		), handleRunesUpdate)

		s.AddTool(mcp.NewTool("runes_delete",
			mcp.WithDescription("Delete a rune"),
			mcp.WithString("name", mcp.Description("Rune name"), mcp.Required()),
		), handleRunesDelete)

		s.AddTool(mcp.NewTool("runes_analyze",
			mcp.WithDescription("Decompose requirements into runes via sandbox agent. Returns proposed runes grouped by namespace for approval."),
			mcp.WithString("requirements", mcp.Description("Plain-text requirements to decompose"), mcp.Required()),
		), handleRunesAnalyze)

		s.AddTool(mcp.NewTool("runes_create_batch",
			mcp.WithDescription("Create multiple runes at once. Accepts a JSON array of rune objects."),
			mcp.WithString("runes", mcp.Description("JSON array of rune objects, each with name, description, and optionally behavior, positive_tests, negative_tests"), mcp.Required()),
		), handleRunesCreateBatch)

		s.AddTool(mcp.NewTool("runes_hydrate",
			mcp.WithDescription("Generate code for a rune via sandbox agent"),
			mcp.WithString("name", mcp.Description("Rune name to hydrate"), mcp.Required()),
		), handleRunesHydrate)

		return server.ServeStdio(s)
	},
}

func toJSON(v any) string {
	data, _ := json.MarshalIndent(v, "", "  ")
	return string(data)
}

func textResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.NewTextContent(text),
		},
	}
}

func errResult(err error) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{
			mcp.NewTextContent(err.Error()),
		},
	}
}

func handleRunesList(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	runes, err := store.List()
	if err != nil {
		return errResult(err), nil
	}
	if runes == nil {
		runes = []runepkg.Rune{}
	}
	return textResult(toJSON(runes)), nil
}

func handleRunesCreate(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	name, _ := args["name"].(string)
	description, _ := args["description"].(string)

	signature, _ := args["signature"].(string)
	r := runepkg.Rune{Name: name, Description: description, Signature: signature}
	if beh, ok := args["behavior"].(string); ok {
		r.Behavior = beh
	}
	if pt, ok := args["positive_tests"].([]any); ok {
		for _, v := range pt {
			if s, ok := v.(string); ok {
				r.PositiveTests = append(r.PositiveTests, s)
			}
		}
	}
	if nt, ok := args["negative_tests"].([]any); ok {
		for _, v := range nt {
			if s, ok := v.(string); ok {
				r.NegativeTests = append(r.NegativeTests, s)
			}
		}
	}
	if err := store.Create(r); err != nil {
		return errResult(err), nil
	}
	created, err := store.Get(name)
	if err != nil {
		return errResult(err), nil
	}
	return textResult(toJSON(created)), nil
}

func handleRunesGet(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	name, _ := args["name"].(string)

	r, err := store.Get(name)
	if err != nil {
		return errResult(err), nil
	}
	return textResult(toJSON(r)), nil
}

func handleRunesUpdate(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	name, _ := args["name"].(string)

	r, err := store.Get(name)
	if err != nil {
		return errResult(err), nil
	}

	changed := false
	if desc, ok := args["description"].(string); ok && desc != "" {
		r.Description = desc
		changed = true
	}
	if sig, ok := args["signature"].(string); ok && sig != "" {
		r.Signature = sig
		changed = true
	}
	if ver, ok := args["version"].(string); ok && ver != "" {
		r.Version = ver
		changed = true
	}
	if !changed {
		return errResult(fmt.Errorf("at least one of description, signature, or version is required")), nil
	}

	if err := store.Update(*r); err != nil {
		return errResult(err), nil
	}
	updated, err := store.Get(name)
	if err != nil {
		return errResult(err), nil
	}
	return textResult(toJSON(updated)), nil
}

func handleRunesDelete(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	name, _ := args["name"].(string)

	if err := store.Delete(name); err != nil {
		return errResult(err), nil
	}
	return textResult(fmt.Sprintf("rune %q deleted", name)), nil
}

func handleRunesAnalyze(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	requirements, _ := args["requirements"].(string)

	run, err := runner.New(cfg.Agent)
	if err != nil {
		return errResult(err), nil
	}

	result, err := ana.Analyze(ctx, requirements, run, os.Stderr)
	if err != nil {
		return errResult(err), nil
	}

	// Auto-create the proposed runes
	var created []string
	for _, p := range result.NewRunes {
		r := p.ToRune()
		if err := store.Create(r); err != nil {
			created = append(created, fmt.Sprintf("FAILED %s: %v", p.Name, err))
			continue
		}
		created = append(created, fmt.Sprintf("created %s", p.Name))
	}

	type analyzeOutput struct {
		*analyzer.Result
		Created []string `json:"created"`
	}

	return textResult(toJSON(analyzeOutput{Result: result, Created: created})), nil
}

func handleRunesHydrate(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	name, _ := args["name"].(string)

	run, err := runner.New(cfg.Agent)
	if err != nil {
		return errResult(err), nil
	}

	result, err := hyd.Hydrate(ctx, name, run, os.Stderr)
	if err != nil {
		return errResult(err), nil
	}
	return textResult(toJSON(result)), nil
}

func handleRunesCreateBatch(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	runesJSON, _ := args["runes"].(string)

	var batch []runepkg.Rune
	if err := json.Unmarshal([]byte(runesJSON), &batch); err != nil {
		return errResult(fmt.Errorf("invalid runes JSON: %w", err)), nil
	}

	var results []string
	for _, r := range batch {
		if err := store.Create(r); err != nil {
			results = append(results, fmt.Sprintf("FAILED %s: %v", r.Name, err))
			continue
		}
		results = append(results, fmt.Sprintf("created %s", r.Name))
	}

	return textResult(toJSON(results)), nil
}
