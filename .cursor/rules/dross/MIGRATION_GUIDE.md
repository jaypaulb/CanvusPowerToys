# Migration Guide: Old Numbered Rules → Atomic Structure

This guide helps you migrate from the old numbered rule structure to the new atomic design-based structure.

## Overview

**Old Structure** (`.cursor/rules/`):
- Numbered files: `1.1-core.mdc`, `2.1-project-management.mdc`, etc.
- Flat structure
- Some duplication between rules

**New Structure** (`.cursor/rules-atomic/`):
- Organized by category: `core/`, `languages/`, `workflows/`, `project/`
- Atomic design principles
- Token-optimized
- TDD integrated
- No duplication

## Migration Mapping

### Core Rules

| Old Rule | New Rule | Notes |
|----------|----------|-------|
| `1.1-core.mdc` | `core/command-execution.mdc`<br>`core/file-operations.mdc`<br>`core/quality-gates.mdc`<br>`core/error-handling.mdc`<br>`core/naming-conventions.mdc`<br>`core/tdd-methodology.mdc`<br>`core/code-structure.mdc` | Split into focused core rules |
| `1.2-dev-env.mdc` | `core/dev-env.mdc` | Updated with atomic design |
| `1.3-project-structure.mdc` | `core/project-structure.mdc` | Updated with atomic design |

### Project Management

| Old Rule | New Rule | Notes |
|----------|----------|-------|
| `2.1-project-management.mdc` | `project/project-management.mdc` | Updated with TDD requirements |
| `2.2-tasks-sop.mdc` | `project/tasks.mdc` | Updated with TDD requirements |
| `2.3-issues-sop.mdc` | `project/issues.mdc` | Updated structure |

### Workflows

| Old Rule | New Rule | Notes |
|----------|----------|-------|
| `3.1-local-agent-workflow.mdc` | `workflows/local-agent-workflow.mdc` | Updated with TDD workflow steps |
| `3.2-background-agent-workflow.mdc` | `workflows/background-agent-workflow.mdc` | Updated with TDD requirements |

### Language Standards

| Old Rule | New Rule | Notes |
|----------|----------|-------|
| `4.1-python-standards.mdc` | `languages/python.mdc` | Updated with TDD patterns and atomic design |
| `4.2-golang-standards.mdc` | `languages/golang.mdc` | Updated with TDD patterns and atomic design |
| `4.3-mern-standards.mdc` | `languages/mern.mdc` | Updated with TDD patterns and atomic design |
| `4.4-cpp-standards.mdc` | `languages/cpp.mdc` | Updated with TDD patterns and atomic design |

**New Language Rules:**
- `languages/flutter.mdc` - Flutter/Dart development (NEW)
- `languages/canvus.mdc` - Canvus API development (NEW, language-agnostic)

### Specialized Workflows

| Old Rule | New Rule | Notes |
|----------|----------|-------|
| `5.1-refactor.mdc` | `workflows/specialized/refactor.mdc` | Updated with TDD for refactoring |
| `5.2-refresh.mdc` | `workflows/specialized/refresh.mdc` | Updated structure |
| `5.3-reflect.mdc` | `workflows/specialized/reflect.mdc` | Updated structure |
| `5.4-research.mdc` | `workflows/specialized/research.mdc` | Updated structure |

## Key Differences

### 1. TDD Integration
**Old**: Testing mentioned but not enforced
**New**: TDD methodology always-applied, enforced in all rules

### 2. Atomic Design
**Old**: No structured code organization
**New**: Atomic design hierarchy (atoms → molecules → organisms → pages) enforced

### 3. Token Optimization
**Old**: Rules could exceed 1,000 tokens
**New**: Always-apply rules: ~728 tokens, context-specific: ~100-200 tokens

### 4. Rule Organization
**Old**: Flat structure with numbered files
**New**: Categorized structure (core/, languages/, workflows/, project/)

### 5. References
**Old**: Some duplication between rules
**New**: Rules reference core rules (tdd-methodology.mdc, code-structure.mdc) instead of duplicating

## Migration Steps

### Step 1: Backup Old Rules
```bash
# Backup old rules
cp -r .cursor/rules .cursor/rules-backup
```

### Step 2: Review New Structure
```bash
# Review new atomic structure
ls -la .cursor/rules-atomic/
tree .cursor/rules-atomic/
```

### Step 3: Compare Functionality
- Review `MIGRATION_GUIDE.md` (this file)
- Check that all functionality from old rules is present in new rules
- Note any customizations you made to old rules

### Step 4: Test New Rules (Optional)
```bash
# Copy atomic rules to test location
cp -r .cursor/rules-atomic .cursor/rules-test
# Test with sample project
```

### Step 5: Deploy New Rules (When Ready)
```bash
# Remove old rules (after backup)
rm -rf .cursor/rules/*.mdc

# Copy new structure
cp -r .cursor/rules-atomic/core .cursor/rules/
cp -r .cursor/rules-atomic/workflows .cursor/rules/
cp -r .cursor/rules-atomic/languages .cursor/rules/
cp -r .cursor/rules-atomic/project .cursor/rules/
```

## What's New

### New Core Rules
- `core/tdd-methodology.mdc` - Test-Driven Development methodology
- `core/code-structure.mdc` - Atomic design code structure

### New Language Rules
- `languages/flutter.mdc` - Flutter/Dart development
- `languages/canvus.mdc` - Canvus API development (language-agnostic)

### Enhanced Rules
- All language rules now include TDD patterns
- All workflow rules include TDD steps
- All rules enforce atomic design principles
- Quality gates enforce TDD requirements

## Breaking Changes

### None
The new rules are **additive** - they add functionality (TDD, atomic design) without removing existing functionality. All old patterns still work.

### Improvements
- Better organization
- Token optimization
- TDD enforcement
- Atomic design structure
- No duplication

## Rollback

If you need to rollback:
```bash
# Restore old rules
cp -r .cursor/rules-backup/* .cursor/rules/
```

## Support

If you encounter issues during migration:
1. Check `VALIDATION_REPORT.md` for validation results
2. Review `RULE_SELECTION_GUIDE.md` for rule selection
3. Check rule references in `.cursor/rules-atomic/README.md`

