package e2e_test

import (
	"strings"
	"testing"
)

func TestCLICreateRune(t *testing.T) {
	dir, cleanup := testEnv(t, "")
	defer cleanup()

	out, code := run(t, dir, "runes", "create", "--name", "test/user-auth", "--description", "Accepts a username and password", "--signature", "(username: string, password: string) -> result[bool, string]")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, out)
	}
	if !strings.Contains(out, "test/user-auth") {
		t.Errorf("expected name in output, got: %s", out)
	}
	if !strings.Contains(out, "0.1.0") {
		t.Errorf("expected default version, got: %s", out)
	}
}

func TestCLICreateRuneDuplicate(t *testing.T) {
	dir, cleanup := testEnv(t, "")
	defer cleanup()

	run(t, dir, "runes", "create", "--name", "test/dup", "--description", "first", "--signature", "(x: i32) -> bool")
	_, code := run(t, dir, "runes", "create", "--name", "test/dup", "--description", "second", "--signature", "(x: i32) -> bool")
	if code == 0 {
		t.Error("expected non-zero exit for duplicate rune")
	}
}

func TestCLIListRunes(t *testing.T) {
	dir, cleanup := testEnv(t, "")
	defer cleanup()

	run(t, dir, "runes", "create", "--name", "test/rune-a", "--description", "first", "--signature", "(x: i32) -> bool")
	run(t, dir, "runes", "create", "--name", "test/rune-b", "--description", "second", "--signature", "(y: string) -> i64")

	out, code := run(t, dir, "runes", "list")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, out)
	}
	if !strings.Contains(out, "test/rune-a") || !strings.Contains(out, "test/rune-b") {
		t.Errorf("expected both runes in output, got: %s", out)
	}
}

func TestCLIGetRune(t *testing.T) {
	dir, cleanup := testEnv(t, "")
	defer cleanup()

	run(t, dir, "runes", "create", "--name", "test/my-rune", "--description", "test rune", "--signature", "(x: i32) -> bool")

	out, code := run(t, dir, "runes", "get", "test/my-rune")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, out)
	}
	if !strings.Contains(out, "test/my-rune") {
		t.Errorf("expected rune name in output, got: %s", out)
	}
}

func TestCLIGetRuneNotFound(t *testing.T) {
	dir, cleanup := testEnv(t, "")
	defer cleanup()

	_, code := run(t, dir, "runes", "get", "does-not-exist")
	if code == 0 {
		t.Error("expected non-zero exit for missing rune")
	}
}

func TestCLIUpdateRune(t *testing.T) {
	dir, cleanup := testEnv(t, "")
	defer cleanup()

	run(t, dir, "runes", "create", "--name", "test/update-me", "--description", "original", "--signature", "(x: i32) -> bool")

	out, code := run(t, dir, "runes", "update", "test/update-me", "--description", "updated description", "--version", "0.2.0")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, out)
	}
	if !strings.Contains(out, "updated description") {
		t.Errorf("expected updated description in output, got: %s", out)
	}
	if !strings.Contains(out, "0.2.0") {
		t.Errorf("expected updated version in output, got: %s", out)
	}
}

func TestCLIDeleteRune(t *testing.T) {
	dir, cleanup := testEnv(t, "")
	defer cleanup()

	run(t, dir, "runes", "create", "--name", "test/delete-me", "--description", "gone soon", "--signature", "(x: i32) -> bool")

	out, code := run(t, dir, "runes", "delete", "test/delete-me")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, out)
	}
	if !strings.Contains(out, "deleted") {
		t.Errorf("expected deleted confirmation, got: %s", out)
	}

	_, code = run(t, dir, "runes", "get", "test/delete-me")
	if code == 0 {
		t.Error("expected non-zero exit after delete")
	}
}
