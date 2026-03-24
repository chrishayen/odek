package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/chrishayen/valkyrie/internal/analyzer"
	"github.com/chrishayen/valkyrie/internal/feature"
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
			mcp.WithDescription("List all runes in the registry. A rune is the atomic unit of functionality — one function described in English. Each rune has a name, description, signature, behavior spec, and test cases. Returns JSON array of all registered runes."),
		), handleRunesList)

		s.AddTool(mcp.NewTool("runes_create",
			mcp.WithDescription("Create a new rune in the registry. A rune describes exactly one function. The name should follow verb-noun pattern and be prefixed with its feature namespace (e.g. auth/validate-email). The description, signature, and behavior together form the contract that the hydration agent uses to generate code. Before creating runes, check features_list to see if the feature namespace exists — if not, create the feature first with features_create."),
			mcp.WithString("name", mcp.Description("Rune name as a slug, prefixed with feature namespace (e.g. auth/validate-email, payment/calculate-total)"), mcp.Required()),
			mcp.WithString("description", mcp.Description("One or two sentences stating what the function does, what it accepts, and what it returns"), mcp.Required()),
			mcp.WithString("signature", mcp.Description("Function signature using precise types, e.g. (email: string) -> bool, (prices: list[f64], tax_rate: f64) -> result[f64, string]"), mcp.Required()),
			mcp.WithString("behavior", mcp.Description("Precise description of inputs, outputs, edge cases, and constraints. Each point on its own line starting with '- '")),
		), handleRunesCreate)

		s.AddTool(mcp.NewTool("runes_get",
			mcp.WithDescription("Retrieve a single rune by name. Returns the full rune record including description, signature, behavior, tests, version, hydration status, and coverage."),
			mcp.WithString("name", mcp.Description("Full rune name including namespace (e.g. auth/validate-email)"), mcp.Required()),
		), handleRunesGet)

		s.AddTool(mcp.NewTool("runes_update",
			mcp.WithDescription("Update an existing rune's description, signature, or version. At least one field must be provided. The rune must already exist in the registry."),
			mcp.WithString("name", mcp.Description("Full rune name including namespace (e.g. auth/validate-email)"), mcp.Required()),
			mcp.WithString("description", mcp.Description("New description for the rune")),
			mcp.WithString("signature", mcp.Description("New function signature with precise types")),
			mcp.WithString("version", mcp.Description("New semantic version (e.g. 0.2.0)")),
		), handleRunesUpdate)

		s.AddTool(mcp.NewTool("runes_delete",
			mcp.WithDescription("Delete a rune from the registry. Removes the rune spec file. Does not remove generated code if the rune was hydrated."),
			mcp.WithString("name", mcp.Description("Full rune name including namespace (e.g. auth/validate-email)"), mcp.Required()),
		), handleRunesDelete)

		s.AddTool(mcp.NewTool("runes_analyze",
			mcp.WithDescription("Decompose plain-text requirements into runes via a sandboxed agent. The agent reads existing runes to avoid duplication, then proposes new runes with names, descriptions, signatures, behavior specs, and test cases. Proposed runes are auto-created in the registry. Returns both new and existing runes that cover the requirements."),
			mcp.WithString("requirements", mcp.Description("Plain-text English description of what you need built. Be specific about inputs, outputs, and expected behavior."), mcp.Required()),
		), handleRunesAnalyze)

		s.AddTool(mcp.NewTool("runes_create_batch",
			mcp.WithDescription("Create multiple runes at once from a JSON array. Each rune object must have name, description, and signature. Useful after manually decomposing requirements into multiple runes."),
			mcp.WithString("runes", mcp.Description("JSON array of rune objects, each with: name (string, required), description (string, required), signature (string, required), behavior (string), positive_tests (string[]), negative_tests (string[])"), mcp.Required()),
		), handleRunesCreateBatch)

		s.AddTool(mcp.NewTool("runes_hydrate",
			mcp.WithDescription("Generate code and tests for a rune via a sandboxed coding agent. The agent reads the rune's English spec and produces implementation files, runs tests, and records coverage. The rune is marked as hydrated on success."),
			mcp.WithString("name", mcp.Description("Full rune name including namespace (e.g. auth/validate-email)"), mcp.Required()),
		), handleRunesHydrate)

		s.AddTool(mcp.NewTool("features_list",
			mcp.WithDescription("List all features in the registry. A feature is a namespace that groups related runes — it describes a domain (e.g. auth, payment) and how its runes compose into components. Returns JSON array of all registered features with their components and connections. Check this before creating runes to see which namespaces exist."),
		), handleFeaturesList)

		s.AddTool(mcp.NewTool("features_create",
			mcp.WithDescription("Create a new feature in the registry. A feature defines a namespace for related runes. The feature name becomes the rune namespace prefix (e.g. feature 'auth' means runes are named auth/validate-email, auth/hash-password, etc.). Pass the complete feature content in the body parameter — description, signature, components with wiring and tests, connections — everything in one call. Do not create a skeleton and edit later. IMPORTANT: Always propose the feature to the user and wait for approval before calling this tool."),
			mcp.WithString("name", mcp.Description("Feature name as a single slug with no slashes (e.g. auth, payment, notifications). This becomes the namespace for its runes."), mcp.Required()),
			mcp.WithString("body", mcp.Description("The full feature content in markdown. Include: description, ## Signature, ## Components (with ### component names, #### Signature, #### Composes, #### Wiring with fenced code blocks, #### Positive tests, #### Negative tests), and ## Connections. This becomes the body of the feature.md file below the # heading."), mcp.Required()),
		), handleFeaturesCreate)

		s.AddTool(mcp.NewTool("features_get",
			mcp.WithDescription("Retrieve a single feature by name. Returns the full feature record including description, components (with their composed runes), connections, version, and status."),
			mcp.WithString("name", mcp.Description("Feature name (e.g. auth, payment)"), mcp.Required()),
		), handleFeaturesGet)

		s.AddTool(mcp.NewTool("features_update",
			mcp.WithDescription("Update a feature's version or status. At least one field must be provided. Use this to promote a feature through the draft → reviewed → stable lifecycle. To modify the feature content (description, components, wiring), edit the feature.md file directly."),
			mcp.WithString("name", mcp.Description("Feature name (e.g. auth, payment)"), mcp.Required()),
			mcp.WithString("version", mcp.Description("New semantic version (e.g. 0.2.0)")),
			mcp.WithString("status", mcp.Description("New status: draft (initial), reviewed (approved), or stable (production-ready)")),
		), handleFeaturesUpdate)

		s.AddTool(mcp.NewTool("features_delete",
			mcp.WithDescription("Delete a feature from the registry. Removes only the feature.md file — runes in the namespace are preserved. Use this to remove a feature definition while keeping its runes intact."),
			mcp.WithString("name", mcp.Description("Feature name (e.g. auth, payment)"), mcp.Required()),
		), handleFeaturesDelete)

		s.AddTool(mcp.NewTool("features_compose",
			mcp.WithDescription("Generate dispatcher and wiring code for a feature by composing its hydrated runes. All runes listed in the feature's components must be hydrated first. The agent reads the feature's wiring pseudocode and rune signatures to generate: (1) a dispatcher that routes all calls by name, (2) wiring functions that implement the pseudocode, (3) registration code, and (4) integration tests. Generated code is stored in runes/<feature>/_composed/."),
			mcp.WithString("name", mcp.Description("Feature name (e.g. auth, payment)"), mcp.Required()),
		), handleFeaturesCompose)

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

