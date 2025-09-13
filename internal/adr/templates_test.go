package adr

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestTemplateGeneration ensures all built-in templates work correctly
func TestTemplateGeneration(t *testing.T) {
	tmpDir := t.TempDir()
	m := Manager{Dir: tmpDir}

	templates := []struct {
		name         string
		template     string
		expectedSections []string
	}{
		{
			name:     "MADR",
			template: "madr",
			expectedSections: []string{
				"Context and Problem Statement",
				"Decision Drivers", 
				"Considered Options",
				"Decision Outcome",
				"Positive Consequences",
				"Negative Consequences",
			},
		},
		{
			name:     "Nygard",
			template: "nygard",
			expectedSections: []string{
				"Status",
				"Date", 
				"Context",
				"Decision",
				"Consequences",
			},
		},
	}

	for _, tt := range templates {
		t.Run(tt.name, func(t *testing.T) {
			title := "Test " + tt.name + " Decision"
			
			path, err := m.WriteNewADR(title, NewOptions{
				Template: tt.template,
				Status:   "Proposed",
				Date:     "2025-01-15",
			})
			if err != nil {
				t.Fatalf("Failed to create ADR with %s template: %v", tt.name, err)
			}

			// Read the generated file
			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("Failed to read generated ADR: %v", err)
			}

			contentStr := string(content)

			// Test 1: Should have frontmatter
			if !strings.HasPrefix(contentStr, "---\n") {
				t.Errorf("%s template should start with frontmatter", tt.name)
			}

			// Test 2: Should contain expected metadata
			expectedMetadata := []string{
				"id: 000", // Should have padded ID
				"title: \"" + title + "\"",
				"status: \"Proposed\"",
				"date: \"2025-01-15\"",
			}

			for _, metadata := range expectedMetadata {
				if !strings.Contains(contentStr, metadata) {
					t.Errorf("%s template should contain metadata: %s", tt.name, metadata)
				}
			}

			// Test 3: Should contain expected sections
			for _, section := range tt.expectedSections {
				if !strings.Contains(contentStr, section) {
					t.Errorf("%s template should contain section: %s", tt.name, section)
				}
			}

			// Test 4: Should be parseable
			meta, err := ParseADR(path)
			if err != nil {
				t.Fatalf("Generated %s ADR should be parseable: %v", tt.name, err)
			}

			if meta.Title != title {
				t.Errorf("%s ADR parsed title: got %q, want %q", tt.name, meta.Title, title)
			}
			if meta.Status != "Proposed" {
				t.Errorf("%s ADR parsed status: got %q, want %q", tt.name, meta.Status, "Proposed")
			}
			if meta.Date != "2025-01-15" {
				t.Errorf("%s ADR parsed date: got %q, want %q", tt.name, meta.Date, "2025-01-15")
			}

			t.Logf("%s template generated valid ADR with %d characters", tt.name, len(contentStr))
		})
	}
}

// TestCustomTemplate demonstrates custom template usage
func TestCustomTemplate(t *testing.T) {
	tmpDir := t.TempDir()
	m := Manager{Dir: tmpDir}

	// Create a custom template file
	customTemplate := `---
id: {{.ID}}
title: "{{.Title}}"
status: "{{.Status}}"
date: "{{.Date}}"
type: "custom"
---

# ADR {{.ID}}: {{.Title}}

**Status**: {{.Status}}  
**Date**: {{.Date}}  

## Problem

Describe the problem here.

## Solution

Describe the solution here.

## Impact

Describe the impact here.

## Notes

Additional notes.
`

	templatePath := filepath.Join(tmpDir, "custom.md")
	err := os.WriteFile(templatePath, []byte(customTemplate), 0644)
	if err != nil {
		t.Fatalf("Failed to write custom template: %v", err)
	}

	// Use the custom template
	title := "Custom Template Test"
	path, err := m.WriteNewADR(title, NewOptions{
		Template: templatePath,
		Status:   "Draft",
		Date:     "2025-01-20",
	})
	if err != nil {
		t.Fatalf("Failed to create ADR with custom template: %v", err)
	}

	// Read and verify the generated file
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read custom ADR: %v", err)
	}

	contentStr := string(content)

	// Should contain custom template elements
	expectedElements := []string{
		`type: "custom"`,
		"## Problem",
		"## Solution", 
		"## Impact",
		"## Notes",
		"**Status**: Draft",
	}

	for _, element := range expectedElements {
		if !strings.Contains(contentStr, element) {
			t.Errorf("Custom template should contain: %s", element)
		}
	}

	// Should be parseable
	meta, err := ParseADR(path)
	if err != nil {
		t.Fatalf("Custom ADR should be parseable: %v", err)
	}

	if meta.Title != title {
		t.Errorf("Custom ADR title: got %q, want %q", meta.Title, title)
	}
	if meta.Status != "Draft" {
		t.Errorf("Custom ADR status: got %q, want %q", meta.Status, "Draft")
	}

	t.Logf("Custom template successfully generated parseable ADR")
}

// TestTemplateVariables ensures all template variables work
func TestTemplateVariables(t *testing.T) {
	tmpDir := t.TempDir()
	m := Manager{Dir: tmpDir}

	title := "Variable Test Decision"
	status := "Experimental"  
	date := "2025-02-01"

	path, err := m.WriteNewADR(title, NewOptions{
		Template: "madr",
		Status:   status,
		Date:     date,
	})
	if err != nil {
		t.Fatalf("Failed to create test ADR: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read test ADR: %v", err)
	}

	contentStr := string(content)

	// Test that all variables are substituted correctly
	variableTests := map[string]string{
		"ID":     "0001",  // Should be zero-padded
		"Title":  title,
		"Status": status,
		"Date":   date,
	}

	for variable, expected := range variableTests {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("Template should contain substituted %s: %s", variable, expected)
		}
	}

	// Test specific template variable patterns
	patterns := []string{
		"title: \"" + title + "\"",       // Frontmatter title as quoted string
		"# ADR 0001: " + title,           // Header with padded ID
		"- Status: " + status,            // MADR status format
		"- Date: " + date,                // MADR date format
	}

	for _, pattern := range patterns {
		if !strings.Contains(contentStr, pattern) {
			t.Errorf("Template should contain pattern: %s", pattern)
		}
	}

	t.Logf("All template variables correctly substituted")
}

// TestTemplateDefaultValues tests default value handling
func TestTemplateDefaultValues(t *testing.T) {
	tmpDir := t.TempDir()  
	m := Manager{Dir: tmpDir}

	// Test with minimal options (relying on defaults)
	path, err := m.WriteNewADR("Default Values Test", NewOptions{})
	if err != nil {
		t.Fatalf("Failed to create ADR with defaults: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read ADR with defaults: %v", err)
	}

	contentStr := string(content)

	// Should use default template (madr)
	if !strings.Contains(contentStr, "Context and Problem Statement") {
		t.Error("Default template should be MADR")
	}

	// Should use default status
	if !strings.Contains(contentStr, "Proposed") {
		t.Error("Default status should be 'Proposed'")
	}

	// Should have auto-generated date
	meta, err := ParseADR(path)
	if err != nil {
		t.Fatalf("ADR with defaults should be parseable: %v", err)
	}

	if meta.Date == "" {
		t.Error("Default ADR should have auto-generated date")
	}

	// Date should be in YYYY-MM-DD format
	if !strings.Contains(meta.Date, "-") || len(meta.Date) != 10 {
		t.Errorf("Auto-generated date should be YYYY-MM-DD format, got: %s", meta.Date)
	}

	t.Logf("Default values work correctly: status=%s, date=%s", meta.Status, meta.Date)
}