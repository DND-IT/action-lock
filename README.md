# action-lock

A GitHub Action for distributed mutex locking using git refs. Serialize concurrent workflow jobs to prevent race conditions.

## Why Use This?

**The Problem:** Multiple workflow runs (e.g., semantic-release across services in a monorepo) try to push to the same branch simultaneously, causing git push conflicts.

**The Solution:** Use git refs as atomic locks. A ref can only be created once — if it already exists, the create call returns 422, making it a natural mutex. Jobs acquire the lock before the critical section and release it when done.

## Quick Example

```yaml
jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: acquire lock
        id: lock
        uses: DND-IT/action-lock@v0
        with:
          action: acquire
          lock_name: release
          token: ${{ secrets.GITHUB_TOKEN }}

      # ... do your critical section work ...

      - name: release lock
        if: always()
        uses: DND-IT/action-lock@v0
        with:
          action: release
          lock_name: release
          token: ${{ secrets.GITHUB_TOKEN }}
```

## Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `action` | Lock action: `acquire` or `release` | Yes | |
| `lock_name` | Name of the lock (used as ref name under `refs/locks/`) | Yes | |
| `timeout` | Maximum time in seconds to wait for lock acquisition | No | `300` |
| `poll_interval` | Seconds between lock acquisition attempts | No | `10` |
| `stale_threshold` | Seconds after which a lock is considered stale and can be force-acquired. Set to `0` to disable stale detection. | No | `600` |
| `fail_on_timeout` | Fail the step if the lock cannot be acquired within timeout. Set to `false` to skip gracefully. | No | `true` |
| `token` | GitHub token with `contents:write` permission | Yes | |

## Outputs

| Output | Description |
|--------|-------------|
| `acquired` | Whether the lock was successfully acquired (`true`/`false`) |
| `lock_ref` | The full git ref used for the lock (e.g., `refs/locks/release`) |

## How It Works

1. **Acquire:** Creates a git ref `refs/locks/<lock_name>` pointing to the current commit SHA. If the ref already exists (HTTP 422), the lock is held by another process — the action retries with exponential backoff until timeout.

2. **Stale Detection:** If a lock has been held longer than `stale_threshold` seconds (based on the commit date of the SHA), it's automatically removed and re-acquired. This prevents deadlocks from crashed workflows.

3. **Release:** Deletes the git ref. Idempotent — releasing a non-existent lock is a no-op.

## Use Cases

### Serializing Semantic Release in a Monorepo

```yaml
jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: acquire release lock
        uses: DND-IT/action-lock@v0
        with:
          action: acquire
          lock_name: semantic-release
          token: ${{ secrets.GITHUB_TOKEN }}

      - name: run semantic release
        uses: cycjimmy/semantic-release-action@v4
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

      - name: release lock
        if: always()
        uses: DND-IT/action-lock@v0
        with:
          action: release
          lock_name: semantic-release
          token: ${{ secrets.GITHUB_TOKEN }}
```

### Locking a Dev Environment to a Pull Request

Hold a lock for the entire lifetime of a PR — the dev environment is exclusively yours until the PR is closed. The main branch workflow waits for the lock before applying.

**PR workflow** (`.github/workflows/terraform-pr.yaml`):

```yaml
on:
  pull_request:
    types: [opened, synchronize, reopened]
    paths: ['terraform/**']

jobs:
  apply:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: acquire dev lock
        uses: DND-IT/action-lock@v0
        with:
          action: acquire
          lock_name: terraform-dev
          stale_threshold: 0  # never expire — held until PR closes
          token: ${{ secrets.GITHUB_TOKEN }}

      - name: terraform apply
        run: terraform apply -auto-approve
```

**PR cleanup workflow** (`.github/workflows/terraform-pr-cleanup.yaml`):

```yaml
on:
  pull_request:
    types: [closed]
    paths: ['terraform/**']

jobs:
  unlock:
    runs-on: ubuntu-latest
    steps:
      - name: release dev lock
        uses: DND-IT/action-lock@v0
        with:
          action: release
          lock_name: terraform-dev
          token: ${{ secrets.GITHUB_TOKEN }}
```

**Main branch workflow** (`.github/workflows/terraform-deploy.yaml`):

```yaml
on:
  push:
    branches: [main]
    paths: ['terraform/**']

jobs:
  apply:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: acquire dev lock
        id: lock
        uses: DND-IT/action-lock@v0
        with:
          action: acquire
          lock_name: terraform-dev
          timeout: 0
          fail_on_timeout: false
          token: ${{ secrets.GITHUB_TOKEN }}

      - name: terraform apply
        if: steps.lock.outputs.acquired == 'true'
        run: terraform apply -auto-approve

      - name: release lock
        if: steps.lock.outputs.acquired == 'true'
        uses: DND-IT/action-lock@v0
        with:
          action: release
          lock_name: terraform-dev
          token: ${{ secrets.GITHUB_TOKEN }}
```

## Development

This action is written in Go and runs as a Docker container.

### Building

```bash
go build -o action-lock ./cmd/action-lock
```

### Testing

```bash
go test -v -race ./...
```

## Releases

This action uses [release-please](https://github.com/googleapis/release-please) for automated versioning based on conventional commits.

### Version Aliases

The release workflow automatically updates version aliases:
- `v0` — Always points to the latest v0.x.x release
- `v0.1` — Always points to the latest v0.1.x release

```yaml
- uses: DND-IT/action-lock@v0        # Always gets latest v0.x.x
- uses: DND-IT/action-lock@v0.1      # Always gets latest v0.1.x
- uses: DND-IT/action-lock@v0.1.0    # Pinned to specific version
```

## License

MIT
