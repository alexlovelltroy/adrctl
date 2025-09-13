package adr

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

//go:embed templates/index.md
var indexTemplate embed.FS

type Entry struct {
	Number int
	ID     string // zero-padded string (e.g., 0001)
	Title  string
	Status string
	Date   string
	File   string // relative path/filename
}

func Scan(dir string) ([]Entry, error) {
	ents := []Entry{}
	items, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, it := range items {
		if it.IsDir() || !strings.HasSuffix(it.Name(), ".md") || strings.EqualFold(it.Name(), "index.md") {
			continue
		}
		path := filepath.Join(dir, it.Name())
		meta, err := ParseADR(path)
		if err != nil {
			// Best-effort: attempt to keep going, but include a minimal entry
			if n, ok := parseLeadingNumber(it.Name()); ok {
				ents = append(ents, Entry{Number: n, ID: fmt.Sprintf("%04d", n), Title: it.Name(), Status: "", Date: "", File: it.Name()})
			}
			continue
		}
		ents = append(ents, Entry{
			Number: meta.Number,
			ID:     fmt.Sprintf("%04d", meta.Number),
			Title:  meta.Title,
			Status: meta.Status,
			Date:   meta.Date,
			File:   it.Name(),
		})
	}
	// sort by Number
	sort.Slice(ents, func(i, j int) bool { return ents[i].Number < ents[j].Number })
	return ents, nil
}

type IndexData struct {
	Entries     []Entry
	ProjectName string
	ProjectURL  string
}

func WriteIndex(out string, entries []Entry, projectName, projectURL string) error {
	// Escape pipe characters in entries
	for i := range entries {
		entries[i].Title = escapePipes(entries[i].Title)
		entries[i].Status = escapePipes(entries[i].Status)
		entries[i].Date = escapePipes(entries[i].Date)
	}

	data := IndexData{
		Entries:     entries,
		ProjectName: projectName,
		ProjectURL:  projectURL,
	}

	// Load and parse template
	tmplContent, err := indexTemplate.ReadFile("templates/index.md")
	if err != nil {
		return fmt.Errorf("failed to read index template: %w", err)
	}

	tmpl, err := template.New("index").Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("failed to parse index template: %w", err)
	}

	// Ensure output directory exists
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return err
	}

	// Create output file
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()

	// Execute template
	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("failed to execute index template: %w", err)
	}

	return nil
}

func escapePipes(s string) string {
	return strings.ReplaceAll(s, "|", "\\|")
}
