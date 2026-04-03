package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrishayen/odek/internal/app"
	"github.com/chrishayen/odek/internal/decomposer"
	"github.com/chrishayen/odek/internal/feature"
	runepkg "github.com/chrishayen/odek/internal/rune"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start MCP server (stdio transport)",
	RunE: func(cmd *cobra.Command, args []string) error {
		s := server.NewMCPServer("odek", "0.2.0", server.WithToolCapabilities(true))

		s.AddTool(mcp.NewTool("runes_list",
			mcp.WithDescription("List all runes in the registry."),
		), handleRunesList)

		s.AddTool(mcp.NewTool("runes_create",
			mcp.WithDescription("Create a new rune. Names use dot-separated paths with snake_case segments (e.g. auth.validate_email)."),
			mcp.WithString("name", mcp.Description("Rune name as dot path"), mcp.Required()),
			mcp.WithString("description", mcp.Description("One or two sentences stating what the function does"), mcp.Required()),
			mcp.WithString("signature", mcp.Description("Function signature using precise types"), mcp.Required()),
			mcp.WithString("behavior", mcp.Description("Precise description of inputs, outputs, edge cases, and constraints")),
		), handleRunesCreate)

		s.AddTool(mcp.NewTool("runes_get",
			mcp.WithDescription("Retrieve a single rune by name."),
			mcp.WithString("name", mcp.Description("Full rune name as dot path"), mcp.Required()),
		), handleRunesGet)

		s.AddTool(mcp.NewTool("runes_update",
			mcp.WithDescription("Update an existing rune's description, signature, or version."),
			mcp.WithString("name", mcp.Description("Full rune name as dot path"), mcp.Required()),
			mcp.WithString("description", mcp.Description("New description")),
			mcp.WithString("signature", mcp.Description("New function signature")),
			mcp.WithString("version", mcp.Description("New semantic version (e.g. 1.1.0)")),
		), handleRunesUpdate)

		s.AddTool(mcp.NewTool("runes_delete",
			mcp.WithDescription("Delete a rune from the registry."),
			mcp.WithString("name", mcp.Description("Full rune name as dot path"), mcp.Required()),
		), handleRunesDelete)

		s.AddTool(mcp.NewTool("runes_decompose",
			mcp.WithDescription("Decompose plain-text requirements into a composition tree of runes. Uses stdlib-first strategy."),
			mcp.WithString("requirements", mcp.Description("Plain-text English description of what you need built."), mcp.Required()),
		), handleRunesDecompose)

		s.AddTool(mcp.NewTool("runes_create_batch",
			mcp.WithDescription("Create multiple runes from a composition tree. Uses the same indented tree format as decompose output."),
			mcp.WithString("tree", mcp.Description("Composition tree text with dot-path names, @ signatures, + positive tests, - negative tests"), mcp.Required()),
		), handleRunesCreateBatch)

		s.AddTool(mcp.NewTool("runes_hydrate",
			mcp.WithDescription("Generate code and tests for a single rune."),
			mcp.WithString("name", mcp.Description("Full rune name as dot path"), mcp.Required()),
		), handleRunesHydrate)

		s.AddTool(mcp.NewTool("runes_hydration_spec",
			mcp.WithDescription("Get the hydration prompt for a rune. Used for sub-agent workflows."),
			mcp.WithString("name", mcp.Description("Full rune name as dot path"), mcp.Required()),
		), handleRunesHydrationSpec)

		s.AddTool(mcp.NewTool("runes_finalize_hydration",
			mcp.WithDescription("Submit generated code for a rune, extract files, run tests, and mark it as hydrated."),
			mcp.WithString("name", mcp.Description("Full rune name as dot path"), mcp.Required()),
			mcp.WithString("output", mcp.Description("The sub-agent's complete text output containing === FILE: === / === END FILE === blocks"), mcp.Required()),
		), handleRunesFinalizeHydration)

		s.AddTool(mcp.NewTool("runes_check",
			mcp.WithDescription("Check for stale references in rune dependencies."),
		), handleRunesCheck)

		s.AddTool(mcp.NewTool("runes_verify",
			mcp.WithDescription("Verify all hydrated runes against their specs."),
		), handleRunesVerify)

		s.AddTool(mcp.NewTool("features_list",
			mcp.WithDescription("List all features in the registry."),
		), handleFeaturesList)

		s.AddTool(mcp.NewTool("features_create",
			mcp.WithDescription("Create a new feature in the registry."),
			mcp.WithString("name", mcp.Description("Feature name as a single slug"), mcp.Required()),
			mcp.WithString("body", mcp.Description("The full feature content in markdown"), mcp.Required()),
		), handleFeaturesCreate)

		s.AddTool(mcp.NewTool("features_get",
			mcp.WithDescription("Retrieve a single feature by name."),
			mcp.WithString("name", mcp.Description("Feature name"), mcp.Required()),
		), handleFeaturesGet)

		s.AddTool(mcp.NewTool("features_update",
			mcp.WithDescription("Update a feature's version or status."),
			mcp.WithString("name", mcp.Description("Feature name"), mcp.Required()),
			mcp.WithString("version", mcp.Description("New semantic version")),
			mcp.WithString("status", mcp.Description("New status: draft, reviewed, or stable")),
		), handleFeaturesUpdate)

		s.AddTool(mcp.NewTool("features_delete",
			mcp.WithDescription("Delete a feature from the registry."),
			mcp.WithString("name", mcp.Description("Feature name"), mcp.Required()),
		), handleFeaturesDelete)

		s.AddTool(mcp.NewTool("features_compose",
			mcp.WithDescription("Generate dispatcher and wiring code for a feature."),
			mcp.WithString("name", mcp.Description("Feature name"), mcp.Required()),
		), handleFeaturesCompose)

		s.AddTool(mcp.NewTool("apps_list",
			mcp.WithDescription("List all apps in the registry."),
		), handleAppsList)

		s.AddTool(mcp.NewTool("apps_create",
			mcp.WithDescription("Create a new app in the registry."),
			mcp.WithString("name", mcp.Description("App name as a single slug"), mcp.Required()),
			mcp.WithString("body", mcp.Description("The full app content in markdown"), mcp.Required()),
		), handleAppsCreate)

		s.AddTool(mcp.NewTool("apps_get",
			mcp.WithDescription("Retrieve a single app by name."),
			mcp.WithString("name", mcp.Description("App name"), mcp.Required()),
		), handleAppsGet)

		s.AddTool(mcp.NewTool("apps_update",
			mcp.WithDescription("Update an app's version, status, or entry point."),
			mcp.WithString("name", mcp.Description("App name"), mcp.Required()),
			mcp.WithString("version", mcp.Description("New semantic version")),
			mcp.WithString("status", mcp.Description("New status: draft, reviewed, or stable")),
			mcp.WithString("entry_point", mcp.Description("Entry point feature name")),
		), handleAppsUpdate)

		s.AddTool(mcp.NewTool("apps_delete",
			mcp.WithDescription("Delete an app from the registry."),
			mcp.WithString("name", mcp.Description("App name"), mcp.Required()),
		), handleAppsDelete)

		return server.ServeStdio(s)
	},
}

