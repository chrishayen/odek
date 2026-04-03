package e2e_test

import (
	"strings"
	"testing"
)

func TestCLICreateFeature(t *testing.T) {
	dir, cleanup := testEnv(t, "")
	defer cleanup()

	out, code := run(t, dir, "features", "create", "--name", "auth", "--description", "Authentication feature")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, out)
	}
	if !strings.Contains(out, "auth") {
		t.Errorf("expected name in output, got: %s", out)
	}
	if !strings.Contains(out, "0.1.0") {
		t.Errorf("expected default version, got: %s", out)
	}
}

func TestCLICreateFeatureDuplicate(t *testing.T) {
	dir, cleanup := testEnv(t, "")
	defer cleanup()

	run(t, dir, "features", "create", "--name", "auth", "--description", "first")
	_, code := run(t, dir, "features", "create", "--name", "auth", "--description", "second")
	if code == 0 {
		t.Error("expected non-zero exit for duplicate feature")
	}
}

func TestCLIListFeatures(t *testing.T) {
	dir, cleanup := testEnv(t, "")
	defer cleanup()

	run(t, dir, "features", "create", "--name", "auth", "--description", "Authentication")
	run(t, dir, "features", "create", "--name", "payment", "--description", "Payment processing")

	out, code := run(t, dir, "features", "list")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, out)
	}
	if !strings.Contains(out, "auth") || !strings.Contains(out, "payment") {
		t.Errorf("expected both features in output, got: %s", out)
	}
}

func TestCLIGetFeature(t *testing.T) {
	dir, cleanup := testEnv(t, "")
	defer cleanup()

	run(t, dir, "features", "create", "--name", "auth", "--description", "Authentication")

	out, code := run(t, dir, "features", "get", "auth")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, out)
	}
	if !strings.Contains(out, "auth") {
		t.Errorf("expected feature name in output, got: %s", out)
	}
}

func TestCLIGetFeatureNotFound(t *testing.T) {
	dir, cleanup := testEnv(t, "")
	defer cleanup()

	_, code := run(t, dir, "features", "get", "nosuchfeature")
	if code == 0 {
		t.Error("expected non-zero exit for missing feature")
	}
}

func TestCLIUpdateFeature(t *testing.T) {
	dir, cleanup := testEnv(t, "")
	defer cleanup()

	run(t, dir, "features", "create", "--name", "auth", "--description", "Authentication")

	out, code := run(t, dir, "features", "update", "auth", "--version", "1.0.0", "--status", "stable")
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

func TestCLIDeleteFeature(t *testing.T) {
	dir, cleanup := testEnv(t, "")
	defer cleanup()

	run(t, dir, "features", "create", "--name", "auth", "--description", "Authentication")

	out, code := run(t, dir, "features", "delete", "auth")
	if code != 0 {
		t.Fatalf("expected exit 0, got %d: %s", code, out)
	}
	if !strings.Contains(out, "deleted") {
		t.Errorf("expected deleted confirmation, got: %s", out)
	}

	_, code = run(t, dir, "features", "get", "auth")
	if code == 0 {
		t.Error("expected non-zero exit after delete")
	}
}
