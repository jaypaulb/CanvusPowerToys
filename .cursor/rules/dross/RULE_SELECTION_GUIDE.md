# Rule Selection Guide

This guide helps you select the appropriate rules for different project types and scenarios.

## Quick Reference

### Always-Applied Rules (No Selection Needed)
These rules are always active:
- `core/command-execution.mdc`
- `core/file-operations.mdc`
- `core/quality-gates.mdc`
- `core/error-handling.mdc`
- `core/naming-conventions.mdc`
- `core/tdd-methodology.mdc`
- `core/code-structure.mdc`
- `workflows/local-agent-workflow.mdc`

### Context-Specific Rules (Auto-Applied Based on Files)
These rules apply automatically when you work with matching files:

**Language Rules** (apply when editing matching file types):
- `languages/python.mdc` - `*.py, *.pyi, requirements*.txt, pyproject.toml, setup.py`
- `languages/golang.mdc` - `*.go, go.mod, go.sum, Dockerfile`
- `languages/mern.mdc` - `*.js, *.jsx, *.ts, *.tsx, package.json, *.json`
- `languages/cpp.mdc` - `*.cpp, *.hpp, *.h, *.cc, CMakeLists.txt, Makefile, *.cmake`
- `languages/flutter.mdc` - `*.dart, pubspec.yaml`
- `languages/canvus.mdc` - `*canvus*.py, *canvus*.ts, *canvus*.js, *canvus*.go, canvus_api/**/*, canvus-client/**/*`

**Project Rules** (apply when editing matching files):
- `project/tasks.mdc` - `TASKS.md`
- `project/project-management.mdc` - `PRD.md, README.md, CHANGELOG.md`
- `project/issues.mdc` - `.gitlab-ci.yml, .github/**`

**Workflow Rules** (apply when context matches):
- `workflows/background-agent-workflow.mdc` - For autonomous agent operation
- `workflows/specialized/refactor.mdc` - When performing refactoring
- `workflows/specialized/refresh.mdc` - When debugging persistent issues
- `workflows/specialized/reflect.mdc` - For post-development analysis
- `workflows/specialized/research.mdc` - `docs/PRD.md, docs/TECH_STACK.md, docs/TASKS.md`

## Project Type Recommendations

### Python Projects

**Automatically Applied:**
- All core rules (always-apply)
- `languages/python.mdc` (when editing `.py` files)
- `project/tasks.mdc` (when editing `TASKS.md`)
- `project/project-management.mdc` (when editing `README.md`, `PRD.md`, `CHANGELOG.md`)

**Recommended for:**
- Web applications (Django, Flask, FastAPI)
- Data science projects
- Automation scripts
- API development

**Example Structure:**
```
src/
├── atoms/          # Basic building blocks (<100 lines)
│   ├── user.py
│   └── validator.py
├── molecules/      # Composed components (<200 lines)
│   └── user_validator.py
└── organisms/      # Complex modules (<500 lines)
    └── auth_system.py

tests/
├── atoms/          # Unit tests
├── molecules/      # Integration tests
└── organisms/      # System tests
```

### Go Projects

**Automatically Applied:**
- All core rules (always-apply)
- `languages/golang.mdc` (when editing `.go` files)
- `project/tasks.mdc` (when editing `TASKS.md`)

**Recommended for:**
- CLI tools
- Microservices
- High-performance applications
- System utilities

**Example Structure:**
```
internal/
├── atoms/          # Basic building blocks
│   ├── user/
│   └── validator/
├── molecules/      # Composed components
│   └── uservalidator/
└── organisms/      # Complex modules
    └── authsystem/
```

### MERN Stack Projects

**Automatically Applied:**
- All core rules (always-apply)
- `languages/mern.mdc` (when editing `.js`, `.jsx`, `.ts`, `.tsx` files)
- `project/tasks.mdc` (when editing `TASKS.md`)

**Recommended for:**
- Full-stack web applications
- React frontends with Express backends
- MongoDB-based applications

**Example Structure:**
```
client/src/
├── atoms/          # React components/hooks
├── molecules/      # Composed components
└── organisms/      # Complex components

server/src/
├── atoms/          # Models/utilities
├── molecules/      # Services
└── organisms/      # Controllers
```

### C++ Projects

**Automatically Applied:**
- All core rules (always-apply)
- `languages/cpp.mdc` (when editing `.cpp`, `.hpp`, `.h` files)
- `project/tasks.mdc` (when editing `TASKS.md`)

**Recommended for:**
- System programming
- Game development
- Performance-critical applications
- Embedded systems

### Flutter Projects

**Automatically Applied:**
- All core rules (always-apply)
- `languages/flutter.mdc` (when editing `.dart` files)
- `project/tasks.mdc` (when editing `TASKS.md`)

**Recommended for:**
- Cross-platform mobile apps
- Desktop applications
- Web applications with Flutter

### Canvus API Client Projects

**Automatically Applied:**
- All core rules (always-apply)
- `languages/canvus.mdc` (when editing Canvus API client files)
- Language-specific rule (Python, Go, TypeScript, etc.) based on file type

**Recommended for:**
- Building Canvus API clients
- Integrating with Canvus API
- Canvus API automation

## Workflow Selection

### For Refactoring Work
When performing large-scale refactoring:
- `workflows/specialized/refactor.mdc` applies automatically
- Includes TDD requirements for refactoring
- Enforces atomic design during refactoring

### For Debugging Persistent Issues
When hunting bugs:
- `workflows/specialized/refresh.mdc` applies automatically
- Systematic debugging approach
- Evidence-based problem solving

### For Project Inception
When starting new projects:
- `workflows/specialized/research.mdc` applies when working with `docs/PRD.md`, `docs/TECH_STACK.md`, `docs/TASKS.md`
- Interactive requirements gathering
- Template generation

### For Autonomous Agents
When running background agents:
- `workflows/background-agent-workflow.mdc` applies automatically
- Multi-agent coordination
- Autonomous task selection

## Customization

### Adjusting File Patterns
Update `globs` in rule frontmatter to match your project:

```markdown
---
globs: src/**/*.py,tests/**/*.py,*.py
alwaysApply: false
---
```

### Adding Custom Rules
Create new rules following the same structure:
1. Place in appropriate directory (`languages/`, `project/`, `workflows/`)
2. Set appropriate `globs` or `alwaysApply`
3. Reference core rules where appropriate
4. Keep focused and minimal

## Token Usage by Scenario

### Typical Python Development Session
- Always-apply: ~728 tokens
- Python rule: ~180 tokens
- Project rules: ~100 tokens
- **Total: ~1,008 tokens** (slightly over, but acceptable)

### Typical Go Development Session
- Always-apply: ~728 tokens
- Go rule: ~200 tokens
- Project rules: ~100 tokens
- **Total: ~1,028 tokens** (acceptable)

### Typical MERN Development Session
- Always-apply: ~728 tokens
- MERN rule: ~200 tokens
- Project rules: ~100 tokens
- **Total: ~1,028 tokens** (acceptable)

**Note**: Token usage may vary slightly, but all scenarios stay well within reasonable limits.

