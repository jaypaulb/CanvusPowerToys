# Building with Version Management

Quick reference for building Canvus PowerToys with automatic version increments.

## One-Command Build Flow

```bash
# Increment version and build for Windows
./scripts/increment-version.sh && make build-windows

# Increment version and build for Linux
./scripts/increment-version.sh && make build-linux
```

## Step-by-Step

### Step 1: Increment Version
```bash
./scripts/increment-version.sh
```

Output:
```
✓ Version incremented successfully
  1.0.67 → 1.0.68

File: internal/atoms/version/version.go
New version line:
	Version     = "1.0.68"
```

### Step 2: Build
```bash
# Windows
make build-windows

# Linux
make build-linux
```

### Step 3: Verify
```bash
# Check the binary was created with new version
ls -lh canvus-powertoys*.exe

# Should see something like:
# canvus-powertoys.1.0.68.exe
```

## Why This Matters

**Before (Version stays the same):**
```
Old build:  canvus-powertoys.1.0.66.exe
New build:  canvus-powertoys.1.0.66.exe  ← Can't tell the difference!
```

**After (Version increments):**
```
Old build:  canvus-powertoys.1.0.66.exe
New build:  canvus-powertoys.1.0.68.exe  ← Clearly different!
```

## What Version Should I Use?

### Auto-Increment (Most Cases)
Use `increment-version.sh` for:
- Bug fixes
- Feature improvements
- UI changes
- Internal refactoring
- Testing builds

This increments the **patch version** (the last number):
- 1.0.66 → 1.0.67 → 1.0.68 → ...

### Manual Version Change (Rare)
Only edit `internal/atoms/version/version.go` directly if:
- You're making a major release (1.0 → 2.0)
- You're adding a significant feature set (1.0 → 1.1)

Example:
```go
// internal/atoms/version/version.go
Version = "2.0.0"  // Major version bump
```

## Complete Workflow Example

```bash
# Day 1: Fix CSS Preview
./scripts/increment-version.sh    # 1.0.66 → 1.0.67
make build-windows
# Test: canvus-powertoys.1.0.67.exe
# ✓ Looks good

# Day 2: Add WebUI logging
./scripts/increment-version.sh    # 1.0.67 → 1.0.68
make build-windows
# Test: canvus-powertoys.1.0.68.exe
# ✓ Logging works

# Day 3: Mark Generate Plugin as coming soon
./scripts/increment-version.sh    # 1.0.68 → 1.0.69
make build-windows
# Test: canvus-powertoys.1.0.69.exe
# ✓ Button is disabled

# Now you can see version progression:
ls -lh canvus-powertoys*.exe
# 1.0.67 - CSS Preview improvements
# 1.0.68 - WebUI logging system
# 1.0.69 - UI status updates
```

## Checking Version Info

**Current version in code:**
```bash
grep "Version" internal/atoms/version/version.go
```

**Version history:**
```bash
git log --oneline -- internal/atoms/version/version.go
```

**Binary version (Windows):**
```bash
# Check file properties (Windows)
# Right-click → Properties → Details

# Or extract from binary (Linux/Mac)
strings canvus-powertoys.1.0.68.exe | grep "1.0.68"
```

## Commit Messages

After building, include version in commit:

```bash
./scripts/increment-version.sh

# Commit your changes
git add -A
git commit -m "feat: add WebUI logging (v1.0.68)"

# Or describe the changes
git commit -m "docs: improve documentation
- Add version management workflow
- Create increment-version script
- Update build instructions
(v1.0.68)"
```

## Troubleshooting

### Script fails to execute
```bash
# Make it executable
chmod +x scripts/increment-version.sh

# Run with explicit bash
bash scripts/increment-version.sh
```

### Want to undo a version increment
```bash
# Option 1: Git restore
git checkout internal/atoms/version/version.go

# Option 2: Edit manually
nano internal/atoms/version/version.go
```

### Check if version was actually incremented
```bash
# Look at file
cat internal/atoms/version/version.go | grep "Version"

# Or check git diff
git diff internal/atoms/version/version.go
```

## Integration with CI/CD

If you set up automated builds, use the script in your pipeline:

```yaml
# Example GitHub Actions workflow
- name: Increment Version
  run: ./scripts/increment-version.sh

- name: Build
  run: make build-windows

- name: Create Release
  uses: actions/create-release@v1
  with:
    # Uses version from version.go
```

## Files Related to Versioning

- `internal/atoms/version/version.go` - Version constants
- `scripts/increment-version.sh` - Auto-increment script
- `scripts/README.md` - Script documentation
- `VERSION_WORKFLOW.md` - This guide
- `BUILD.md` - Build instructions
- `CLAUDE.md` - Project overview

## Quick Reference

| Task | Command |
|------|---------|
| Check current version | `grep "Version" internal/atoms/version/version.go` |
| Increment version | `./scripts/increment-version.sh` |
| Build Windows | `make build-windows` |
| Build Linux | `make build-linux` |
| Both steps at once | `./scripts/increment-version.sh && make build-windows` |
| Undo version change | `git checkout internal/atoms/version/version.go` |

---

**Key Takeaway:** Always run `./scripts/increment-version.sh` before building, and your version numbers will automatically tell you which binary has the latest changes!
