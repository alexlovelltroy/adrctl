# adrctl Examples

This directory contains practical examples for integrating adrctl into your development workflow.

## 📁 Directory Structure

```
examples/
├── pre-commit-hook/          # Git pre-commit hook for automatic index updates
│   ├── pre-commit           # Hook script
│   └── README.md           # Installation and usage instructions
├── github-workflows/        # GitHub Actions workflows
│   ├── adr-index-with-signing.yml    # Full-featured workflow with GPG signing
│   └── README.md           # Workflow documentation
└── README.md               # This file
```

## 🚀 Quick Setup

### For Simple Projects
Use the basic GitHub workflow example in the main README.

### For Projects Requiring Signed Commits
1. Use the workflow in `github-workflows/adr-index-with-signing.yml`
2. Configure GPG secrets as documented

### For Local Development
Install the pre-commit hook to automatically update indices before commits.

## 📋 Integration Options

| Method | Use Case | Complexity | Features |
|--------|----------|------------|----------|
| **Basic GitHub Workflow** | Simple projects, no signing requirements | Low | Automatic index updates |
| **GPG Workflow** | Enterprise/secure repos requiring signed commits | Medium | Signed commits, PR comments |
| **Pre-commit Hook** | Local development, immediate feedback | Low | Local validation, fast feedback |

## 🛠️ Customization

All examples use the default `ADRs/` directory. To customize:

1. **Change ADR directory**: Update paths in workflows and hooks
2. **Custom index location**: Modify `--out` parameter in commands
3. **Additional validation**: Add steps for linting, formatting, etc.

## 💡 Best Practices

- **Combine approaches**: Use pre-commit hooks for local development + GitHub workflows for CI
- **Test before deploying**: Try workflows on feature branches first
- **Monitor workflow runs**: Check Actions tab for any issues
- **Keep secrets secure**: Rotate GPG keys periodically if using signed commits

## 🔗 Related Documentation

- [Main README](../README.md) - Basic usage and simple GitHub workflow
- [adrctl Installation](../README.md#installation) - Getting started with adrctl