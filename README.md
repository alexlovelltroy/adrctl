# adrctl

A small, dependency-light Go CLI for managing Architecture Decision Records (ADRs). Works great locally and inside GitHub Actions. Includes built-in templates for **MADR** and **Nygard**, with the option to point at a custom markdown template of your own.

## Installation

### Pre-built Binaries (Recommended)
Download the latest release for your platform from [GitHub Releases](https://github.com/alexlovelltroy/adrctl/releases).

**Quick install on macOS/Linux:**
```bash
# macOS (Apple Silicon)
curl -L https://github.com/alexlovelltroy/adrctl/releases/latest/download/adrctl_0.2.0_Darwin_arm64.tar.gz | tar xz && sudo mv adrctl /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/alexlovelltroy/adrctl/releases/latest/download/adrctl_0.2.0_Darwin_x86_64.tar.gz | tar xz && sudo mv adrctl /usr/local/bin/

# Linux (x86_64)
curl -L https://github.com/alexlovelltroy/adrctl/releases/latest/download/adrctl_0.2.0_Linux_x86_64.tar.gz | tar xz && sudo mv adrctl /usr/local/bin/

# Linux (ARM64)
curl -L https://github.com/alexlovelltroy/adrctl/releases/latest/download/adrctl_0.2.0_Linux_arm64.tar.gz | tar xz && sudo mv adrctl /usr/local/bin/
```

**Supported platforms:**
- **Linux**: x86_64, ARM64 (with deb, rpm, apk packages)
- **macOS**: Intel (x86_64), Apple Silicon (ARM64)

### Go Install
For Go developers:
```bash
go install github.com/alexlovelltroy/adrctl/cmd/adrctl@latest
```

### From Source
```bash
git clone https://github.com/alexlovelltroy/adrctl.git
cd adrctl
go build -o adrctl ./cmd/adrctl
```

## Features
- `adrctl init` — scaffold an ADR directory (defaults to `ADRs/`).
- `adrctl new "Title"` — create a new ADR with incremental ID and selected template.
- `adrctl index` — scan ADRs and generate/update `index.md`.
- Built-in templates or bring your own: `madr`, `nygard`; or `--template path/to/template.md`.  See [Nygard](https://www.cognitect.com/blog/2011/11/15/documenting-architecture-decisions) and [madr Release](https://github.com/adr/madr/releases) for details on these popular formats.
- Parses title, number, status, and date from ADR files using YAML frontmatter or markdown parsing.

## Quick start
```bash
# initialize (uses default ADRs/ directory)
adrctl init

# create ADR using MADR (default template)
adrctl new "Adopt DuckDB for local analytics"

# create ADR using Nygard template
adrctl new "JWT issuer and JWKS layout" --template nygard --status Accepted

# or with a custom template
adrctl new "Pick core logging library" --template ./my-template.md

# generate the index (outputs to ADRs/index.md by default)
adrctl index

# generate branded index with project info
adrctl index --project-name "My Project" --project-url "https://github.com/myorg/project"

# specify custom output location
adrctl index --out docs/decisions/index.md
```

## GitHub Actions
Use `actions/setup-go` and run `adrctl index` on every PR/push to keep the index up to date.

```yaml
name: ADR Index
on:
  push:
    paths:
      - 'ADRs/**.md'
      - '.github/workflows/adr-index.yml'
  pull_request:
    paths:
      - 'ADRs/**.md'

jobs:
  adr-index:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Install adrctl
        run: |
          curl -L https://github.com/alexlovelltroy/adrctl/releases/latest/download/adrctl_0.2.0_Linux_x86_64.tar.gz | tar xz
          sudo mv adrctl /usr/local/bin/
      - name: Generate ADR index
        run: adrctl index --project-name "${{ github.event.repository.name }}" --project-url "${{ github.event.repository.html_url }}"
      - name: Commit changes (if any)
        run: |
          if [[ -n "$(git status --porcelain)" ]]; then
            git config user.name github-actions
            git config user.email github-actions@github.com
            git add ADRs/index.md
            git commit -m "chore(adr): update index"
            git push
          fi
```

> **💡 Need more advanced integration?** Check out the [`examples/`](examples/) directory for:
> - Pre-commit hooks for local development
> - GitHub workflows with GPG commit signing
> - Additional integration patterns

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

## Contributing

### Development
```bash
git clone https://github.com/alexlovelltroy/adrctl.git
cd adrctl
go mod tidy
go build -o adrctl ./cmd/adrctl
```

### Releases
Releases are automated using GoReleaser:
1. Create and push a git tag: `git tag v1.0.0 && git push origin v1.0.0`
2. GitHub Actions automatically builds and publishes cross-platform binaries
3. Binaries are available on the [releases page](https://github.com/alexlovelltroy/adrctl/releases)
4. Package managers (Homebrew, apt, etc.) are automatically updated

## License

[MIT License](/LICENSE)
