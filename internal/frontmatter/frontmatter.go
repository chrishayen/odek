package frontmatter

import (
	"strings"

	"gopkg.in/yaml.v3"
)

// Parse extracts YAML frontmatter from content delimited by "---\n"
// and unmarshals it into dest. Returns the body after the closing delimiter.
func Parse(content string, dest any) string {
	if !strings.HasPrefix(content, "---\n") {
		return content
	}
	end := strings.Index(content[4:], "\n---")
	if end == -1 {
		return content
	}
	fm := content[4 : 4+end]
	_ = yaml.Unmarshal([]byte(fm), dest)
	return content[4+end+4:]
}

// Strip returns everything after the closing "---" delimiter,
// discarding the frontmatter block.
func Strip(content string) string {
	if !strings.HasPrefix(content, "---\n") {
		return content
	}
	end := strings.Index(content[4:], "\n---")
	if end == -1 {
		return content
	}
	return content[4+end+4:]
}
