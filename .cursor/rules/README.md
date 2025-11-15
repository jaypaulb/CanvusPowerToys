# CursorRules - Modular Development Standards

A comprehensive set of focused, reusable Cursor Rules for rapid development setup across different project types. Each rule is designed to be concise (~100 lines), focused on specific concerns, and actionable for both human developers and AI agents.

## Quick Start

### Method 1: GitHub Template (Recommended for New Projects)

1. **Use this template** on GitHub:
   - Click "Use this template" button on the repository page
   - Create a new repository from the template
   - Clone your new repository

2. **Copy the rules** to your project:
   ```bash
   cp -r JaypaulsCursorRules/.cursor your-project/
   ```

3. **Select appropriate rules** based on your project type (see Rule Selection Guide below)

4. **Customize frontmatter** in rule files to match your project's file patterns

5. **Start developing** with AI assistance following established standards

### Method 2: Git Subtree (Recommended for Rule Updates)

For projects where you want to keep rules synchronized with updates:

**Prerequisites:**
- Project must be a git repository (`git init` if needed)
- `.cursor/rules/` directory must not exist or be empty

1. **Add rules as subtree** (one-time setup):
   ```bash
   # Ensure you're in a git repository
   git init  # if not already initialized
   
   # Remove existing rules if they exist
   rm -rf .cursor/rules/
   
   # Add subtree
   git subtree add --prefix=.cursor/rules https://github.com/jaypaulb/JaypaulsCursorRules.git main --squash
   ```

2. **Update rules** when new versions are available:
   ```bash
   git subtree pull --prefix=.cursor/rules https://github.com/jaypaulb/JaypaulsCursorRules.git main --squash
   ```

3. **Select and customize rules** as needed for your project

### Method 3: Manual Clone (Simple Setup)

1. **Clone this repository** to your new project:
   ```bash
   git clone https://github.com/jaypaulb/JaypaulsCursorRules.git
   cp -r JaypaulsCursorRules/.cursor .
   rm -rf JaypaulsCursorRules  # Clean up
   ```

2. **Select appropriate rules** based on your project type (see Rule Selection Guide below)

3. **Customize frontmatter** in rule files to match your project's file patterns

4. **Start developing** with AI assistance following established standards

## Deployment Method Comparison

| Method | Best For | Pros | Cons |
|--------|----------|------|------|
| **GitHub Template** | New projects, one-time setup | Clean separation, no git history pollution | Manual updates required |
| **Git Subtree** | Projects needing rule updates | Version controlled, can pull latest changes | More complex setup |
| **Manual Clone** | Quick setup, simple projects | Simple, no dependencies | No update mechanism |

**Recommendation:**
- Use **GitHub Template** for most new projects
- Use **Git Subtree** if you want to keep rules synchronized with updates
- Use **Manual Clone** for quick prototyping or simple projects

## Rule Categories

### 1. Core Foundation Rules
- **1.1 core.mdc**: Essential AI agent behavior and decision-making principles
- **1.2 dev-env.mdc**: Environment detection and OS-specific best practices

### 2. Project Management & Workflow
- **2.1 project-management.mdc**: High-level project management standards
- **2.2 tasks-sop.mdc**: TASKS.md file management and task lifecycle
- **2.3 issues-sop.mdc**: Issue tracking and GitLab/GitHub integration

### 3. Agent Collaboration Workflows
- **3.1 local-agent-workflow.mdc**: Guidelines for AI agents working with humans
- **3.2 background-agent-workflow.mdc**: Autonomous agent behavior for unattended tasks

### 4. Language-Specific Standards
- **4.1 python-standards.mdc**: Python development best practices
- **4.2 golang-standards.mdc**: Go development standards and patterns
- **4.3 mern-standards.mdc**: MERN stack (MongoDB, Express, React, Node.js) guidelines

### 5. Specialized Development Workflows
- **5.1 refactor.mdc**: Guidelines for large-scale code refactoring
- **5.2 refresh.mdc**: Bug hunting and problem-solving approaches
- **5.3 reflect.mdc**: Post-development analysis and rule refinement

## Rule Selection Guide

