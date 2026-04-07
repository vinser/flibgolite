# Architecture Guidelines

These rules define the structural decisions of the project.

## Project Structure
Follow the modular layout:

- `internal/core` — models, config, errors, interfaces  
- `internal/parsers` — format parsers (fb2, epub, etc.)  
- `internal/store` — storage and metadata  
- `internal/index` — indexing logic
- `internal/converter` — book format conversion
- `internal/opds` — OPDS feed generation  
- `internal/app` — application orchestration  
- `internal/web` — embedded static web UI  
- `internal/admin` — embedded admin UI  

## Rules
- Do not create new top-level directories without explicit instruction.  
  *Why:* keeps the project predictable and maintainable.
- Keep modules small and cohesive.
- Avoid introducing new dependencies unless explicitly requested.
- Follow existing patterns in `internal/core/errors` for error handling.
