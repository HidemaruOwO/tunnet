**LANGUAGE**: Japanese (日本語)

## Workflow

1. Plan First
   - Create and tell me your plan and wait for “approve plan” I says. And call `plan tools mcp` built in OpenCode or codex or claude code if user says "approve".

2. Branching
   - From `main`: `feature/<short>` or `fix/<short>` ..etc (Refer `Branch Naming` section).
3. Tests
   - Add/update unit tests in `go language best practice location` for every non-trivial function.
   - If tests fail, reject your PR (`/abort`).
4. Quality
   - Run linters/formatters and commit resulting fixes.
   - Run build and tests locally before pushing.

## (IMPORTANT) Coding Best Practices

- Keep logic small and focused; split into multiple functions only if it improves clarity or reusability.
- Prefer early returns; avoid using `else` when not necessary.
- Handle errors explicitly; avoid `panic` except in initialization or unrecoverable situations.
- Limit `defer` usage scope to what is truly needed.
- Avoid using `any` or empty interfaces (`interface{}`); prefer concrete types or generics.
- Use short and descriptive variable names (e.g., `ctx`, `r`, `w`, `n`, `i`).
- Favor the Go standard library first (`os`, `io`, `context`, `net/http`, `encoding/*`, etc.).
- For file I/O, prefer `os.ReadFile`, `os.WriteFile`, and similar standard APIs.
- Pass `context.Context` as the first argument when cancellation or timeouts may be required.
- Use synchronization (`sync.Mutex`, `sync.WaitGroup`, channels) to protect shared state.
- Add comments for all exported symbols; comments should start with the identifier name.
- Minimize reflection usage.
- Wrap errors using `%w` and propagate them properly.
- Use structured logging where possible (`log/slog` or similar).

## Commands

- **Linter:** `golangci-lint run --issues-exit-code=1`
- **Formatter:** `gofmt -s -w . && goimports -w .`
- **Build:** `go build ./...`
- **Test:** `go test -race -cover ./...`

## Branch Naming

- feature: new feature
- fix: bug fix
- docs: documentation only
- refactor: code change without feature change
- test: adding/updating tests
- chore: build process or auxiliary tools

## Review Checklist

- [ ] Plan approved
- [ ] New branch from `main`
- [ ] One feature/bug-fix
- [ ] Tests added/updated
- [ ] Lint/format pass
- [ ] Build/test pass
- [ ] Docs updated if needed

## MCP (Optional)

- `@context7 doc "<lib>@<ver> <topic>"` — fetch latest docs. In answers, include
  a 7-sentence summary, a minimal canonical snippet, and the official URL. See
  “System Prompt Defaults” for when to use.
  Always use context7 when to generation codes, setup or configuration steps, or
  library/API documentation. This means you should automatically use the Context7 MCP
  tools to resolve library id and get library docs without me having to explicitly ask.
- `@fetch url` - supply `url` (with optional `headers`);
  When ask have urls, always use `@fetch` MCP tool to get the content of the url and use it to answer.
