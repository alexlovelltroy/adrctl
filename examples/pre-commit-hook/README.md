# Pre-commit Hook for ADR Index Updates

This pre-commit hook automatically updates your ADR index file whenever ADR files are modified and committed.

## Installation

1. **Install adrctl** (if not already installed):
   ```bash
   go install github.com/alexlovelltroy/adrctl/cmd/adrctl@latest
   ```

2. **Copy the hook**:
   ```bash
   cp examples/pre-commit-hook/pre-commit .git/hooks/pre-commit
   ```

3. **Make it executable**:
   ```bash
   chmod +x .git/hooks/pre-commit
   ```

## How it works

- Runs before each commit
- Detects if any `.md` files in the `ADRs/` directory (excluding `index.md`) are being committed
- If ADR files are modified, automatically runs `adrctl index` to update the index
- Adds the updated index file to the commit

## Configuration

Edit the hook file to adjust:
- `ADR_DIR`: Change from `ADRs` if your ADR directory is different
- `INDEX_FILE`: Change the index file location if needed

## Disabling

To temporarily disable the hook for a commit:
```bash
git commit --no-verify
```

To permanently remove:
```bash
rm .git/hooks/pre-commit
```