func toJSON(v any) string {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("{\"error\": %q}", err.Error())
	}
	return string(data)
}

func textResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{Content: []mcp.Content{mcp.NewTextContent(text)}}
}

func errResult(err error) *mcp.CallToolResult {
	return &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.NewTextContent(err.Error())}}
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
	name, _ := req.GetArguments()["name"].(string)
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
		r.Version = runepkg.ParseSemver(ver)
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
	name, _ := req.GetArguments()["name"].(string)
	if err := store.Delete(name); err != nil {
		return errResult(err), nil
	}
	return textResult(fmt.Sprintf("rune %q deleted", name)), nil
}

func handleRunesDecompose(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	requirements, _ := req.GetArguments()["requirements"].(string)
	result, err := dec.Decompose(ctx, requirements, "", nil)
	if err != nil {
		return errResult(err), nil
	}
	var created []string
	for _, p := range result.NewRunes {
		r := p.ToRune()
		if err := store.Create(r); err != nil {
			created = append(created, fmt.Sprintf("FAILED %s: %v", p.Name, err))
			continue
		}
		created = append(created, fmt.Sprintf("created %s", p.Name))
	}
	type decomposeOutput struct {
		*decomposer.Result
		Created []string `json:"created"`
	}
	return textResult(toJSON(decomposeOutput{Result: result, Created: created})), nil
}

func handleRunesHydrate(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, _ := req.GetArguments()["name"].(string)
	result, err := hyd.Hydrate(ctx, name)
	if err != nil {
		return errResult(err), nil
	}
	return textResult(toJSON(result)), nil
}