### For Python Projects
```
✅ 1.1 core.mdc (always apply)
✅ 1.2 dev-env.mdc
✅ 4.1 python-standards.mdc
✅ 2.1 project-management.mdc
✅ 2.2 tasks-sop.mdc
✅ 2.3 issues-sop.mdc
⚠️ 3.1 local-agent-workflow.mdc (if working with AI agents)
```

### For Go Projects
```
✅ 1.1 core.mdc (always apply)
✅ 1.2 dev-env.mdc
✅ 4.2 golang-standards.mdc
✅ 2.1 project-management.mdc
✅ 2.2 tasks-sop.mdc
✅ 2.3 issues-sop.mdc
⚠️ 3.1 local-agent-workflow.mdc (if working with AI agents)
```

### For MERN Stack Projects
```
✅ 1.1 core.mdc (always apply)
✅ 1.2 dev-env.mdc
✅ 4.3 mern-standards.mdc
✅ 2.1 project-management.mdc
✅ 2.2 tasks-sop.mdc
✅ 2.3 issues-sop.mdc
⚠️ 3.1 local-agent-workflow.mdc (if working with AI agents)
```

### For Background Agent Projects
```
✅ 1.1 core.mdc (always apply)
✅ 1.2 dev-env.mdc
✅ 3.2 background-agent-workflow.mdc
✅ 2.2 tasks-sop.mdc
✅ 2.3 issues-sop.mdc
⚠️ 4.x Language-specific rule based on project
```

## Rule Structure

Each rule file follows this structure:

```markdown
---
description: Brief description of what this rule covers
globs: *.py,*.js  # File patterns this rule applies to
alwaysApply: true/false  # Whether to apply regardless of file type
---
# Rule Title

## Section 1: Core Concepts
- Clear, actionable guidelines
- Specific examples where helpful
- No theoretical discussions

## Section 2: Implementation Details
- Step-by-step procedures
- Command examples
- Code snippets

## Section 3: Quality Standards
- Measurable criteria
- Validation steps
- Common pitfalls to avoid
```

## Customization

### Adjusting File Patterns
Update the `globs` field in rule frontmatter to match your project structure:

```markdown
---
globs: src/**/*.py,tests/**/*.py,*.py
---
```

### Enabling/Disabling Rules
- Set `alwaysApply: true` for rules that should always be active
- Set `alwaysApply: false` and use specific `globs` for conditional application
- Remove or rename `.mdc` extension to disable rules

### Project-Specific Customizations
1. Copy rules to your project's `.cursor/rules/` directory
2. Modify content to match your team's specific needs
3. Add project-specific rules following the same structure
4. Update this README with your customizations

## Working with AI Agents

### Human + AI Development
Use **local-agent-workflow.mdc** when working interactively with AI agents. This rule provides:
- Guidelines for autonomous AI behavior
- Clear escalation criteria
- Progress tracking standards
- Quality gate requirements

### Background Agent Development
Use **background-agent-workflow.mdc** for unattended AI agents working through task queues. This rule provides:
- Multi-agent coordination protocols
- Conflict resolution strategies
- Autonomous task selection
- Self-healing behaviors

## Quality Assurance

### Rule Validation
Each rule is designed to be:
- **Actionable**: Contains specific, executable instructions
- **Focused**: Addresses one primary concern (~100 lines max)
- **Measurable**: Includes clear success criteria
- **Reusable**: Works across different projects of the same type

### Testing Rules
Validate rules by:
1. Testing with sample projects of each type
2. Ensuring AI agents can follow instructions autonomously
3. Verifying no conflicts between rule combinations
4. Confirming quality gates are achievable

## Contributing

### Adding New Rules
1. Follow the established structure and style
2. Keep rules focused and under 100 lines
3. Include specific examples and commands
4. Test with real projects before submitting
5. Update this README with the new rule

### Improving Existing Rules
1. Maintain backward compatibility when possible
2. Update all related rules consistently
3. Test changes with multiple project types
4. Document breaking changes in commit messages

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

- **Issues**: Report problems or suggest improvements via GitHub issues
- **Discussions**: Share experiences and ask questions in GitHub discussions
- **Documentation**: All rules are self-documenting with examples and procedures

---

**Version**: 1.0.0
**Last Updated**: January 2024
**Compatibility**: Cursor IDE with Cursor Rules support
