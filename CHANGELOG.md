# Changelog

All notable changes to gocode are documented here.

---

## v0.8.0 — The Universal Model Layer

### 200+ Model Support
- Universal model layer: 4 native providers + 7 proxy services via OpenAI-compatible shim
- Native: Anthropic (Claude), OpenAI (GPT), Google (Gemini), xAI (Grok)
- Proxy: DeepSeek, Mistral, Groq, Together AI, OpenRouter, Azure OpenAI
- Local: Ollama, LM Studio, vLLM via `OPENAI_BASE_URL`
- Codex backend integration with `~/.codex/auth.json` auth
- Full model list with aliases in [docs/supported-models.md](docs/supported-models.md)

### Provider Launch Profiles
- `gocode profile init` — create default profile
- `gocode profile auto --goal coding` — auto-detect best provider from env vars
- `gocode profile recommend --goal latency` — preview without saving
- `gocode profile show` — show current profile
- Goal-based model selection: `--goal coding`, `--goal latency`, `--goal balanced`

### Persistent Memory
- Cross-session memory system persisted to `.gocode/memory.json`
- `/memory set key value` — store a memory
- `/memory get key` — retrieve a memory
- `/memory list` — show all memories
- `/memory delete key` — remove a memory
- Memories injected into system prompt automatically

### Task Management
- Persistent task tracking in `.gocode/tasks.json`
- `/tasks add Fix the auth bug` — create a task
- `/tasks list` — show all tasks
- `/tasks done 1` — mark task complete

### Runtime Hardening
- `gocode smoke` — quick runtime smoke test (version, registries, skills, plugins)
- `gocode hardening` — security audit (file permissions, API key exposure in shell history)
- `gocode doctor` expanded to check all 10 provider env vars + Codex auth + profile

### New Tools
- WebFetchTool — fetch URL content (10KB max, no API key)
- WebSearchTool — DuckDuckGo instant answer API (no API key, no MCP required)
- NotebookEditTool — Jupyter notebook (.ipynb) editing
- `orchestrator_delegate` — delegate tasks to specialist sub-agents
- `orchestrator_delegate_bg` — background sub-agent delegation

### New Slash Commands
- `/commit` — auto-generate commit message with `Co-Authored-By: gocoder6969`
- `/memory` — manage persistent cross-session memories
- `/tasks` — manage persistent task list

### Multi-Agent Orchestration (Wired)
- Orchestrator delegation tools now registered in tool registry (previously dead code)
- LLM can call `orchestrator_delegate` and `orchestrator_delegate_bg`
- 4 built-in sub-agent profiles: coordinator, deep-worker, planner, debugger
- Permission-filtered tool access per sub-agent
- BackgroundManager with 5-slot concurrency limit

---

## v0.7.0 — The Agent Operating System

### Full Terminal UI (Bubbletea TUI)
- Split panels: chat on left, git diff viewer on right (Ctrl+D toggle)
- Tab to switch between Build mode (full access) and Plan mode (read-only)
- 4 built-in themes: `golang`, `monokai`, `dracula`, `nord`
- Custom themes via `.gocode/theme.json`
- Custom keybinds via `.gocode/keybinds.json`
- REPL is default; `--tui` flag for TUI mode

### Multi-Agent Orchestration
- Orchestrator with 4 built-in sub-agent profiles
- ModelRouter for category-based provider selection (deep, quick, visual, ultrabrain)
- BackgroundManager for concurrent agent execution (max 5 slots)
- FallbackProvider with automatic failover on 429/500/502/503/504

### IDE-Level Tools
- LSP integration — rename, go-to-definition, find-references
- AST-grep — structural code search and rewrite
- Tmux sessions — persistent terminal sessions
- MCP client — connect to external MCP servers

### UX Parity with Claude Code
- Real-time streaming with token-by-token display
- Thinking block display (Claude extended thinking)
- Model-aware max tokens
- Git context in system prompt (branch, changed files)
- GOCODE.md / CLAUDE.md project config loading
- Cost estimation with `/cost`
- Ctrl+C interrupt support

### 21 Slash Commands
`/help` `/exit` `/clear` `/compact` `/cost` `/model` `/skill` `/plan` `/init-deep` `/diff` `/undo` `/redo` `/status` `/review` `/permissions` `/doctor` `/connect` `/share` `/commit` `/memory` `/tasks`

### CLI Commands
- `gocode serve` — headless HTTP REST API server
- `gocode stats` — usage statistics across sessions
- `gocode export/import` — session import/export
- `gocode pr` — GitHub PR creation via `gh` CLI
- `gocode github` — GitHub issue listing via `gh` CLI
- `gocode auth generate/list/delete` — remote access key management
- `gocode config` — show runtime configuration
- `gocode plugin list/install/uninstall` — plugin management

### Plugin System
- Hook pipeline with pre/post tool-use interception
- 2 bundled plugins: safety-guard, git-auto-commit
- Plugin install/uninstall via CLI

### Editor Compatibility
- Editor detection (VS Code, Cursor, Kiro, Neovim, etc.)
- Editor-specific configuration hints

### Skills System
- 8 built-in skills with community attributions
- `--skill` flag and `/skill` command for mid-session switching
- Custom skills via `.gocode/skills/` JSON files

### Auto-Format
- gofmt/goimports for Go
- prettier for JS/TS
- black for Python
- rustfmt for Rust

---

## v0.3.0 — The Foundation

### Core Architecture
- Complete Go reimplementation of Claude Code agent runtime
- 38 internal packages: agent, apiclient, apitypes, bootstrap, commands, context, execution, hashline, history, manifest, mcp, models, modes, permissions, queryengine, repl, runtime, session, setup, tools, toolimpl, toolpool, transcript
- Single binary (~12MB), zero runtime dependencies
- <10ms startup time

### Multi-Model Support
- 4 native providers: Anthropic (Claude), OpenAI (GPT), Google (Gemini), xAI (Grok)
- Model aliases: `sonnet`, `opus`, `haiku`, `gpt5`, `gpt4o`, `gemini`, `grok`
- `--model` flag for provider selection

### Agent Mode
- Interactive REPL (`gocode chat`)
- One-shot mode (`gocode prompt`)
- Multi-turn tool-use loops
- Permission system (workspace-write, full-access)

### MCP Server
- Full MCP protocol compliance (initialize, tools/list, tools/call, ping, resources/list)
- 14 built-in tools: BashTool, FileReadTool, FileEditTool, FileWriteTool, GlobTool, GrepTool, ListDirectoryTool
- Dual transport: stdio (for IDEs) and HTTP
- Works with Cursor, Kiro, VS Code, Antigravity, Claude Desktop

### Session Management
- Session persistence and resume
- Transcript flushing
- Session store with JSON serialization

### CLI Commands
- `gocode chat` — interactive agent
- `gocode prompt` — one-shot agent
- `gocode mcp-serve` — MCP server
- `gocode summary` — workspace summary
- `gocode manifest` — port manifest
- `gocode subsystems` — module listing
- `gocode commands` — command registry
- `gocode tools` — tool registry
- `gocode route` — prompt routing
- `gocode bootstrap` — session bootstrapping
- `gocode turn-loop` — turn loop execution
- `gocode setup-report` — environment setup report

### Installation
- `go install` for all platforms
- One-line install scripts (macOS/Linux bash, Windows PowerShell)
- Binary downloads for all platforms (darwin/linux/windows, amd64/arm64)
- deb/rpm packages for Linux
- Build from source
