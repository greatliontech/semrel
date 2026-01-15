# semrel

A semantic versioning and release automation tool, based on conventional commits.

## Features

- Automatic version calculation based on conventional commit messages
- Support for GitHub and GitLab releases
- Configurable commit type mappings
- Release notes generation with filtering and regex transformations
- Development mode for pre-1.0.0 projects
- GitHub Actions integration
- JSON output for CI/CD pipelines

## Installation

### Homebrew

```bash
brew install greatliontech/tap/semrel
```

### Go Install

```bash
go install github.com/greatliontech/semrel@latest
```

### Binary Downloads

Download pre-built binaries from the [releases page](https://github.com/greatliontech/semrel/releases).

## Quick Start

```bash
# Show the next version based on commits
semrel

# Show the current version
semrel current

# Create a release on GitHub/GitLab
semrel release

# Validate a version string
semrel validate 1.2.3

# Compare versions
semrel compare --ge 1.0.0
```

## Commands

### Default Command

Calculate and display the next semantic version based on conventional commits since the last tag.

```bash
semrel [flags]

Flags:
  -p, --prerelease string    Add prerelease suffix (e.g., alpha.1)
  -b, --build string         Add build metadata
      --create-tag           Create the git tag
      --push-tag             Push the tag to remote
      --current-branch-only  Only consider tags from current branch
      --auth-username        Username for basic auth
      --auth-password        Password for basic auth
      --auth-token           Token for auth
  -v, --verbose              Enable verbose output
      --dry-run              Simulate actions without making changes
```

### current

Print the current (latest) release version.

```bash
semrel current [flags]

Flags:
      --current-branch-only  Only consider tags from current branch
      --json                 Output in JSON format
```

### release

Create a release on GitHub or GitLab with generated release notes.

```bash
semrel release [flags]

Flags:
  -p, --prerelease string    Add prerelease suffix
  -b, --build string         Add build metadata
      --current-branch-only  Only consider tags from current branch
      --json                 Output in JSON format
```

### compare

Compare the current or supplied version against constraints.

```bash
semrel compare [version] [flags]

Flags:
      --le strings           Less than or equal to
      --ge strings           Greater than or equal to
      --lt strings           Less than
      --gt strings           Greater than
      --current-branch-only  Only compare the current branch
      --json                 Output in JSON format
```

### validate

Validate a semantic version string.

```bash
semrel validate <version> [flags]

Flags:
      --strict         Strict semver validation
      --noPrerelease   Do not allow pre-release versions
      --noBuild        Do not allow build metadata
      --json           Output in JSON format
```

## Configuration

Create a `.semrel.yaml` file in your repository root:

```yaml
# Commit types that trigger a patch version bump
patchTypes:
  - fix
  - chore
  - perf

# Commit types that trigger a minor version bump
minorTypes:
  - feat

# Commit types that trigger a major version bump
majorTypes:
  - breaking

# Default bump when no commits match types (none, patch, minor, major)
defaultBump: none

# Initial version for new repositories
initialVersion: 1.0.0

# Development mode: starts at 0.1.0, breaking changes are minor bumps
development: false

# Version prefix (e.g., "v" for v1.0.0)
prefix: "v"

# Automatically create git tag
createTag: false

# Automatically push tag to remote
pushTag: false

# Platform for releases (github, gitlab)
platform: github

# Filters for release notes
filters:
  types:
    - docs
    - test
  scopes:
    - internal

# Regex rules for transforming commit messages in release notes
matchRules:
  - match: "TRACK-(\\d+)"
    replace: "[#$1](https://tracker.example.com/TRACK-$1)"
```

## Conventional Commits

semrel parses commit messages following the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

### Breaking Changes

Breaking changes can be indicated in two ways:

1. **Attention marker** - Add `!` after the type/scope:

   ```
   feat!: remove deprecated API
   feat(api)!: change response format
   ```

2. **Footer** - Add `BREAKING CHANGE:` in the commit body:

   ```
   feat: update authentication

   BREAKING CHANGE: JWT tokens now expire after 1 hour
   ```

### Version Bumps

| Commit Type | Version Bump |
|-------------|--------------|
| `fix:` | Patch (1.0.0 -> 1.0.1) |
| `feat:` | Minor (1.0.0 -> 1.1.0) |
| `feat!:` or `BREAKING CHANGE:` | Major (1.0.0 -> 2.0.0) |

### Development Mode

When `development: true` is set:

- Initial version starts at `0.1.0` instead of `1.0.0`
- Breaking changes on `0.x.x` versions bump the patch version instead of major
- Once you reach `1.0.0`, breaking changes behave normally

## Environment Variables

| Variable | Description |
|----------|-------------|
| `SEMREL_PRERELEASE` | Prerelease suffix |
| `SEMREL_BUILD` | Build metadata |
| `SEMREL_PLATFORM` | Override platform detection |
| `SEMREL_TOKEN` | Platform authentication token |
| `SEMREL_PROJECT` | Repository identifier (owner/repo) |
| `SEMREL_BRANCH` | Release target branch |
| `SEMREL_AUTH_USERNAME` | Basic auth username |
| `SEMREL_AUTH_PASSWORD` | Basic auth password |
| `SEMREL_AUTH_TOKEN` | Token for git operations |

## GitHub Actions

semrel can be used as a GitHub Action:

```yaml
name: Release

on:
  push:
    branches: [main]

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - uses: greatliontech/semrel@v0
        id: semrel
        with:
          token: ${{ secrets.GITHUB_TOKEN }}
          release: true

      - run: echo "Next version is ${{ steps.semrel.outputs.next-version }}"
```

### Action Versioning

The action follows standard GitHub Action versioning conventions:

| Reference | Description |
|-----------|-------------|
| `@v0` | Latest v0.x.x release (recommended for stability) |
| `@v0.5` | Latest v0.5.x release |
| `@v0.5.2` | Exact version (immutable) |
| `@main` | Latest development (not recommended for production) |

Major version tags (e.g., `v0`, `v1`) are automatically updated on each release to point to the latest minor/patch version, so you get bug fixes automatically while maintaining API stability.

### Action Inputs

| Input | Description | Default |
|-------|-------------|---------|
| `token` | GitHub token for releases | `${{ github.token }}` |
| `tag` | Tool version to download | `latest` |
| `prerelease` | Prerelease suffix | |
| `build` | Build metadata | |
| `prefix` | Version prefix | |
| `release` | Create GitHub release | `false` |

### Action Outputs

| Output | Description |
|--------|-------------|
| `next-version` | The calculated next version |
| `current-version` | The current/latest version |

## Exit Codes

| Code | Description |
|------|-------------|
| 0 | Success |
| 1 | General error |
| 2 | Validation error |
| 3 | Configuration error |
| 4 | Repository error |
| 5 | Release error |
| 6 | Authentication error |

## JSON Output

All commands support `--json` flag for machine-readable output:

```bash
# Current version
semrel current --json
# {"version":"1.2.3","tag":"v1.2.3"}

# Validate
semrel validate 1.2.3 --json
# {"input":"1.2.3","valid":true,"version":"1.2.3"}

# Compare
semrel compare --ge 1.0.0 --json
# {"version":"1.2.3","result":true}

# Release
semrel release --json
# {"current":"v1.2.3","next":"v1.3.0","created":true}
```

## License

Apache License 2.0 - see [LICENSE](LICENSE) for details.
