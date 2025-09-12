package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/alexlovelltroy/adrctl/internal/adr"
)

var (
	flagDir      string
	flagTemplate string
	flagStatus   string
	flagDate     string
	flagOut      string
)

func main() {
	root := &cobra.Command{
		Use:   "adr",
		Short: "Manage Architecture Decision Records (ADRs)",
	}

	root.PersistentFlags().StringVar(&flagDir, "dir", "docs/adr", "ADR directory")

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
				return errors.New("--out is required (e.g., --out docs/adr/index.md)")
			}
			entries, err := adr.Scan(flagDir)
			if err != nil {
				return err
			}
			if err := adr.WriteIndex(flagOut, entries); err != nil {
				return err
			}
			fmt.Fprintln(os.Stdout, flagOut)
			return nil
		},
	}
	cmdIndex.Flags().StringVar(&flagOut, "out", "", "Output index path (e.g., docs/adr/index.md)")

	root.AddCommand(cmdInit, cmdNew, cmdIndex)

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
