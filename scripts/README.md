# Build Scripts

Utility scripts for building and maintaining Canvus PowerToys.

## Scripts

### `increment-version.sh`

Automatically increments the patch version number before building.

**Usage:**
```bash
./scripts/increment-version.sh
```

**What it does:**
- Reads current version from `internal/atoms/version/version.go`
- Increments patch version (e.g., `1.0.66` → `1.0.67`)
- Updates the version file
- Displays confirmation

**Example workflow:**
```bash
# Before building, increment version
./scripts/increment-version.sh
# 1.0.66 → 1.0.67

# Build for Windows
make build-windows
# Creates: canvus-powertoys.1.0.67.exe

# Build for Linux
make build-linux
# Creates: canvus-powertoys-linux (versioned internally)
```

**Version Format:**
- Uses semantic versioning: `MAJOR.MINOR.PATCH`
- Currently: `1.0.x`
- Only patches are incremented automatically
- To change MAJOR or MINOR, edit `internal/atoms/version/version.go` manually

**Why increment versions?**
- Each build gets a unique version number
- Makes it clear which binary has the latest changes
- Prevents confusion when testing multiple builds
- Version appears in:
  - Windows executable filename
  - Internal version tracking
  - Build metadata
  - Any about/help dialogs

## Typical Build Workflow

```bash
# 1. Make code changes
# (edit files...)

# 2. Increment version
./scripts/increment-version.sh

# 3. Build for target platform
make build-windows    # or make build-linux

# 4. Binary is ready for testing with new version number
ls -lh canvus-powertoys*.exe
# canvus-powertoys.1.0.67.exe
```

## Git Integration

The version script is safe for git workflows:
- Only modifies `internal/atoms/version/version.go`
- Changes should be committed together with code changes
- Example commit:
  ```bash
  ./scripts/increment-version.sh
  git add internal/atoms/version/version.go
  git commit -m "refactor: improve WebUI logging"
  ```

## Troubleshooting

**"File not found" error:**
- Run from project root directory
- Make sure you're in `/home/jaypaulb/Documents/gh/CanvusPowerToys/`

**Version didn't update:**
- Check that `internal/atoms/version/version.go` exists
- Verify file permissions
- Try running with `bash` explicitly: `bash scripts/increment-version.sh`

**Want to manually set version?**
Edit `internal/atoms/version/version.go` directly and change the Version string:
```go
Version = "1.2.3"
```

## Version History

See version changes in git log:
```bash
git log --oneline -- internal/atoms/version/version.go
```
