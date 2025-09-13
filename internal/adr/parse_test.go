package adr

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFrontmatter(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected Meta
	}{
		{
			name: "frontmatter with all fields",
			content: `---
id: 42
title: "Test ADR"
status: "Accepted" 
date: "2025-01-15"
---

# ADR 42: Test ADR

Some content here.`,
			expected: Meta{
				Number: 42,
				Title:  "Test ADR",
				Status: "Accepted",
				Date:   "2025-01-15",
			},
		},
		{
			name: "frontmatter with string id",
			content: `---
id: "0001"
title: "String ID Test"
status: "Proposed"
date: "2025-01-15"
---

Content here.`,
			expected: Meta{
				Number: 1,
				Title:  "String ID Test",
				Status: "Proposed",
				Date:   "2025-01-15",
			},
		},
		{
			name: "legacy MADR format",
			content: `# ADR 0003: Use MADR format

- Status: Accepted  
Date: 2025-01-15

## Context and Problem Statement`,
			expected: Meta{
				Number: 3,
				Title:  "Use MADR format",
				Status: "Accepted",
				Date:   "2025-01-15",
			},
		},
		{
			name: "legacy Nygard format", 
			content: `# ADR 0005: Use microservices

## Status
Superseded

Date: 2025-01-15

## Context`,
			expected: Meta{
				Number: 5,
				Title:  "Use microservices", 
				Status: "Superseded",
				Date:   "2025-01-15",
			},
		},
		{
			name: "bold status format",
			content: `# ADR 0007: Bold Status Format

**Status:** Proposed

Date: 2025-01-15

## Context`,
			expected: Meta{
				Number: 7,
				Title:  "Bold Status Format",
				Status: "Proposed",
				Date:   "2025-01-15",
			},
		},
		{
			name: "mixed frontmatter and content parsing",
			content: `---
id: 8
status: "Draft"
---

# ADR 0008: Mixed Format

Content with different date format.
Date: 2025-01-20`,
			expected: Meta{
				Number: 8,
				Title:  "Mixed Format", // Should be parsed from header
				Status: "Draft",       // From frontmatter
				Date:   "2025-01-20", // From content
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "test.md")
			
			err := os.WriteFile(tmpFile, []byte(tt.content), 0644)
			if err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			// Parse the file
			result, err := ParseADR(tmpFile)
			if err != nil {
				t.Fatalf("ParseADR failed: %v", err)
			}

			// Check results
			if result.Number != tt.expected.Number {
				t.Errorf("Number: got %d, want %d", result.Number, tt.expected.Number)
			}
			if result.Title != tt.expected.Title {
				t.Errorf("Title: got %q, want %q", result.Title, tt.expected.Title)
			}
			if result.Status != tt.expected.Status {
				t.Errorf("Status: got %q, want %q", result.Status, tt.expected.Status)
			}
			if result.Date != tt.expected.Date {
				t.Errorf("Date: got %q, want %q", result.Date, tt.expected.Date)
			}
		})
	}
}

func TestParseEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		filename string
		expected Meta
	}{
		{
			name:     "no frontmatter, derive from filename",
			content:  `Some content without proper headers`,
			filename: "1234-this-is-a-test.md",
			expected: Meta{
				Number: 0, // Parser doesn't extract numbers from filenames without headers
				Title:  "this is a test", // Derived from filename
				Status: "",
				Date:   "", // Will be file mod time, we'll skip checking this
			},
		},
		{
			name:     "malformed frontmatter falls back to content",
			filename: "0099-fallback-test.md",
			content: `---
invalid yaml: [
---

# ADR 0099: Fallback Test

Status: Working`,
			expected: Meta{
				Number: 99,
				Title:  "Fallback Test",
				Status: "Working",
				Date:   "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, tt.filename)
			
			err := os.WriteFile(tmpFile, []byte(tt.content), 0644)
			if err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			result, err := ParseADR(tmpFile)
			if err != nil {
				t.Fatalf("ParseADR failed: %v", err)
			}

			if result.Number != tt.expected.Number {
				t.Errorf("Number: got %d, want %d", result.Number, tt.expected.Number)
			}
			if result.Title != tt.expected.Title {
				t.Errorf("Title: got %q, want %q", result.Title, tt.expected.Title)
			}
			if result.Status != tt.expected.Status {
				t.Errorf("Status: got %q, want %q", result.Status, tt.expected.Status)
			}
			// Don't check date for fallback cases as it uses file mod time
		})
	}
}

// TestParseADR_Documentation demonstrates supported ADR formats
func TestParseADR_Documentation(t *testing.T) {
	examples := map[string]string{
		"Modern with Frontmatter": `---
id: 1
title: "Use React for Frontend"
status: "Accepted"
date: "2025-01-15"
---

# ADR 1: Use React for Frontend

## Context
We need a frontend framework.

## Decision
We will use React.`,

		"MADR Format": `# ADR 0002: Use PostgreSQL

- Status: Proposed
- Date: 2025-01-15

## Context and Problem Statement

## Decision Drivers

## Considered Options

## Decision Outcome`,

		"Nygard Format": `# ADR 0003: Use Microservices

## Status
Accepted

## Date  
2025-01-15

## Context

## Decision

## Consequences`,
	}

	for name, content := range examples {
		t.Run(name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "test.md")
			
			err := os.WriteFile(tmpFile, []byte(content), 0644)
			if err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			result, err := ParseADR(tmpFile)
			if err != nil {
				t.Fatalf("ParseADR failed for %s: %v", name, err)
			}

			// All examples should parse successfully with non-empty basic fields
			if result.Number == 0 {
				t.Errorf("%s: Expected non-zero Number", name)
			}
			if result.Title == "" {
				t.Errorf("%s: Expected non-empty Title", name)
			}
			if result.Status == "" {
				t.Errorf("%s: Expected non-empty Status", name) 
			}
			if result.Date == "" {
				t.Errorf("%s: Expected non-empty Date", name)
			}

			t.Logf("%s parsed: Number=%d, Title=%q, Status=%q, Date=%q", 
				name, result.Number, result.Title, result.Status, result.Date)
		})
	}
}