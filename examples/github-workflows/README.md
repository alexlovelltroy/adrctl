# GitHub Workflows for adrctl

This directory contains example GitHub workflows for automatically maintaining ADR indices.

## adr-index-with-signing.yml

A comprehensive workflow that includes:

- **GPG commit signing** for repositories requiring signed commits
- **Automatic index updates** when ADR files are modified
- **PR comments** to notify about index updates
- **Proper permissions** and error handling
- **Full git history** for reliable operations

### Setup

1. **Copy the workflow**:
   ```bash
   cp examples/github-workflows/adr-index-with-signing.yml .github/workflows/
   ```

2. **Configure GPG secrets** (required for signed commits):
   - `GPG_PRIVATE_KEY`: Your GPG private key (export with `gpg --armor --export-secret-keys KEY_ID`)
   - `GPG_PASSPHRASE`: Passphrase for your GPG key

3. **Add GPG public key** to your GitHub account for signature verification

### Features

- ✅ **Signed commits**: Complies with repositories requiring commit signing
- ✅ **PR support**: Works on both push and pull request events
- ✅ **Change detection**: Only commits when the index actually changes
- ✅ **Proper attribution**: Includes co-author attribution for GitHub Actions
- ✅ **PR notifications**: Comments on PRs when index is updated
- ✅ **Path filtering**: Only runs when ADR files are modified

### Workflow Behavior

**On Push to main/master:**
- Generates new ADR index
- Commits and pushes changes (if any) with GPG signature

**On Pull Request:**
- Generates new ADR index
- Commits changes to the PR branch (if any)
- Adds a comment to the PR about the index update

### Customization

Edit the workflow file to adjust:
- **ADR directory**: Change `'ADRs/**.md'` if your directory is different
- **Index output**: Modify `--out ADRs/index.md` for different index location
- **Trigger paths**: Add additional paths that should trigger index updates
- **Commit message**: Customize the commit message format

### Alternative: Simple Workflow

For repositories that don't require signed commits, see the main README for a simpler workflow example without GPG complexity.