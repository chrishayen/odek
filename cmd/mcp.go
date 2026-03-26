package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chrishayen/valkyrie/internal/server/store"
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
			mcp.WithDescription("List runes on the rune server. Filterable by project or namespace prefix. Returns JSON array of rune specs with FQN, description, signature, status."),
			mcp.WithString("project", mcp.Description("Filter by project name")),
			mcp.WithString("namespace", mcp.Description("Filter by namespace prefix (e.g. 'net.http' or 'myapp.auth')")),
		), handleRunesList)

		s.AddTool(mcp.NewTool("runes_get",
			mcp.WithDescription("Get a single rune by its fully-qualified name (dot notation). Returns the full spec including description, signature, behavior, tests, version, and status."),
			mcp.WithString("fqn", mcp.Description("Fully-qualified rune name in dot notation (e.g. net.http.parse_url, myapp.auth.validate_token)"), mcp.Required()),
		), handleRunesGet)

		s.AddTool(mcp.NewTool("runes_search",
			mcp.WithDescription("Search runes by keyword. Matches against FQN, description, and signature. Use this to find existing runes before creating new ones."),
			mcp.WithString("query", mcp.Description("Search query string"), mcp.Required()),
		), handleRunesSearch)

		s.AddTool(mcp.NewTool("runes_approve",
			mcp.WithDescription("Approve a rune spec, changing its status from draft to approved. Use this after the user confirms they're satisfied with the spec."),
			mcp.WithString("fqn", mcp.Description("Fully-qualified rune name to approve"), mcp.Required()),
		), handleRunesApprove)

		s.AddTool(mcp.NewTool("runes_reject",
			mcp.WithDescription("Reject a rune spec with feedback. The feedback is used to guide redesign of the spec."),
			mcp.WithString("fqn", mcp.Description("Fully-qualified rune name to reject"), mcp.Required()),
			mcp.WithString("feedback", mcp.Description("Explanation of what needs to change"), mcp.Required()),
		), handleRunesReject)

		s.AddTool(mcp.NewTool("requirements_submit",
			mcp.WithDescription("Submit refined requirements to the rune server for decomposition. The server breaks requirements into runes, classifies them, searches for existing matches, and designs new specs. Returns a job ID to poll with requirements_status."),
			mcp.WithString("requirements", mcp.Description("The refined requirements text"), mcp.Required()),
		), handleRequirementsSubmit)

		s.AddTool(mcp.NewTool("requirements_status",
			mcp.WithDescription("Check the status of a requirements decomposition job. Returns the current status and, when complete, the list of proposed rune specs and existing matches."),
			mcp.WithString("id", mcp.Description("Job ID returned by requirements_submit"), mcp.Required()),
		), handleRequirementsStatus)

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

func handleRunesList(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	project, _ := args["project"].(string)
	namespace, _ := args["namespace"].(string)

	runes, err := apiClient.RunesList(project, namespace)
	if err != nil {
		return errResult(err), nil
	}
	if runes == nil {
		runes = []store.Rune{}
	}
	return textResult(toJSON(runes)), nil
}

func handleRunesGet(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	fqn, _ := args["fqn"].(string)

	r, err := apiClient.RunesGet(fqn)
	if err != nil {
		return errResult(err), nil
	}
	return textResult(toJSON(r)), nil
}

func handleRunesSearch(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	query, _ := args["query"].(string)

	runes, err := apiClient.RunesSearch(query)
	if err != nil {
		return errResult(err), nil
	}
	if runes == nil {
		runes = []store.Rune{}
	}
	return textResult(toJSON(runes)), nil
}

func handleRunesApprove(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	fqn, _ := args["fqn"].(string)

	r, err := apiClient.RunesApprove(fqn)
	if err != nil {
		return errResult(err), nil
	}
	return textResult(toJSON(r)), nil
}

func handleRunesReject(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	fqn, _ := args["fqn"].(string)
	feedback, _ := args["feedback"].(string)

	// Update the rune status back to draft with feedback in behavior
	patch := map[string]any{
		"status":   "draft",
		"behavior": fmt.Sprintf("[FEEDBACK] %s", feedback),
	}
	r, err := apiClient.RunesUpdate(fqn, patch)
	if err != nil {
		return errResult(err), nil
	}
	return textResult(toJSON(r)), nil
}

func handleRequirementsSubmit(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	requirements, _ := args["requirements"].(string)

	result, err := apiClient.RequirementsSubmit(cfg.Project, requirements)
	if err != nil {
		return errResult(err), nil
	}
	return textResult(toJSON(result)), nil
}

func handleRequirementsStatus(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	id, _ := args["id"].(string)

	result, err := apiClient.RequirementsStatus(id)
	if err != nil {
		return errResult(err), nil
	}
	return textResult(toJSON(result)), nil
}
