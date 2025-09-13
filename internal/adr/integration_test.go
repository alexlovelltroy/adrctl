package adr

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCompleteWorkflow demonstrates the typical ADR workflow
func TestCompleteWorkflow(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Initialize ADR manager
	m := Manager{Dir: tmpDir}
	
	// Test 1: Create first ADR with MADR template
	path1, err := m.WriteNewADR("Use PostgreSQL Database", NewOptions{
		Template: "madr",
		Status:   "Proposed",
		Date:     "2025-01-15",
	})
	if err != nil {
		t.Fatalf("Failed to create first ADR: %v", err)
	}
	
	// Verify the file exists and has correct name
	expectedFile1 := filepath.Join(tmpDir, "0001-use-postgresql-database.md")
	if path1 != expectedFile1 {
		t.Errorf("Expected path %s, got %s", expectedFile1, path1)
	}
	
	// Test 2: Create second ADR with Nygard template
	path2, err := m.WriteNewADR("Implement Microservices Architecture", NewOptions{
		Template: "nygard", 
		Status:   "Accepted",
		Date:     "2025-01-16",
	})
	if err != nil {
		t.Fatalf("Failed to create second ADR: %v", err)
	}
	
	expectedFile2 := filepath.Join(tmpDir, "0002-implement-microservices-architecture.md")
	if path2 != expectedFile2 {
		t.Errorf("Expected path %s, got %s", expectedFile2, path2)
	}
	
	// Test 3: Parse the created ADRs
	meta1, err := ParseADR(path1)
	if err != nil {
		t.Fatalf("Failed to parse first ADR: %v", err)
	}
	
	if meta1.Number != 1 {
		t.Errorf("First ADR: expected number 1, got %d", meta1.Number)
	}
	if meta1.Title != "Use PostgreSQL Database" {
		t.Errorf("First ADR: expected title 'Use PostgreSQL Database', got %q", meta1.Title)
	}
	if meta1.Status != "Proposed" {
		t.Errorf("First ADR: expected status 'Proposed', got %q", meta1.Status)
	}
	
	meta2, err := ParseADR(path2)
	if err != nil {
		t.Fatalf("Failed to parse second ADR: %v", err)
	}
	
	if meta2.Number != 2 {
		t.Errorf("Second ADR: expected number 2, got %d", meta2.Number)
	}
	
	// Test 4: Scan directory and generate index
	entries, err := Scan(tmpDir)
	if err != nil {
		t.Fatalf("Failed to scan directory: %v", err)
	}
	
	if len(entries) != 2 {
		t.Fatalf("Expected 2 entries, got %d", len(entries))
	}
	
	// Verify entries are sorted by number
	if entries[0].Number != 1 || entries[1].Number != 2 {
		t.Errorf("Entries not sorted correctly: got numbers %d, %d", entries[0].Number, entries[1].Number)
	}
	
	// Test 5: Generate index file
	indexPath := filepath.Join(tmpDir, "index.md")
	err = WriteIndex(indexPath, entries, "Test Project", "https://example.com/project")
	if err != nil {
		t.Fatalf("Failed to write index: %v", err)
	}
	
	// Verify index file exists and has expected content
	indexContent, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("Failed to read index file: %v", err)
	}
	
	indexStr := string(indexContent)
	if !strings.Contains(indexStr, "Test Project - Architecture Decision Records") {
		t.Error("Index should contain project name in header")
	}
	if !strings.Contains(indexStr, "https://example.com/project") {
		t.Error("Index should contain project URL")
	}
	if !strings.Contains(indexStr, "Use PostgreSQL Database") {
		t.Error("Index should contain first ADR title")
	}
	if !strings.Contains(indexStr, "Implement Microservices Architecture") {
		t.Error("Index should contain second ADR title")
	}
	
	t.Logf("Successfully completed workflow with 2 ADRs and generated index")
}

// TestTemplateFunctionality demonstrates template usage
func TestTemplateFunctionality(t *testing.T) {
	tmpDir := t.TempDir()
	m := Manager{Dir: tmpDir}
	
	templates := []struct {
		name     string
		template string
	}{
		{"MADR", "madr"},
		{"Nygard", "nygard"},
	}
	
	for i, tt := range templates {
		t.Run(tt.name, func(t *testing.T) {
			title := "Test " + tt.name + " Template"
			path, err := m.WriteNewADR(title, NewOptions{
				Template: tt.template,
				Status:   "Draft", 
				Date:     "2025-01-15",
			})
			if err != nil {
				t.Fatalf("Failed to create %s ADR: %v", tt.name, err)
			}
			
			// Read and verify the created file
			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("Failed to read %s ADR: %v", tt.name, err)
			}
			
			contentStr := string(content)
			
			// All templates should have frontmatter
			if !strings.HasPrefix(contentStr, "---\n") {
				t.Errorf("%s template should start with frontmatter", tt.name)
			}
			
			// Should contain the title and status
			if !strings.Contains(contentStr, title) {
				t.Errorf("%s template should contain title %q", tt.name, title)
			}
			if !strings.Contains(contentStr, "Draft") {
				t.Errorf("%s template should contain status 'Draft'", tt.name)
			}
			
			// Parse to ensure it's valid
			meta, err := ParseADR(path)
			if err != nil {
				t.Fatalf("Failed to parse %s ADR: %v", tt.name, err)
			}
			
			if meta.Number != i+1 {
				t.Errorf("%s ADR: expected number %d, got %d", tt.name, i+1, meta.Number)
			}
			if meta.Title != title {
				t.Errorf("%s ADR: expected title %q, got %q", tt.name, title, meta.Title)
			}
			if meta.Status != "Draft" {
				t.Errorf("%s ADR: expected status 'Draft', got %q", tt.name, meta.Status)
			}
		})
	}
}

