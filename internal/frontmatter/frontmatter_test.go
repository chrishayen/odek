package frontmatter

import "testing"

type testMeta struct {
	Version string `yaml:"version"`
	Status  string `yaml:"status"`
}

func TestParse(t *testing.T) {
	content := "---\nversion: 1.0.0\nstatus: draft\n---\n\n# Title\n\nBody text\n"
	var m testMeta
	body := Parse(content, &m)

	if m.Version != "1.0.0" {
		t.Errorf("Version = %q", m.Version)
	}
	if m.Status != "draft" {
		t.Errorf("Status = %q", m.Status)
	}
	if body != "\n\n# Title\n\nBody text\n" {
		t.Errorf("body = %q", body)
	}
}

func TestParseNoFrontmatter(t *testing.T) {
	content := "Just plain text"
	var m testMeta
	body := Parse(content, &m)

	if body != content {
		t.Errorf("body = %q, want %q", body, content)
	}
	if m.Version != "" {
		t.Errorf("Version should be empty, got %q", m.Version)
	}
}

func TestStrip(t *testing.T) {
	content := "---\nversion: 1.0.0\n---\n\nBody\n"
	body := Strip(content)
	if body != "\n\nBody\n" {
		t.Errorf("Strip() = %q", body)
	}
}

func TestStripNoFrontmatter(t *testing.T) {
	content := "No frontmatter here"
	body := Strip(content)
	if body != content {
		t.Errorf("Strip() = %q, want %q", body, content)
	}
}
