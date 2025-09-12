package adr

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

//go:embed templates/madr.md templates/nygard.md
var builtinTemplates embed.FS

// Manager holds settings for ADR operations.
type Manager struct {
	Dir string
}

// NewOptions controls ADR creation.
type NewOptions struct {
	Template string // "madr" | "nygard" | "/path/to/template.md"
	Status   string // default: Proposed
	Date     string // ISO date; default today
}

func EnsureDir(dir string) error {
	if dir == "" {
		return errors.New("empty dir")
	}
	return os.MkdirAll(dir, 0o755)
}

func (m Manager) nextID() (int, error) {
	entries, err := os.ReadDir(m.Dir)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return 0, err
	}
	max := 0
	for _, e := range entries {
		name := e.Name()
		if !e.IsDir() && len(name) >= 4 {
			if n, ok := parseLeadingNumber(name); ok && n > max {
				max = n
			}
		}
	}
	return max + 1, nil
}

func parseLeadingNumber(name string) (int, bool) {
	var n int
	for i := 0; i < len(name) && i < 6; i++ { // up to 6 digits
		if name[i] < '0' || name[i] > '9' {
			break
		}
		n = n*10 + int(name[i]-'0')
	}
	if n == 0 {
		return 0, false
	}
	return n, true
}

func sanitizeTitle(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.ReplaceAll(s, "_", "-")
	s = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, s)
	return s
}

func (m Manager) WriteNewADR(title string, opt NewOptions) (string, error) {
	if err := EnsureDir(m.Dir); err != nil {
		return "", err
	}
	id, err := m.nextID()
	if err != nil {
		return "", err
	}

	file := fmt.Sprintf("%04d-%s.md", id, sanitizeTitle(title))
	path := filepath.Join(m.Dir, file)

	if opt.Date == "" {
		opt.Date = time.Now().Format("2006-01-02")
	}
	if opt.Status == "" {
		opt.Status = "Proposed"
	}

	tpl, err := m.loadTemplate(opt.Template)
	if err != nil {
		return "", err
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		return "", err
	}
	defer f.Close()

	data := map[string]any{
		"ID":     fmt.Sprintf("%04d", id),
		"Title":  title,
		"Status": opt.Status,
		"Date":   opt.Date,
	}

	if err := tpl.Execute(f, data); err != nil {
		return "", err
	}
	return path, nil
}

func (m Manager) loadTemplate(name string) (*template.Template, error) {
	var content []byte
	var err error

	switch strings.ToLower(strings.TrimSpace(name)) {
	case "", "madr":
		content, err = builtinTemplates.ReadFile("templates/madr.md")
	case "nygard":
		content, err = builtinTemplates.ReadFile("templates/nygard.md")
	default:
		content, err = os.ReadFile(name)
	}
	if err != nil {
		return nil, err
	}
	return template.New("adr").Parse(string(content))
}