// TestDirectoryOperations shows how to work with ADR directories
func TestDirectoryOperations(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Test creating ADR directory
	adrDir := filepath.Join(tmpDir, "architecture-decisions")
	err := EnsureDir(adrDir)
	if err != nil {
		t.Fatalf("Failed to create ADR directory: %v", err)
	}
	
	// Verify directory exists
	info, err := os.Stat(adrDir)
	if err != nil {
		t.Fatalf("ADR directory should exist: %v", err)
	}
	if !info.IsDir() {
		t.Error("ADR path should be a directory")
	}
	
	// Test scanning empty directory
	entries, err := Scan(adrDir)
	if err != nil {
		t.Fatalf("Failed to scan empty directory: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("Empty directory should have 0 entries, got %d", len(entries))
	}
	
	// Create some ADRs
	m := Manager{Dir: adrDir}
	for i := 1; i <= 3; i++ {
		_, err := m.WriteNewADR("Test ADR", NewOptions{
			Template: "madr",
			Status:   "Proposed",
		})
		if err != nil {
			t.Fatalf("Failed to create ADR %d: %v", i, err)
		}
	}
	
	// Test scanning populated directory
	entries, err = Scan(adrDir)
	if err != nil {
		t.Fatalf("Failed to scan populated directory: %v", err)
	}
	if len(entries) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(entries))
	}
	
	// Verify entries are properly numbered and sorted
	for i, entry := range entries {
		expectedNum := i + 1
		if entry.Number != expectedNum {
			t.Errorf("Entry %d: expected number %d, got %d", i, expectedNum, entry.Number)
		}
		if entry.ID != "000"+string(rune('0'+expectedNum)) {
			t.Errorf("Entry %d: expected ID with zero padding", i)
		}
	}
}

// TestIndexGeneration demonstrates index generation options
func TestIndexGeneration(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create some test ADRs with different statuses
	entries := []Entry{
		{Number: 1, ID: "0001", Title: "First Decision", Status: "Accepted", Date: "2025-01-15", File: "0001-first.md"},
		{Number: 2, ID: "0002", Title: "Second Decision", Status: "Proposed", Date: "2025-01-16", File: "0002-second.md"},
		{Number: 3, ID: "0003", Title: "Third Decision", Status: "Superseded", Date: "2025-01-17", File: "0003-third.md"},
	}
	
	tests := []struct {
		name        string
		projectName string
		projectURL  string
		checkFunc   func(content string) error
	}{
		{
			name:        "Basic index",
			projectName: "",
			projectURL:  "",
			checkFunc: func(content string) error {
				if !strings.Contains(content, "# Architecture Decision Records") {
					return nil // Expected default header
				}
				if strings.Contains(content, "Project:") {
					t.Error("Basic index should not contain project info")
				}
				return nil
			},
		},
		{
			name:        "Branded index",
			projectName: "My Project",
			projectURL:  "https://github.com/user/project",
			checkFunc: func(content string) error {
				if !strings.Contains(content, "My Project - Architecture Decision Records") {
					t.Error("Branded index should contain project name in header")
				}
				if !strings.Contains(content, "https://github.com/user/project") {
					t.Error("Branded index should contain project URL")
				}
				return nil
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			indexPath := filepath.Join(tmpDir, tt.name+"-index.md")
			
			err := WriteIndex(indexPath, entries, tt.projectName, tt.projectURL)
			if err != nil {
				t.Fatalf("Failed to write index: %v", err)
			}
			
			content, err := os.ReadFile(indexPath)
			if err != nil {
				t.Fatalf("Failed to read index: %v", err)
			}
			
			contentStr := string(content)
			
			// Common checks
			for _, entry := range entries {
				if !strings.Contains(contentStr, entry.Title) {
					t.Errorf("Index should contain entry title: %s", entry.Title)
				}
				if !strings.Contains(contentStr, entry.Status) {
					t.Errorf("Index should contain entry status: %s", entry.Status)
				}
			}
			
			// Specific checks
			if tt.checkFunc != nil {
				if err := tt.checkFunc(contentStr); err != nil {
					t.Error(err)
				}
			}
			
			t.Logf("Generated %s index with %d entries", tt.name, len(entries))
		})
	}
}