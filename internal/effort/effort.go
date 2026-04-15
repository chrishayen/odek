package effort

import (
	"context"
	"encoding/json"
	"fmt"

	openai "shotgun.dev/odek/openai"
)

// Result is the model's complexity estimate for a software requirement.
type Result struct {
	Level  int    `json:"level"`
	Reason string `json:"reason"`
}

const systemPrompt = "You are a software-complexity estimator. Given a software requirement, rate it 1-5 by calling the rate_effort tool. Reply only via the tool call."

var rateEffortTool = openai.Tool{
	Type: openai.ToolTypeFunction,
	Function: &openai.FunctionDefinition{
		Name:        "rate_effort",
		Description: "Rate the complexity of a software requirement on a 1-5 scale.",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"level": map[string]any{
					"type":        "integer",
					"minimum":     1,
					"maximum":     5,
					"description": "1=trivial (hello world, single function); 2=small (one file or simple CLI); 3=medium (a few modules); 4=large (subsystem with several integration points); 5=very large (full application stack)",
				},
				"reason": map[string]any{
					"type":        "string",
					"description": "One short sentence justifying the level.",
				},
			},
			"required": []string{"level", "reason"},
		},
	},
}

// Estimate asks the model to rate the requirement's complexity 1-5 via a forced tool call.
func Estimate(ctx context.Context, client *openai.Client, requirement string) (Result, error) {
	args, err := client.AskTool(ctx, systemPrompt,
		"Rate the complexity of this requirement: "+requirement, rateEffortTool)
	if err != nil {
		return Result{}, fmt.Errorf("effort completion failed: %w", err)
	}
	var est Result
	if err := json.Unmarshal([]byte(args), &est); err != nil {
		return Result{}, fmt.Errorf("parsing effort args: %w (raw: %s)", err, args)
	}
	if est.Level < 1 || est.Level > 5 {
		return Result{}, fmt.Errorf("level out of range: %d", est.Level)
	}
	return est, nil
}
