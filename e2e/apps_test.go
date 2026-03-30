package e2e_test

import (
	"strings"
	"testing"
)

func TestCLICreateApp(t *testing.T) {
	dir, cleanup := testEnv(t, "")
	defer cleanup()

	out, code := run(t, dir, "apps", "create", "--name", "myapp", "--description", "A test application")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, out)
	}
	if !strings.Contains(out, "myapp") {
		t.Errorf("expected name in output, got: %s", out)
	}
	if !strings.Contains(out, "0.1.0") {
		t.Errorf("expected default version, got: %s", out)
	}
}

func TestCLICreateAppDuplicate(t *testing.T) {
	dir, cleanup := testEnv(t, "")
	defer cleanup()

	run(t, dir, "apps", "create", "--name", "dupapp", "--description", "first")
	_, code := run(t, dir, "apps", "create", "--name", "dupapp", "--description", "second")
	if code == 0 {
		t.Error("expected non-zero exit for duplicate app")
	}
}

func TestCLIListApps(t *testing.T) {
	dir, cleanup := testEnv(t, "")
	defer cleanup()

	run(t, dir, "apps", "create", "--name", "app1", "--description", "first app")
	run(t, dir, "apps", "create", "--name", "app2", "--description", "second app")

	out, code := run(t, dir, "apps", "list")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, out)
	}
	if !strings.Contains(out, "app1") || !strings.Contains(out, "app2") {
		t.Errorf("expected both apps in output, got: %s", out)
	}
}

func TestCLIGetApp(t *testing.T) {
	dir, cleanup := testEnv(t, "")
	defer cleanup()

	run(t, dir, "apps", "create", "--name", "getme", "--description", "test app")

	out, code := run(t, dir, "apps", "get", "getme")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, out)
	}
	if !strings.Contains(out, "getme") {
		t.Errorf("expected app name in output, got: %s", out)
	}
}

func TestCLIGetAppNotFound(t *testing.T) {
	dir, cleanup := testEnv(t, "")
	defer cleanup()

	_, code := run(t, dir, "apps", "get", "nosuchapp")
	if code == 0 {
		t.Error("expected non-zero exit for missing app")
	}
}

func TestCLIUpdateApp(t *testing.T) {
	dir, cleanup := testEnv(t, "")
	defer cleanup()

	run(t, dir, "apps", "create", "--name", "updateme", "--description", "original")

	out, code := run(t, dir, "apps", "update", "updateme", "--version", "1.0.0", "--status", "stable")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, out)
	}
	if !strings.Contains(out, "1.0.0") {
		t.Errorf("expected updated version in output, got: %s", out)
	}
	if !strings.Contains(out, "stable") {
		t.Errorf("expected updated status in output, got: %s", out)
	}
}

func TestCLIDeleteApp(t *testing.T) {
	dir, cleanup := testEnv(t, "")
	defer cleanup()

	run(t, dir, "apps", "create", "--name", "deleteme", "--description", "gone soon")

	out, code := run(t, dir, "apps", "delete", "deleteme")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, out)
	}
	if !strings.Contains(out, "deleted") {
		t.Errorf("expected deleted confirmation, got: %s", out)
	}

	_, code = run(t, dir, "apps", "get", "deleteme")
	if code == 0 {
		t.Error("expected non-zero exit after delete")
	}
}
