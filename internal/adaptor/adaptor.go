package adaptor

import "context"

// Adaptor runs a prompt through an AI model and returns the response text.
type Adaptor interface {
	Run(ctx context.Context, prompt string) (string, error)
}
