package openai

import "testing"

func TestNewClientNormalizesBaseURL(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{name: "default", in: "", want: "http://127.0.0.1:1234/v1"},
		{name: "host without version", in: "http://localhost:8080", want: "http://localhost:8080/v1"},
		{name: "trailing slash", in: "http://localhost:8080/", want: "http://localhost:8080/v1"},
		{name: "already versioned", in: "https://api.openai.com/v1", want: "https://api.openai.com/v1"},
		{name: "already versioned with slash", in: "https://api.openai.com/v1/", want: "https://api.openai.com/v1"},
		{name: "scheme omitted", in: "localhost:8080", want: "http://localhost:8080/v1"},
		{name: "nested API path", in: "http://localhost:8080/api", want: "http://localhost:8080/api/v1"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			client, err := NewClient(tc.in)
			if err != nil {
				t.Fatalf("NewClient(%q): %v", tc.in, err)
			}
			if client.baseURL != tc.want {
				t.Fatalf("baseURL = %q, want %q", client.baseURL, tc.want)
			}
		})
	}
}

func TestNewSystemPromptBuilderUsesBaseAndSortedMetadata(t *testing.T) {
	msg := NewSystemPromptBuilder("base").
		AddMetadata("z", "last").
		AddMetadata("a", "first").
		Build()

	if msg.Role != RoleSystem {
		t.Fatalf("Role = %q, want %q", msg.Role, RoleSystem)
	}
	want := "base | a=first, z=last"
	if msg.Content != want {
		t.Fatalf("Content = %q, want %q", msg.Content, want)
	}
}
