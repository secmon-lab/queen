# Queen

Agentic SAST triage tool. Takes SARIF findings and uses LLM to determine true/false positives.

## Build & Test

```bash
go vet ./...
go test ./...
```

## Architecture

DDD + Clean Architecture. CLI -> UseCase -> Service -> Domain.

- `pkg/cli/` — CLI commands (urfave/cli/v3)
- `pkg/cli/config/` — LLM provider configuration (Anthropic, OpenAI, Gemini)
- `pkg/domain/types/` — Type-safe IDs and enums
- `pkg/domain/model/` — Domain models (Finding, Issue, TriageResult, ScanSession)
- `pkg/domain/interfaces/` — Repository interface
- `pkg/repository/fs/` — Filesystem-based repository
- `pkg/service/sarif/` — SARIF parser
- `pkg/service/triage/` — LLM triage service (gollem Agent)
- `pkg/tool/` — Agent tools (read_file, find, grep). One file per tool, implements `gollem.Tool`
- `pkg/usecase/` — Use cases with functional options for DI

## Git & PR

- Commit messages: English, one line, semantic commit format (`feat:`, `fix:`, `refactor:`, `test:`, `docs:`, `chore:`, `ci:`)
- Example: `feat: add SARIF parser for semgrep output`
- Do NOT add `Co-Authored-By` or any trailer lines
- PR title and description: always in English

## Key Dependencies

- `github.com/m-mizutani/gollem` — LLM abstraction (Agent, Tool, Session)
- `github.com/m-mizutani/goerr/v2` — Error wrapping
- `github.com/urfave/cli/v3` — CLI framework
- `github.com/m-mizutani/gt` — Test assertions
