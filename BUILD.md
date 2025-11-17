# Build Instructions

This document provides detailed instructions for building Canvus PowerToys from source.

## Prerequisites

### Required
- **Go 1.24.5 or later** - [Download Go](https://go.dev/dl/)
- **Git** - For version control and dependency management

### For Windows Cross-Compilation (from Linux)
- **gcc-mingw-w64-x86-64** - Windows cross-compilation toolchain
  ```bash
  sudo apt-get install gcc-mingw-w64-x86-64
  ```

### Optional
- **goversioninfo** - For Windows executable version info and icons
  ```bash
  go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest
  ```
- **golangci-lint** - For code linting (optional)
  ```bash
  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  ```

## Build Methods

### Method 1: Using Makefile (Recommended)

The Makefile provides convenient build targets:

```bash
# Build for current platform
make build

# Build for Linux
make build-linux

# Build for Windows (from Linux, requires mingw-w64)
make build-windows
```

### Method 2: Using Build Scripts

#### Linux Build
```bash
chmod +x build.sh
./build.sh
```

Output: `canvus-powertoys`

#### Windows Build (Cross-Compilation from Linux)
```bash
chmod +x build-windows.sh
./build-windows.sh
```

Output: `canvus-powertoys.X.X.X.exe` (version number included)

### Method 3: Manual Go Build

#### Linux Build
```bash
# Process WebUI assets first (minifies CSS, JS, HTML)
go run webui/build/process-assets.go

# Build
go build -ldflags="-s -w" -o canvus-powertoys ./cmd/powertoys
```

#### Windows Build (Cross-Compilation from Linux)
```bash
# Process WebUI assets first
go run webui/build/process-assets.go

# Generate version info resource (optional, requires goversioninfo)
goversioninfo -64 -o cmd/powertoys/resource.syso versioninfo.json

# Get version info
VERSION=$(grep 'Version.*=' internal/atoms/version/version.go | sed -n 's/.*"\([^"]*\)".*/\1/p')
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build with CGO enabled for Fyne
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build \
    -ldflags="-s -w -H windowsgui -X github.com/jaypaulb/CanvusPowerToys/internal/atoms/version.Version=$VERSION -X github.com/jaypaulb/CanvusPowerToys/internal/atoms/version.BuildDate=$BUILD_DATE -X github.com/jaypaulb/CanvusPowerToys/internal/atoms/version.GitCommit=$GIT_COMMIT" \
    -o canvus-powertoys.$VERSION.exe ./cmd/powertoys

# Clean up resource.syso (if generated)
rm -f cmd/powertoys/resource.syso
```

## Build Flags Explained

### Standard Flags
- `-ldflags="-s -w"`: Strip debug symbols and disable DWARF generation (reduces binary size)
- `-o <output>`: Output filename

### Windows-Specific Flags
- `-H windowsgui`: Sets Windows subsystem to GUI (hides console window)
  - To show console for debugging: Remove `-H windowsgui` from `-ldflags`
- `-X <package.Variable>=<value>`: Set build-time variables for version info
  - `Version`: Application version
  - `BuildDate`: Build timestamp (UTC)
  - `GitCommit`: Git commit hash

### Cross-Compilation Flags
- `CGO_ENABLED=1`: Enable CGO (required for Fyne framework)
- `GOOS=windows`: Target operating system
- `GOARCH=amd64`: Target architecture
- `CC=x86_64-w64-mingw32-gcc`: C compiler for Windows cross-compilation

## Version Info (Windows)

The Windows build can include version information and icons using `goversioninfo`:

1. Install goversioninfo:
   ```bash
   go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest
   ```

2. Generate resource file:
   ```bash
   goversioninfo -64 -o cmd/powertoys/resource.syso versioninfo.json
   ```

3. Build normally - Go will automatically include `resource.syso` if present

4. Clean up after build:
   ```bash
   rm -f cmd/powertoys/resource.syso
   ```

The `versioninfo.json` file contains:
- Application version
- Company name
- Product name
- Icon file path
- File description

## WebUI Asset Processing

Before building, WebUI assets (CSS, JS, HTML) must be processed (minified):

```bash
go run webui/build/process-assets.go
```

This step is automatically included in:
- `make build`
- `make build-linux`
- `make build-windows`

The processed assets are embedded in the binary at build time.

## Build Output

### Linux
- Output: `canvus-powertoys` (or `canvus-powertoys-linux` with Makefile)
- Executable: Yes (chmod +x may be required)
- Size: ~15-20 MB (stripped)

### Windows
- Output: `canvus-powertoys.X.X.X.exe` (version number included)
- Executable: Yes (.exe file)
- Size: ~20-25 MB (stripped, with version info)

## Troubleshooting

### Windows Cross-Compilation Issues

**Error: `x86_64-w64-mingw32-gcc: command not found`**
- Solution: Install mingw-w64: `sudo apt-get install gcc-mingw-w64-x86-64`

**Error: CGO compilation fails**
- Solution: Ensure `CGO_ENABLED=1` is set and `CC=x86_64-w64-mingw32-gcc` is specified

**Error: Fyne framework issues**
- Solution: Fyne requires CGO for Windows. Ensure CGO is enabled and mingw-w64 is installed

### Version Info Issues

**Warning: goversioninfo not found**
- This is not critical - build will succeed without version info
- To include version info: `go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest`

**Error: resource.syso not found**
- This is normal if goversioninfo is not installed
- Build will succeed without version info (no icon, no version details in file properties)

### Asset Processing Issues

**Error: Asset processing fails**
- Solution: Ensure `webui/build/process-assets.go` exists and is runnable
- Check that source assets exist in `webui/src/`

### Build Size Issues

**Binary is too large**
- Ensure `-ldflags="-s -w"` is included (strips debug symbols)
- Check that assets are minified (run `process-assets.go`)
- Verify no unnecessary dependencies are included

## Development Builds

For development, you may want to:
- Remove `-s -w` flags to keep debug symbols
- Remove `-H windowsgui` to show console for debugging
- Skip asset minification for faster iteration

Example development build:
```bash
# Process assets (or skip for faster iteration)
go run webui/build/process-assets.go

# Build with debug symbols and console
go build -o canvus-powertoys-dev ./cmd/powertoys
```

## Continuous Integration

For CI/CD pipelines:

```yaml
# Example GitHub Actions workflow
- name: Build Linux
  run: make build-linux

- name: Build Windows
  run: |
    sudo apt-get update
    sudo apt-get install -y gcc-mingw-w64-x86-64
    make build-windows
```

## Version Increment

To increment version before building:

1. Update version in `internal/atoms/version/version.go`
2. Optionally use version increment script (if available):
   ```bash
   ./dross/scripts/increment-version.sh
   ```
3. Build normally - version will be included in binary

## Testing Build

After building, verify the build:

```bash
# Linux
./canvus-powertoys --version  # If version flag exists

# Windows
canvus-powertoys.X.X.X.exe  # Run executable
```

## Clean Build Artifacts

To clean build artifacts:

```bash
make clean
```

Or manually:
```bash
rm -f canvus-powertoys canvus-powertoys.exe canvus-powertoys-linux
rm -f cmd/powertoys/resource.syso
```

## Additional Resources

- [Go Build Documentation](https://pkg.go.dev/cmd/go#hdr-Compile_packages_and_dependencies)
- [Fyne Framework Documentation](https://developer.fyne.io/)
- [goversioninfo Documentation](https://github.com/josephspurrier/goversioninfo)

