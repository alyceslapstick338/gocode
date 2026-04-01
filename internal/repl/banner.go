package repl

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ANSI color codes — blue gradient matching Claude Code's style
const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	dim     = "\033[2m"
	b0      = "\033[38;5;17m"  // darkest blue (background chars)
	b1      = "\033[38;5;18m"  // dark blue
	b2      = "\033[38;5;25m"  // medium-dark blue
	b3      = "\033[38;5;33m"  // medium blue
	b4      = "\033[38;5;39m"  // bright blue
	b5      = "\033[38;5;45m"  // light blue
	b6      = "\033[38;5;51m"  // cyan-blue (highlights)
	bw      = "\033[38;5;255m" // white (eyes)
	white   = "\033[38;5;255m"
	gray    = "\033[38;5;242m"
	green   = "\033[38;5;114m"
	yellow  = "\033[38;5;221m"
	boxChar = "\033[38;5;60m"
)

// Scaled-down Go gopher in blue gradient, ~12 lines tall
// Based on the classic gopher silhouette, colored like Claude Code's claw
var gopher = []string{
	b1 + `      .=+#%` + b3 + `@@@@@@` + b1 + `%#+=.` + reset,
	b2 + `   .+%%*=` + b3 + `----------` + b2 + `=*%%+.` + reset,
	b2 + ` .*@%` + b4 + `@+--+#@@#=` + b3 + `----` + b4 + `*@#+--=@%` + b2 + `*.` + reset,
	b3 + ` *@+--=%` + b4 + `#=.` + reset + `    ` + b4 + `.+%#` + b3 + `--*%=` + reset + `    ` + b4 + `*%--%` + b3 + `@*` + reset,
	b3 + ` %=-=%` + b4 + `@%+-%*` + reset + `       ` + b4 + `=%--#%` + reset + `     ` + b4 + `#%--#%` + reset,
	b3 + ` %+--=@` + b4 + `@+--%-` + bw + `@@` + b4 + `+` + reset + `    ` + b4 + `:%--#+` + bw + `.@@` + b4 + `+` + reset + `  ` + b4 + `#%--+@` + reset,
	b3 + `  =@*=` + b4 + `=@*--%-*@` + b5 + `+=` + reset + `    ` + b4 + `:%--##` + b5 + `.%@+` + reset + `  ` + b4 + `.#%--%%` + reset,
	b3 + `    =%` + b4 + `@%----%#` + reset + `       ` + b4 + `+#--=@*` + reset + `     ` + b4 + `+@---%%` + reset,
	b4 + `     =%---` + b5 + `---+@#-..` + b4 + `=%@+=%` + b5 + `@@*#@@` + b4 + `#*+*%@*` + reset,
	b4 + `     *#---` + b5 + `--------*%-:::--:::*%=---` + b4 + `---*#` + reset,
	b4 + `     **---` + b5 + `--------*%::::+**=::-%+--` + b4 + `---+%` + reset,
	b4 + `     *#---` + b5 + `---------+%%#+=@=-#%%+---` + b4 + `---=%` + reset,
	b4 + `     +#=--` + b5 + `----------=%- -%- +%=----` + b4 + `---=%` + reset,
	b3 + `      =%=--` + b4 + `----------------------------` + b3 + `=%` + reset,
	b3 + `      -%=--` + b4 + `----------------------------` + b3 + `=%` + reset,
	b2 + `      :#*--` + b3 + `----------------------------` + b2 + `=%` + reset,
	b2 + `       +%--` + b3 + `----------------------------` + b2 + `+%` + reset,
	b1 + `        :%*` + b2 + `-------------------------` + b1 + `=#%` + reset,
	b1 + `         -%` + b2 + `@%+-------------------=#%` + b1 + `*@` + reset,
	b0 + `          :@` + b1 + `#::*@%+----------=#@%-::+@` + reset,
	b0 + `           =@` + b1 + `%+:-%+.  .-+*#%%#*+-.  *%-%@` + reset,
	b0 + `             :#@@%..             ..=*=` + reset,
}

