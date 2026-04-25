# Cline Rules
# Minimal, strict rules for safe and controlled refactoring.

============================================================
1. CODE STYLE
============================================================

- All comments must be in English.
- All identifiers must use English names.
- Only English text is allowed in source files (except test fixtures).
- Comments must be short, technical, and neutral.
- Prefer explicit, readable code.
- Formatting and import order are handled by Go tools.

============================================================
2. REFACTORING SCOPE & SAFETY
============================================================

- Modify only files directly related to the current task.
- Do not perform unrelated optimizations or cosmetic changes.
- Each refactoring step must be small, isolated, and verifiable.
- The project must build after every step: `go build ./...`.
- Program behavior must remain unchanged unless explicitly allowed.
- Do not change architecture, package structure, or public APIs.
- Do not modify main.go or service.go unless explicitly requested.

============================================================
3. CLINE EXECUTION RULES
============================================================

- Do not modify files outside the declared task scope.
- Do not rewrite large sections of the project.
- Create helper functions only when explicitly required by the task.
- Do not change formatting beyond what Go tools produce.
- Do not modify comments unless the task requires it.

============================================================
4. NO UNREQUESTED CHANGES
============================================================

- Do not rename identifiers unless required by the task.
- Do not reorder code unless required for correctness.
- Do not introduce new abstractions unless required by the task.
