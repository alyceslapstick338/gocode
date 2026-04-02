# UX Features Guide

[← Back to README](../README.md)

gocode v0.4.0 brings Claude Code-level UX to every model. Streaming responses, thinking blocks, git awareness, project config, and smart token limits — all working out of the box.

---

## Streaming Responses

Chat mode now streams tokens as they arrive. You see the response forming in real-time instead of waiting for the full response.

The spinner shows while waiting for the first token. Once text starts flowing, it switches to live streaming. During tool execution, the spinner pauses and shows tool progress, then resumes for the next LLM call.

---

## Slash Commands

| Command | What It Does |
|---------|-------------|
| `/help` | Show all available commands |
| `/exit` | Quit the session (also Ctrl+D) |
| `/clear` | Reset conversation history |
| `/compact` | Compact conversation — keeps system prompt + last 5 message pairs |
| `/cost` | Show token usage with estimated dollar cost |
| `/model` | Show current model (or suggest restart with a different one) |
| `/skill` | List or activate skills mid-session |
| `/plan` | Start an interview-style planning session |
| `/init-deep` | Generate AGENTS.md context files |
| `/diff` | Show git diff of changes made this session |

### /compact

When your conversation gets long and the context window fills up, `/compact` trims it down while keeping the most recent context:

```
you> /compact
Compacted: 47 → 10 messages
```

### /diff

See what the agent changed during this session:

```
you> /diff
diff --git a/main.go b/main.go
index abc1234..def5678 100644
--- a/main.go
+++ b/main.go
@@ -10,6 +10,7 @@ func main() {
+    fmt.Println("new line")
```

---

## Project Config (GOCODE.md / CLAUDE.md)

Drop a `GOCODE.md` or `CLAUDE.md` file in your project root with project-specific instructions. gocode reads it on startup and injects it into the system prompt.

```markdown
# Project Instructions

This is a Go project using Chi router and PostgreSQL.
- Always use structured logging with slog
- Run `make test` after changes
- Database migrations are in db/migrations/
```

gocode checks for `GOCODE.md` first, then falls back to `CLAUDE.md`. This means your existing Claude Code project configs work automatically.

---

## Thinking Block Display

When using models that support extended thinking (Claude Opus, Claude Sonnet with thinking enabled), gocode displays the model's internal reasoning in a dimmed 💭 section:

```
💭 Let me analyze the codebase structure first. The user wants to refactor
the auth module, so I should check how authentication is currently implemented
across the different packages...

assistant> I've analyzed your auth module. Here's what I found...
```

Thinking blocks are shown in gray/dim text so they're visible but don't compete with the actual response.

---

## Model-Aware Token Limits

gocode automatically sets the right `max_tokens` for each model. No more sending 8192 to a model that supports 128K output.

| Model | Max Output Tokens |
|-------|------------------|
| GPT-5.4 | 128,000 |
| Grok 4.20 | 131,072 |
| o3 / o4 | 100,000 |
| Claude Sonnet 4.6 | 64,000 |
| Gemini 3.1 Pro | 65,536 |
| Claude Opus 4.6 | 32,000 |
| GPT-4o | 16,384 |
| Claude Haiku 4.5 | 8,192 |

You can still override with `--max-tokens` if needed.

---

## Git Context

When you're in a git repository, gocode automatically includes the current branch and number of changed files in the system prompt. The agent knows your git state without you telling it.

```
# Git
- Branch: feat/new-feature
- Changed files: 3
```

---

## Cost Estimation

`/cost` now shows estimated dollar cost alongside token counts:

```
you> /cost
Tokens: 12,450 in / 3,200 out (15,650 total, 5 turns) — est. $0.0854
```

Pricing is a blended average ($3/1M input, $15/1M output). Actual cost varies by provider and model.

---

[← Back to README](../README.md)