func handleRunesHydrationSpec(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, _ := req.GetArguments()["name"].(string)
	spec, err := hyd.GetHydrationSpec(name)
	if err != nil {
		return errResult(err), nil
	}
	return textResult(toJSON(spec)), nil
}

func handleRunesFinalizeHydration(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	name, _ := args["name"].(string)
	output, _ := args["output"].(string)
	result, err := hyd.FinalizeHydration(name, output)
	if err != nil {
		return errResult(err), nil
	}
	return textResult(toJSON(result)), nil
}

func handleRunesCheck(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	stale, ok, err := store.CheckStaleRefs()
	if err != nil {
		return errResult(err), nil
	}
	if stale == 0 {
		return textResult(fmt.Sprintf("All %d references up to date.", ok)), nil
	}
	return textResult(fmt.Sprintf("%d stale, %d ok", stale, ok)), nil
}

func handleRunesVerify(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result, err := hyd.VerifyAll(ctx, cfg.Concurrency, nil)
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
	name, _ := req.GetArguments()["name"].(string)
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
	name, _ := req.GetArguments()["name"].(string)
	if err := featureStore.Delete(name); err != nil {
		return errResult(err), nil
	}
	return textResult(fmt.Sprintf("feature %q deleted", name)), nil
}

func handleFeaturesCompose(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, _ := req.GetArguments()["name"].(string)
	result, err := comp.Compose(ctx, name)
	if err != nil {
		return errResult(err), nil
	}
	return textResult(toJSON(result)), nil
}

func handleAppsList(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	apps, err := appStore.List()
	if err != nil {
		return errResult(err), nil
	}
	if apps == nil {
		apps = []app.App{}
	}
	return textResult(toJSON(apps)), nil
}

func handleAppsCreate(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	name, _ := args["name"].(string)
	body, _ := args["body"].(string)
	if err := appStore.Create(name, body); err != nil {
		return errResult(err), nil
	}
	created, err := appStore.Get(name)
	if err != nil {
		return errResult(err), nil
	}
	return textResult(toJSON(created)), nil
}

func handleAppsGet(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, _ := req.GetArguments()["name"].(string)
	a, err := appStore.Get(name)
	if err != nil {
		return errResult(err), nil
	}
	return textResult(toJSON(a)), nil
}

func handleAppsUpdate(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	name, _ := args["name"].(string)
	a, err := appStore.Get(name)
	if err != nil {
		return errResult(err), nil
	}
	changed := false
	if ver, ok := args["version"].(string); ok && ver != "" {
		a.Version = ver
		changed = true
	}
	if status, ok := args["status"].(string); ok && status != "" {
		a.Status = status
		changed = true
	}
	if ep, ok := args["entry_point"].(string); ok && ep != "" {
		a.EntryPoint = ep
		changed = true
	}
	if !changed {
		return errResult(fmt.Errorf("at least one of version, status, or entry_point is required")), nil
	}
	if err := appStore.Update(*a); err != nil {
		return errResult(err), nil
	}
	updated, err := appStore.Get(name)
	if err != nil {
		return errResult(err), nil
	}
	return textResult(toJSON(updated)), nil
}

func handleAppsDelete(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, _ := req.GetArguments()["name"].(string)
	if err := appStore.Delete(name); err != nil {
		return errResult(err), nil
	}
	return textResult(fmt.Sprintf("app %q deleted", name)), nil
}

func handleRunesCreateBatch(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tree, _ := req.GetArguments()["tree"].(string)
	nodes := runepkg.ParseTree(tree)
	var results []string
	for _, n := range nodes {
		// Skip reference-only nodes (-> refs with no own signature/tests)
		if len(n.Refs) > 0 && n.Signature == "" && len(n.Pos) == 0 && len(n.Neg) == 0 {
			continue
		}
		r := runepkg.Rune{
			Name:          n.Path,
			Signature:     n.Signature,
			PositiveTests: n.Pos,
			NegativeTests: n.Neg,
			Dependencies:  n.Refs,
		}
		if len(n.Pos) > 0 {
			r.Description = n.Pos[0]
		}
		if err := store.Create(r); err != nil {
			results = append(results, fmt.Sprintf("FAILED %s: %v", r.Name, err))
			continue
		}
		results = append(results, fmt.Sprintf("created %s", r.Name))
	}
	return textResult(toJSON(results)), nil
}
