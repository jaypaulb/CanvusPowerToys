# Version Management Workflow

This guide explains how to manage versions when building Canvus PowerToys with your changes.

## Quick Start

**Before each build, increment the version:**

```bash
./scripts/increment-version.sh
make build-windows
```

That's it! Your binary will now have a new version number showing it contains the latest changes.

## The Problem (What We're Solving)

Without version increment:
- `canvus-powertoys.1.0.66.exe` (old build from yesterday)
- `canvus-powertoys.1.0.66.exe` (new build with your changes)
- ❌ **Can't tell which is which!**

With version increment:
- `canvus-powertoys.1.0.66.exe` (old build from yesterday)
- `canvus-powertoys.1.0.67.exe` (new build with your changes)
- ✅ **Instantly clear which is newer**

## How It Works

### 1. Increment Version

```bash
./scripts/increment-version.sh
```

This automatically:
- Reads current version: `1.0.66`
- Increments patch: `1.0.67`
- Updates `internal/atoms/version/version.go`
- Shows you the result

### 2. Build

```bash
# Windows
make build-windows
# Output: canvus-powertoys.1.0.67.exe

# Linux
make build-linux
# Version embedded in binary internally
```

### 3. Done!

Your binary is ready with the new version number.

## Typical Development Flow

```bash
# 1. Make code changes
# Edit files, test locally, etc...

# 2. Ready to build? Increment version
./scripts/increment-version.sh
# 1.0.66 → 1.0.67

# 3. Build for testing
make build-windows

# 4. Test the new version
# canvus-powertoys.1.0.67.exe

# 5. If more changes needed, repeat from step 1
./scripts/increment-version.sh
# 1.0.67 → 1.0.68
make build-windows
```

## Git Integration

The version file should be committed with your changes:

```bash
# Make your code changes
git add -A
./scripts/increment-version.sh

# Now commit together
git add internal/atoms/version/version.go
git commit -m "feat: add new feature and bump version to 1.0.67"

# Or if already staged:
git add .
git commit -m "feat: add new feature"
```

## Version Numbering

**Format:** `MAJOR.MINOR.PATCH`

Currently: `1.0.x`

- **MAJOR** (1): Only change for major breaking changes (very rarely)
- **MINOR** (0): For feature additions (change manually in version.go if needed)
- **PATCH** (x): Auto-incremented by `increment-version.sh` for bug fixes and improvements

To manually change MAJOR or MINOR:

```bash
# Edit file directly
nano internal/atoms/version/version.go

# Change Version line to whatever you want:
# Version = "2.0.0"
# or
# Version = "1.1.0"
```

## Checking Versions

**See current version:**
```bash
grep "Version" internal/atoms/version/version.go
```

**See version history:**
```bash
git log --oneline -- internal/atoms/version/version.go
```

**See all version files created:**
```bash
ls -lh canvus-powertoys*.exe
```

## Troubleshooting

### Script not found
Make sure you're in the project root directory:
```bash
cd /home/jaypaulb/Documents/gh/CanvusPowerToys
./scripts/increment-version.sh
```

### Permission denied
Make the script executable:
```bash
chmod +x scripts/increment-version.sh
```

### Version didn't change
Run the script with bash explicitly:
```bash
bash scripts/increment-version.sh
```

### Want to undo a version increment
```bash
# Edit the file manually
nano internal/atoms/version/version.go

# Or use git to revert
git checkout internal/atoms/version/version.go
```

## Benefits

✅ **Clear version tracking** - Know exactly which binary has which changes
✅ **Easy testing** - Compare different builds side by side
✅ **Git-friendly** - Works perfectly with git commits
✅ **Automatic** - One command to increment
✅ **Semantic versioning** - Standard format that tools understand
✅ **Build metadata** - Version embeds into the Windows executable

## Real Example

```bash
# Session 1: Fix CSS Preview
./scripts/increment-version.sh
# 1.0.66 → 1.0.67
make build-windows
# Created: canvus-powertoys.1.0.67.exe
# Test it...

# Session 2: Add WebUI logging
./scripts/increment-version.sh
# 1.0.67 → 1.0.68
make build-windows
# Created: canvus-powertoys.1.0.68.exe
# Test it...

# Session 3: Mark Generate Plugin as coming soon
./scripts/increment-version.sh
# 1.0.68 → 1.0.69
make build-windows
# Created: canvus-powertoys.1.0.69.exe
# Test it...

# Now you have clear versions of each build!
ls -lh canvus-powertoys*.exe
# 1.0.67 - CSS Preview fix
# 1.0.68 - WebUI logging
# 1.0.69 - Generate Plugin status
```

---

**Summary:** Run `./scripts/increment-version.sh` before each build, and you'll always know which binary contains which changes!
