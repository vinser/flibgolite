# Refactoring Workflow Rules

These rules ensure safe, incremental refactoring without breaking the project.

## Scope Control
- Modify only files relevant to the current refactoring step.
- Do not perform unrelated optimizations or stylistic changes.
- Keep each refactoring step small and isolated.

## Build Safety
- After moving files, always update imports.
- After each step, ensure the project builds:
  `go build ./...`
- Preserve existing behavior unless explicitly instructed.

## Branching Workflow
Follow the long-running refactor flow described in `REFACTORING_FLOW.md`:
- All refactoring goes into `refactor/architecture`.
- Bugfixes go into `master` and are merged back into the refactor branch.

## Cline Execution Rules
- Never modify files outside the scope of the current task.
- Do not rewrite large parts of the project unless explicitly requested.
- Maintain consistency with existing module structure.
