package adr

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestScanFiltersNonADR verifies that Scan only returns ADR markdown files
// that begin with a numeric prefix (e.g. 0001-some-decision.md) and skips
// README.md, template.md, index.md and other non-conforming files.
func TestScanFiltersNonADR(t *testing.T) {
	dir := t.TempDir()

	// Helper to write a file
	write := func(name, content string) {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatalf("write %s failed: %v", name, err)
		}
	}

	// Non-ADR markdown files that should be ignored
	write("README.md", "# Project README\n")
	write("template.md", "Template guidance\n")
	write("index.md", "# ADR Index\n")
	write("design-notes.md", "Random notes\n")

	// Valid ADR files
	write("0001-first-decision.md", `---\nid: 1\ntitle: First Decision\nstatus: Accepted\ndate: 2025-01-01\n---\n\n# ADR 1: First Decision\n`)
	write("0002-second-decision.md", `---\nid: 2\ntitle: Second Decision\nstatus: Proposed\ndate: 2025-01-02\n---\n\n# ADR 2: Second Decision\n`)

	// Malformed ADR file (should still appear with fallback minimal metadata)
	write("0003-third-decision.md", `---\n: bad frontmatter [\n---\n# ADR 3: Third Decision\nStatus: Draft\nDate: 2025-01-03\n`)

	entries, err := Scan(dir)
	if err != nil {
		t.Fatalf("Scan failed: %v", err)
	}

	// Collect filenames for assertion
	gotFiles := map[string]Entry{}
	for _, e := range entries {
		gotFiles[e.File] = e
		// Ensure numeric prefix is present
		if !strings.HasPrefix(e.File, "000") && len(entries) <= 3 { // loose check for these fixtures
			t.Errorf("Entry file %s does not start with expected numeric prefix", e.File)
		}
	}

	// Expected ADR entries (README/template/index/design-notes must be absent)
	expected := []string{
		"0001-first-decision.md",
		"0002-second-decision.md",
		"0003-third-decision.md", // malformed but still included
	}

	if len(entries) != len(expected) {
		t.Fatalf("Expected %d ADR entries, got %d (files: %v)", len(expected), len(entries), expected)
	}

	for _, f := range expected {
		if _, ok := gotFiles[f]; !ok {
			t.Errorf("Missing expected ADR file in results: %s", f)
		}
	}

	// Verify ordering by Number (ascending)
	for i := 1; i < len(entries); i++ {
		if entries[i-1].Number > entries[i].Number {
			t.Errorf("Entries not sorted by Number ascending: %d before %d", entries[i-1].Number, entries[i].Number)
		}
	}

	// Check fallback meta for malformed file retains filename-derived basics
	third := gotFiles["0003-third-decision.md"]
	if third.Number != 3 {
		t.Errorf("Malformed ADR fallback: expected Number 3, got %d", third.Number)
	}
	if third.Title == "" {
		t.Errorf("Malformed ADR fallback: expected non-empty Title")
	}
}