// BannerConfig holds the info displayed in the welcome banner.
type BannerConfig struct {
	Version  string
	Model    string
	MaxTurns int
	Cwd      string
}

// PrintBanner renders the welcome screen with a blue Go gopher.
func PrintBanner(w io.Writer, cfg BannerConfig) {
	cwd := cfg.Cwd
	if cwd == "" {
		cwd, _ = os.Getwd()
	}
	if home, err := os.UserHomeDir(); err == nil {
		if rel, err := filepath.Rel(home, cwd); err == nil && !strings.HasPrefix(rel, "..") {
			cwd = "~/" + rel
		}
	}

	version := cfg.Version
	if version == "" {
		version = "dev"
	}

	width := 72
	title := fmt.Sprintf(" gocode %s ", version)
	padTotal := width - 2 - len(title)
	padLeft := padTotal / 2
	padRight := padTotal - padLeft

	// Top border
	fmt.Fprintf(w, "%s╭%s%s%s%s%s%s%s╮%s\n",
		boxChar, strings.Repeat("─", padLeft), reset, bold+b6, title, reset, boxChar, strings.Repeat("─", padRight), reset)

	// Welcome line
	welcome := fmt.Sprintf("  %sWelcome to gocode%s", bold+white, reset)
	printPaddedLine(w, welcome, width)
	printEmptyLine(w, width)

	// Gopher art (centered)
	for _, line := range gopher {
		printPaddedLine(w, "  "+line, width)
	}

	printEmptyLine(w, width)

	// Model + turns info
	modelLine := fmt.Sprintf("  %s%s%s • Max %s%d%s turns", bold+b5, cfg.Model, reset, bold+b5, cfg.MaxTurns, reset)
	printPaddedLine(w, modelLine, width)

	// Working directory
	cwdLine := fmt.Sprintf("  %s%s%s", dim+gray, cwd, reset)
	printPaddedLine(w, cwdLine, width)

	printEmptyLine(w, width)

	// Tips section
	tipsHeader := fmt.Sprintf("  %s%sTips%s", bold, green, reset)
	printPaddedLine(w, tipsHeader, width)
	tips := []struct{ cmd, desc string }{
		{"/exit", "quit session"},
		{"/clear", "reset history"},
		{"/cost", "token usage"},
		{"Ctrl+D", "quit (EOF)"},
		{`line\`, "multi-line input"},
	}
	for _, t := range tips {
		line := fmt.Sprintf("    %s%-12s%s %s%s%s", gray, t.cmd, reset, white, t.desc, reset)
		printPaddedLine(w, line, width)
	}

	printEmptyLine(w, width)

	// Models section
	modelsHeader := fmt.Sprintf("  %s%sModels%s", bold, yellow, reset)
	printPaddedLine(w, modelsHeader, width)
	models := []struct{ alias, provider string }{
		{"sonnet", "Claude"},
		{"gpt4o", "GPT-4o"},
		{"gemini", "Gemini"},
		{"grok", "Grok 3"},
		{"o3", "OpenAI o3"},
		{"codex", "Codex"},
	}
	for _, m := range models {
		line := fmt.Sprintf("    %s--model %s%-8s %s%s%s", gray, white, m.alias, gray, m.provider, reset)
		printPaddedLine(w, line, width)
	}

	printEmptyLine(w, width)

	// Bottom border
	fmt.Fprintf(w, "%s╰%s╯%s\n\n", boxChar, strings.Repeat("─", width-2), reset)
}

func printPaddedLine(w io.Writer, content string, width int) {
	vis := visibleLen(content)
	pad := width - 2 - vis
	if pad < 0 {
		pad = 0
	}
	fmt.Fprintf(w, "%s│%s%s%s %s│%s\n", boxChar, reset, content, strings.Repeat(" ", pad), boxChar, reset)
}

func printEmptyLine(w io.Writer, width int) {
	fmt.Fprintf(w, "%s│%s%s%s│%s\n", boxChar, reset, strings.Repeat(" ", width-2), boxChar, reset)
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
