package repl

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ANSI color codes
const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	dim     = "\033[2m"
	blue    = "\033[38;5;39m"
	cyan    = "\033[38;5;51m"
	white   = "\033[38;5;255m"
	gray    = "\033[38;5;242m"
	green   = "\033[38;5;114m"
	yellow  = "\033[38;5;221m"
	boxChar = "\033[38;5;60m"
)

// gopher is the Go gopher in blue ASCII art
var gopher = []string{
	blue + `    ╭───╮  ╭───╮` + reset,
	blue + `    │ ` + white + `◉` + blue + ` │  │ ` + white + `◉` + blue + ` │` + reset,
	blue + `    ╰─┬─╯  ╰─┬─╯` + reset,
	blue + `   ╭──┴──────┴──╮` + reset,
	blue + `   │  ` + cyan + `▀▀▀▀▀▀` + blue + `  │` + reset,
	blue + `   │    ` + white + `╰──╯` + blue + `    │` + reset,
	blue + `   ╰──┬──────┬──╯` + reset,
	blue + `      │ ` + cyan + `╭──╮` + blue + ` │` + reset,
	blue + `      │ ` + cyan + `╰──╯` + blue + ` │` + reset,
	blue + `      ╰──┬┬──╯` + reset,
	blue + `         ││` + reset,
	blue + `        ╰╯╰╯` + reset,
}

// BannerConfig holds the info displayed in the welcome banner.
type BannerConfig struct {
	Version string
	Model   string
	MaxTurns int
	Cwd     string
}

// PrintBanner renders the Claude-Code-style welcome screen with a Go gopher.
func PrintBanner(w io.Writer, cfg BannerConfig) {
	cwd := cfg.Cwd
	if cwd == "" {
		cwd, _ = os.Getwd()
	}
	// Shorten home dir
	if home, err := os.UserHomeDir(); err == nil {
		if rel, err := filepath.Rel(home, cwd); err == nil && !strings.HasPrefix(rel, "..") {
			cwd = "~/" + rel
		}
	}

	width := 64
	version := cfg.Version
	if version == "" {
		version = "dev"
	}

	// Top border with version
	title := fmt.Sprintf(" gocode %s ", version)
	padTotal := width - 2 - len(title)
	padLeft := padTotal / 2
	padRight := padTotal - padLeft
	topBorder := boxChar + "╭" + strings.Repeat("─", padLeft) + reset + bold + cyan + title + reset + boxChar + strings.Repeat("─", padRight) + "╮" + reset

	fmt.Fprintln(w, topBorder)
	fmt.Fprintln(w, boxChar+"│"+reset+strings.Repeat(" ", width-2)+boxChar+"│"+reset)

	// Left side: gopher + model info
	// Right side: tips
	leftLines := make([]string, 0, 16)
	leftLines = append(leftLines, "")
	leftLines = append(leftLines, gopher...)
	leftLines = append(leftLines, "")
	leftLines = append(leftLines, fmt.Sprintf("  %s%s%s • Max %s%d%sx", bold+white, cfg.Model, reset, bold+white, cfg.MaxTurns, reset))
	leftLines = append(leftLines, fmt.Sprintf("  %s%s%s", dim+gray, cwd, reset))

	rightLines := []string{
		green + bold + "Tips" + reset,
		gray + "  /exit       " + white + "quit session" + reset,
		gray + "  /clear      " + white + "reset history" + reset,
		gray + "  /cost       " + white + "token usage" + reset,
		gray + "  Ctrl+D      " + white + "quit (EOF)" + reset,
		gray + "  line\\       " + white + "multi-line input" + reset,
		"",
		yellow + bold + "Supported models" + reset,
		gray + "  --model " + white + "sonnet" + gray + "   Claude" + reset,
		gray + "  --model " + white + "gpt4o" + gray + "    GPT-4o" + reset,
		gray + "  --model " + white + "gemini" + gray + "   Gemini" + reset,
		gray + "  --model " + white + "grok" + gray + "     Grok 3" + reset,
		gray + "  --model " + white + "o3" + gray + "       OpenAI o3" + reset,
		gray + "  --model " + white + "codex" + gray + "    Codex" + reset,
	}

	// Pad to same length
	maxLines := len(leftLines)
	if len(rightLines) > maxLines {
		maxLines = len(rightLines)
	}
	for len(leftLines) < maxLines {
		leftLines = append(leftLines, "")
	}
	for len(rightLines) < maxLines {
		rightLines = append(rightLines, "")
	}

	leftColWidth := 30
	for i := 0; i < maxLines; i++ {
		left := leftLines[i]
		right := rightLines[i]
		leftVisible := visibleLen(left)
		padding := leftColWidth - leftVisible
		if padding < 0 {
			padding = 0
		}
		rightPad := width - 2 - leftColWidth - visibleLen(right)
		if rightPad < 0 {
			rightPad = 0
		}
		fmt.Fprintf(w, "%s│%s %s%s%s%s %s│%s\n",
			boxChar, reset,
			left, strings.Repeat(" ", padding),
			right, strings.Repeat(" ", rightPad),
			boxChar, reset)
	}

	fmt.Fprintln(w, boxChar+"│"+reset+strings.Repeat(" ", width-2)+boxChar+"│"+reset)

	// Bottom border
	fmt.Fprintln(w, boxChar+"╰"+strings.Repeat("─", width-2)+"╯"+reset)
	fmt.Fprintln(w)
}

// visibleLen returns the length of a string excluding ANSI escape codes.
func visibleLen(s string) int {
	inEscape := false
	count := 0
	for _, r := range s {
		if r == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEscape = false
			}
			continue
		}
		count++
	}
	return count
}
