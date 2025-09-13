package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/alexlovelltroy/adrctl/internal/adr"
)

var (
	// Build-time variables (set by goreleaser)
	version = "dev"
	commit  = "none"
	date    = "unknown"

	// Command flags
	flagDir         string
	flagTemplate    string
	flagStatus      string
	flagDate        string
	flagOut         string
	flagProjectName string
	flagProjectURL  string
)

func main() {
	root := &cobra.Command{
		Use:     "adrctl",
		Short:   "Manage Architecture Decision Records (ADRs)",
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
	}

	root.PersistentFlags().StringVar(&flagDir, "dir", "ADRs", "ADR directory")

	cmdInit := &cobra.Command{
		Use:   "init",
		Short: "Initialize ADR directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			return adr.EnsureDir(flagDir)
		},
	}

	cmdNew := &cobra.Command{
		Use:   "new [title]",
		Short: "Create a new ADR from a template",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			m := adr.Manager{Dir: flagDir}
			title := args[0]
			opt := adr.NewOptions{Template: flagTemplate, Status: flagStatus, Date: flagDate}
			path, err := m.WriteNewADR(title, opt)
			if err != nil {
				return err
			}
			fmt.Println(path)
			return nil
		},
	}
	cmdNew.Flags().StringVar(&flagTemplate, "template", "madr", "Template to use: madr|nygard|/path/to/template.md")
	cmdNew.Flags().StringVar(&flagStatus, "status", "Proposed", "Initial ADR status")
	cmdNew.Flags().StringVar(&flagDate, "date", "", "ISO date (YYYY-MM-DD); defaults to today")

	cmdIndex := &cobra.Command{
		Use:   "index",
		Short: "Generate or update index.md for ADRs",
		RunE: func(cmd *cobra.Command, args []string) error {
			if flagOut == "" {
				flagOut = filepath.Join(flagDir, "index.md")
			}
			entries, err := adr.Scan(flagDir)
			if err != nil {
				return err
			}
			if err := adr.WriteIndex(flagOut, entries, flagProjectName, flagProjectURL); err != nil {
				return err
			}
			fmt.Fprintln(os.Stdout, flagOut)
			return nil
		},
	}
	cmdIndex.Flags().StringVar(&flagOut, "out", "", "Output index path (defaults to <dir>/index.md)")
	cmdIndex.Flags().StringVar(&flagProjectName, "project-name", "", "Project name to display in index header")
	cmdIndex.Flags().StringVar(&flagProjectURL, "project-url", "", "Project URL to link in index header")

	root.AddCommand(cmdInit, cmdNew, cmdIndex)

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
