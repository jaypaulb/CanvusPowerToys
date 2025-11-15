# CursorRules Repository - Product Requirements Document

## Project Overview
Create a comprehensive, modular set of Cursor Rules that can be cloned to any new project for rapid development setup. Each rule should be focused, concise (~100 lines), and address specific aspects of the development workflow.

## Goals
- **Reusability**: Rules that work across different project types and environments
- **Modularity**: Each rule addresses a specific concern with minimal overlap
- **Clarity**: Concise, actionable guidelines that AI agents can follow autonomously
- **Completeness**: Cover all major aspects of software development workflow

## Target Users
- Developers using Cursor with AI pair programming
- Background AI agents working autonomously on tasks
- Teams needing consistent development standards across projects

## Core Requirements

### Rule Structure
- Maximum ~100 lines per rule file
- Clear, focused scope for each rule
- Consistent markdown formatting with frontmatter
- Actionable guidelines, not theoretical concepts

### Rule Categories

#### 1. Core Foundation Rules
- **1.1 core.mdc**: Essential AI agent behavior and decision-making principles
- **1.2 dev-env.mdc**: Environment detection and OS-specific best practices (Windows/Ubuntu)
- **1.3 project-structure.mdc**: Standardized project folder structure and organization

#### 2. Project Management & Workflow
- **2.1 project-management.mdc**: High-level project management standards
- **2.2 tasks-sop.mdc**: TASKS.md file management and task lifecycle
- **2.3 issues-sop.mdc**: Issue tracking and GitLab/GitHub integration

#### 3. Agent Collaboration Workflows
- **3.1 local-agent-workflow.mdc**: Guidelines for AI agents working with humans
- **3.2 background-agent-workflow.mdc**: Autonomous agent behavior for unattended tasks

#### 4. Language-Specific Standards
- **4.1 python-standards.mdc**: Python development best practices
- **4.2 golang-standards.mdc**: Go development standards and patterns
- **4.3 mern-standards.mdc**: MERN stack (MongoDB, Express, React, Node.js) guidelines
- **4.4 cpp-standards.mdc**: C++ development standards and best practices

#### 5. Specialized Development Workflows
- **5.1 refactor.mdc**: Guidelines for large-scale code refactoring
- **5.2 refresh.mdc**: Bug hunting and problem-solving approaches
- **5.3 reflect.mdc**: Post-development analysis and rule refinement
- **5.4 research.mdc**: Project inception, research, and design methodology

## Success Criteria
- Each rule file is self-contained and actionable
- Rules can be mixed and matched based on project needs
- AI agents can operate autonomously following the guidelines
- Development quality and consistency improves across projects
- Onboarding time for new projects is reduced significantly

## Technical Requirements
- Compatible with Cursor IDE rule system
- Proper frontmatter for rule targeting (globs, alwaysApply, etc.)
- Clear separation of concerns between rules
- No duplicate content across rules

## Deliverables
1. Refactored rule files in `.cursor/rules/` directory
2. Comprehensive README.md with usage instructions
3. This PRD documenting the vision and requirements

## Usage Model
This repository serves as a **template** that developers clone and copy the `.cursor/rules/` folder to their new projects. The rules then provide consistent development standards and AI agent behavior across all projects.

## Maintenance
- Rules should be updated based on real-world usage feedback
- Regular review of rule effectiveness and clarity
- Version control for rule changes to track improvements
