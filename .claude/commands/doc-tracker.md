---
description: Generate and maintain comprehensive documentation from the code
argument-hint: [--api or --readme]
allowed-tools: Bash(ls:*), Bash(cat:*), Bash(test:*), Bash(grep:*), Bash(find:*)
---

Generate and maintain documentation from the code, keeping it in sync with implementation.

## Usage Examples

**Basic documentation generation:**

```
/doc-tracker
```

**Generate README:**

```
/doc-tracker --readme
```

**Check documentation coverage:**

```
/doc-tracker --check
```

**Help and options:**

```
/doc-tracker --help
```

## Implementation

If $ARGUMENT contains "help" or "--help":
Display this usage information and exit.

Parse documentation options from $ARGUMENTS (--readme, --help)

## 1. Analyze the current documentation

Check existing documentation
!find . -name "_.md" | grep -v node_modules | head -20
!test -f README.md && echo "README exists" || echo "No README.md found"
!find . -name "_.go" -exec grep -l '"""' {} \; | wc -l

Think step by step about documentation needs and:

1. Identify what documentation is missing
2. Generate appropriate documentation based on code analysis
3. Create templates for missing documentation
4. Ensure examples are included where helpful

Generate documentation in this format:

For README.md:

```markdown
# Project Name

Brief description of the project

## Installation

## Usage
```
