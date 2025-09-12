package adr

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	reADRTitle = regexp.MustCompile(`(?i)^#\s*ADR\s+(\d+)\s*:\s*(.+)$`)
	reStatus   = regexp.MustCompile(`(?i)^##\s*Status\s*$`)
	reStatusKV = regexp.MustCompile(`(?i)^(\*\*Status:\*\*|[-*]\s*Status:?|\s*Status:)\s*(.+)$`)
	reDateKV   = regexp.MustCompile(`(?i)^(Date|Date\s*:\s*)\s*:?[\s]*([0-9]{4}-[0-9]{2}-[0-9]{2}).*$`)
)

type Meta struct {
	Number int
	Title  string
	Status string
	Date   string // YYYY-MM-DD
}

type Frontmatter struct {
	ID     any    `yaml:"id"`
	Title  string `yaml:"title"`
	Status string `yaml:"status"`
	Date   string `yaml:"date"`
}

// parseFrontmatter extracts YAML frontmatter from file content
func parseFrontmatter(content []byte) (*Frontmatter, []byte, error) {
	if !bytes.HasPrefix(content, []byte("---\n")) {
		return nil, content, nil
	}

	// Find the end of frontmatter
	end := bytes.Index(content[4:], []byte("\n---\n"))
	if end == -1 {
		return nil, content, nil
	}

	yamlContent := content[4 : end+4]
	remaining := content[end+9:] // Skip past the closing ---

	var fm Frontmatter
	if err := yaml.Unmarshal(yamlContent, &fm); err != nil {
		return nil, content, err
	}

	return &fm, remaining, nil
}

// ParseADR parses minimal metadata from an ADR file.
func ParseADR(path string) (Meta, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return Meta{}, err
	}

	var m Meta

	// Try to parse frontmatter first
	if fm, remaining, err := parseFrontmatter(content); err == nil && fm != nil {
		// Convert frontmatter to Meta
		if fm.Title != "" {
			m.Title = fm.Title
		}
		if fm.Status != "" {
			m.Status = fm.Status
		}
		if fm.Date != "" {
			m.Date = fm.Date
		}
		// Handle ID field which can be int or string
		if fm.ID != nil {
			switch id := fm.ID.(type) {
			case int:
				m.Number = id
			case string:
				m.Number = atoi(id)
			case float64: // YAML can decode numbers as float64
				m.Number = int(id)
			}
		}

		// If we have all required fields from frontmatter, use them
		if m.Title != "" && m.Status != "" && m.Date != "" && m.Number != 0 {
			return m, nil
		}

		// Otherwise, fall back to parsing the remaining content
		content = remaining
	}

	// Fallback to original parsing logic for backward compatibility
	s := bufio.NewScanner(bytes.NewReader(content))
	s.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)

	var sawStatusHeader bool

	for s.Scan() {
		line := s.Text()

		if m.Number == 0 || m.Title == "" {
			if g := reADRTitle.FindStringSubmatch(line); len(g) == 3 {
				if m.Number == 0 {
					m.Number = atoi(g[1])
				}
				if m.Title == "" {
					m.Title = strings.TrimSpace(g[2])
				}
				continue
			}
		}
		if m.Status == "" {
			if g := reStatusKV.FindStringSubmatch(line); len(g) == 3 {
				m.Status = strings.TrimSpace(g[2])
			}
			if reStatus.MatchString(line) {
				sawStatusHeader = true
				continue
			}
			if sawStatusHeader {
				// The first non-empty line after a "## Status" header is the status
				trim := strings.TrimSpace(line)
				if trim != "" {
					m.Status = trim
					sawStatusHeader = false
				}
			}
		}
		if m.Date == "" {
			if g := reDateKV.FindStringSubmatch(line); len(g) == 3 {
				m.Date = strings.TrimSpace(g[2])
			}
		}
	}

	if m.Date == "" {
		// fall back to file mod time
		if fi, err := os.Stat(path); err == nil {
			m.Date = fi.ModTime().Format("2006-01-02")
		}
	}

	if m.Title == "" {
		// fallback: derive title from filename
		base := filepath.Base(path)
		base = strings.TrimSuffix(base, filepath.Ext(base))
		parts := strings.SplitN(base, "-", 2)
		if len(parts) == 2 {
			m.Title = strings.ReplaceAll(parts[1], "-", " ")
		}
	}

	return m, nil
}

func atoi(s string) int {
	var n int
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		n = n*10 + int(c-'0')
	}
	return n
}
