# Deployment Guide

## Overview

Rules are developed in `.cursor/rules-atomic/` for validation, but must be deployed to `.cursor/rules/` to work in Cursor.

## Important: Reference Paths

**All references in rule files point to `.cursor/rules/` (the live folder), not `rules-atomic/`.**

This means:
- ✅ Rules work immediately after copying to `.cursor/rules/`
- ✅ No manual path updates needed
- ✅ References are correct from the start

## Deployment Steps

### Option 1: Copy Structure (Recommended)

```bash
# Copy entire atomic structure to rules folder
cp -r .cursor/rules-atomic/* .cursor/rules/

# Or use rsync to preserve structure
rsync -av .cursor/rules-atomic/ .cursor/rules/
```

### Option 2: Symlink (Development)

```bash
# Create symlink for development (rules-atomic → rules)
ln -s .cursor/rules-atomic .cursor/rules-atomic-link
# Then copy structure when ready
```

### Option 3: Selective Deployment

Copy only what you need:
```bash
# Copy atoms
cp -r .cursor/rules-atomic/atoms .cursor/rules/

# Copy specific molecules
cp .cursor/rules-atomic/molecules/python-standards.mdc .cursor/rules/molecules/

# Copy organisms
cp -r .cursor/rules-atomic/organisms .cursor/rules/

# Copy pages
cp -r .cursor/rules-atomic/pages .cursor/rules/
```

## Directory Structure After Deployment

```
.cursor/rules/
├── atoms/
│   ├── command-execution.mdc
│   ├── quality-gates.mdc
│   ├── file-operations.mdc
│   ├── error-handling.mdc
│   ├── naming-conventions.mdc
│   ├── documentation-standards.mdc
│   ├── git-workflow.mdc
│   └── task-management.mdc
├── molecules/
│   └── python-standards.mdc
├── organisms/
│   ├── core-foundation.mdc
│   ├── workflow-standards.mdc
│   └── project-management.mdc
├── pages/
│   └── python-project.mdc
└── templates/
```

## Verification

After deployment, verify:
1. All rule files are in `.cursor/rules/` with proper structure
2. Cursor recognizes the rules (check Settings → Rules)
3. Rules apply correctly when working with matching files
4. References in rule files resolve correctly

## Rollback

If issues occur:
1. Keep original rules in `.cursor/rules/` (backup first)
2. Atomic rules are additive - can coexist
3. Remove atomic rules if needed: `rm -rf .cursor/rules/atoms .cursor/rules/molecules .cursor/rules/organisms .cursor/rules/pages`

## Notes

- Rules in `rules-atomic/` are for development/validation only
- All path references already point to `.cursor/rules/`
- No manual path updates needed after deployment
- Structure is preserved during copy

