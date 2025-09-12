# adrctl

A small, dependency-light Go CLI for managing Architecture Decision Records (ADRs). Works great locally and inside GitHub Actions. Includes built-in templates for **MADR** and **Nygard**, with the option to point at a custom `template.md`.

## Features
- `adrctl init` — scaffold an ADR directory (defaults to `docs/adr`).
- `adrctl new "Title"` — create a new ADR with incremental ID and selected template.
- `adrctl index` — scan ADRs and generate/update `index.md`.
- Built-in templates: `madr`, `nygard`; or `--template path/to/template.md`.
- Parses title, number, status, and date from ADR files using YAML frontmatter or markdown parsing.

## Quick start
```bash
# build
go build -o adrctl ./cmd/adrctl

# initialize
./adrctl init --dir docs/adr

# create ADR using MADR (default)
./adrctl new "Adopt DuckDB for local analytics" --dir docs/adr --template madr --status Proposed

# create ADR using Nygard
./adrctl new "JWT issuer and JWKS layout" --dir docs/adr --template nygard --status Accepted

# or with a custom template
./adrctl new "Pick core logging library" --dir docs/adr --template ./my-template.md

# generate the index
./adr index --dir docs/adr --out docs/adr/index.md
```

## GitHub Actions
Use `actions/setup-go` and run `adr index` on every PR/push to keep the index up to date.

```yaml
name: ADR Index
on:
  push:
    paths:
      - 'docs/adr/**.md'
      - '.github/workflows/adr-index.yml'
  pull_request:
    paths:
      - 'docs/adr/**.md'

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22.x'
      - name: Import GPG key
        uses: crazy-max/ghaction-import-gpg@v6
        with:
          gpg_private_key: ${{ secrets.GPG_PRIVATE_KEY }}
          passphrase: ${{ secrets.GPG_PASSPHRASE }}
          git_user_signingkey: true
          git_commit_gpgsign: true
      - name: Build adrctl
        run: go build -o adrctl ./cmd/adrctl
      - name: Generate ADR index
        run: ./adr index --dir docs/adr --out docs/adr/index.md
      - name: Commit changes (if any)
        run: |
          if [[ -n "$(git status --porcelain)" ]]; then
            git add docs/adr/index.md
            git commit -m "chore(adr): update index"
          fi
```

## Conventions
- Filenames: `NNNN-kebab-title.md` (e.g., `0001-adopt-duckdb.md`).
- Title header: `# ADR NNNN: Title`.
- **Frontmatter support**: All templates now include YAML frontmatter for structured metadata:
  ```yaml
  ---
  id: 1
  title: "ADR Title"
  status: "Proposed"
  date: "2025-01-15"
  ---
  ```
- **Backward compatibility**: Legacy parsing still supports various markdown formats:
  - Status heading: `## Status` followed by status value
  - Inline status: `**Status:** value`, `- Status: value`, or `Status: value`
- Date: extracted from frontmatter, or derived from file mtime or `Date:` line in the document; can be overridden on `adr new`.

## Exit codes (CI-friendly)
- `0`: success
- `1`: usage error or invalid flags
- `2`: filesystem/template issues