func handleFeaturesList(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	features, err := featureStore.List()
	if err != nil {
		return errResult(err), nil
	}
	if features == nil {
		features = []feature.Feature{}
	}
	return textResult(toJSON(features)), nil
}

func handleFeaturesCreate(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	name, _ := args["name"].(string)
	body, _ := args["body"].(string)

	if err := featureStore.Create(name, body); err != nil {
		return errResult(err), nil
	}
	created, err := featureStore.Get(name)
	if err != nil {
		return errResult(err), nil
	}
	return textResult(toJSON(created)), nil
}

func handleFeaturesGet(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	name, _ := args["name"].(string)

	f, err := featureStore.Get(name)
	if err != nil {
		return errResult(err), nil
	}
	return textResult(toJSON(f)), nil
}

func handleFeaturesUpdate(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	name, _ := args["name"].(string)

	f, err := featureStore.Get(name)
	if err != nil {
		return errResult(err), nil
	}

	changed := false
	if ver, ok := args["version"].(string); ok && ver != "" {
		f.Version = ver
		changed = true
	}
	if status, ok := args["status"].(string); ok && status != "" {
		f.Status = status
		changed = true
	}
	if !changed {
		return errResult(fmt.Errorf("at least one of version or status is required")), nil
	}

	if err := featureStore.Update(*f); err != nil {
		return errResult(err), nil
	}
	updated, err := featureStore.Get(name)
	if err != nil {
		return errResult(err), nil
	}
	return textResult(toJSON(updated)), nil
}

func handleFeaturesDelete(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	name, _ := args["name"].(string)

	if err := featureStore.Delete(name); err != nil {
		return errResult(err), nil
	}
	return textResult(fmt.Sprintf("feature %q deleted", name)), nil
}

func handleFeaturesCompose(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	name, _ := args["name"].(string)

	run, err := runner.New(cfg.Agent)
	if err != nil {
		return errResult(err), nil
	}

	result, err := comp.Compose(ctx, name, run, os.Stderr)
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
