# Cline Rules for FLibGoLite  
# Strict, deterministic workflow for safe refactoring and architecture work.

============================================================  
1. CODING STANDARDS  
============================================================

## Language
- All code comments must be written in English.  
- All identifiers (packages, variables, functions, structs) must use English names.  
- Only English text is allowed inside source files (except test fixtures).

## Style
- Keep comments concise and technical.  
- Avoid conversational or emotional tone.  
- Prefer explicit, readable code over clever constructs.  
- Follow existing patterns in `/internal/core/` for naming and structure.

## Documentation
- All documentation must be in English.  
- Place new docs in the project root or `/docs/`.

============================================================  
2. ARCHITECTURE GUIDELINES  
============================================================

## Project Structure
Follow the modular layout:

- internal/core       — models, config, errors, interfaces  
- internal/parsers    — format parsers (fb2, epub, etc.)  
- internal/store      — storage and metadata  
- internal/index      — indexing logic (core indexing operations)  
- internal/converter  — book format conversion  
- internal/opds       — OPDS feed generation and HTTP handlers  
- internal/app        — application orchestration and lifecycle management  
- internal/service    — integration with kardianos/service  
- internal/web        — embedded static web UI  
- internal/admin      — embedded admin UI  

## Rules
- Do not create new **top-level** directories without explicit instruction.  
- It is allowed to create new subdirectories under `/internal/` when:
  - they are listed in “Project Structure”, or  
  - they are required by the current refactoring step.  
- Keep modules small and cohesive.  
- Avoid introducing new dependencies unless explicitly requested.  
- Follow existing patterns in `internal/core/errors` for error handling.  
- Preserve service behavior and shutdown semantics:
  - same kardianos/service integration behavior,  
  - same console mode behavior in Docker,  
  - graceful stop for OPDS server and indexer.

============================================================  
3. REFACTORING WORKFLOW RULES  
============================================================

## Scope Control
- Modify only files relevant to the current refactoring step.  
- Do not perform unrelated optimizations or stylistic changes.  
- Keep each refactoring step small and isolated.  
- When extracting logic from `cmd/flibgolite`:
  - prefer “move implementation to internal/* and keep a thin wrapper in cmd”.

## Build Safety
- After moving files, always update imports.  
- After each step, ensure the project builds: `go build ./...`  
- Preserve existing behavior unless explicitly instructed.

## Branching Workflow
Follow the long-running refactor flow described in `REFACTORING_FLOW.md`:
- All refactoring goes into `refactor/architecture`.  
- Bugfixes go into `master` and are merged back into the refactor branch.

## Cline Execution Rules
- Never modify files outside the scope of the current task.  
- Do not rewrite large parts of the project unless explicitly requested.  
- Maintain consistency with existing module structure.

============================================================  
4. PLAN / ACT MODEL BEHAVIOR  
============================================================

## Planning Model (Qwen2.5 Coder 32B)
- Provide analysis, reasoning, and step-by-step plans only.  
- Do not generate code or patches.  
- Keep responses concise and structured.  
- Never propose changes outside the current task scope.  
- When analyzing the project, output:
  - architecture overview  
  - dependency mapping  
  - refactoring opportunities  
  - risks and constraints  
- Do not propose code changes during analysis.

## Action Model (Qwen2.5 Coder 7B)
- Generate patches only when explicitly requested.  
- Keep patches minimal and isolated.  
- Do not modify unrelated files.  
- Do not rewrite large sections of code.  
- Always show a unified diff before applying changes.  
- Do not include explanations inside the patch.  
- After generating a patch, stop and wait for confirmation.

============================================================  
5. NO UNREQUESTED OPTIMIZATIONS  
============================================================

- Do not rename identifiers unless:
  - the rename is strictly required by the current refactoring step, and  
  - the new name follows existing naming patterns.  
- Do not reorder code unless required for correctness.  
- Do not change formatting beyond gofmt/goimports.  
- Do not introduce new abstractions or helper functions **except** when:
  - required by the current refactoring step, and  
  - directly related to the architecture plan (App, OPDS server, indexer, service integration).  
- Do not modify comments unless explicitly instructed.

============================================================  
6. FILE OPERATIONS  
============================================================

- When moving files, always update imports.  
- When creating new files, follow existing naming conventions.  
- Do not delete files unless explicitly instructed.  
- Creating new directories under `/internal/` is allowed when required by the current refactoring step or listed in the architecture structure.

============================================================  
7. GO LANGUAGE RULES  
============================================================

- Follow Go module boundaries.  
- Use explicit imports; avoid wildcard imports.  
- Maintain existing error-handling patterns.  
- Follow patterns in `internal/core/errors` for error creation.

