package claude

import "testing"

func TestStripCodeFences(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"go fence", "```go\nfunc main() {}\n```", "func main() {}"},
		{"ts fence", "```typescript\nconst x = 1\n```", "const x = 1"},
		{"plain fence", "```\nhello\n```", "hello"},
		{"no fence", "just text", "just text"},
		{"multiple fences", "```go\nfunc a() {}\n```\n```go\nfunc b() {}\n```", "func a() {}\nfunc b() {}"},
	}
	for _, tt := range tests {
		got := StripCodeFences(tt.input)
		if got != tt.want {
			t.Errorf("%s: StripCodeFences() = %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestMockResponse(t *testing.T) {
	c := New("", "", true)

	// Decompose
	resp, err := c.Call("decompose composition tree", "test requirements")
	if err != nil {
		t.Fatal(err)
	}
	if resp == "" {
		t.Error("expected non-empty decompose response")
	}

	// Verify
	resp, err = c.Call("RESULT: ALL PASS verify", "check it")
	if err != nil {
		t.Fatal(err)
	}
	if resp == "" {
		t.Error("expected non-empty verify response")
	}

	// Hydrate (default)
	resp, err = c.Call("", "generate code for test.hello")
	if err != nil {
		t.Fatal(err)
	}
	if resp == "" {
		t.Error("expected non-empty hydrate response")
	}
}

func TestClassifyError(t *testing.T) {
	err := classifyError(401, []byte("unauthorized"))
	if err == nil {
		t.Fatal("expected error")
	}
	if got := err.Error(); got != "auth error: token expired — run 'odek login'" {
		t.Errorf("got %q", got)
	}

	err = classifyError(500, []byte("internal server error"))
	if err == nil {
		t.Fatal("expected error")
	}
	if got := err.Error(); got != "api error 500: internal server error" {
		t.Errorf("got %q", got)
	}
}
