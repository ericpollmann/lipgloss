package lipgloss

import (
	"os"
	"strings"
	"sync"

	"github.com/charmbracelet/x/ansi"
)

var (
	needsEmojiFixCache     bool
	needsEmojiFixCacheOnce sync.Once
)

// needsEmojiWidthFix returns true for terminals that render emoji with variation
// selectors as 1 column but ansi.StringWidth reports 2. This affects iTerm2,
// Terminal.app, Alacritty, and VSCode. Kitty correctly adjusts width with
// variation selectors, so we skip the fix there.
//
// Note: This uses environment variable detection which is fragile but necessary
// until terminals standardize their emoji+VS rendering behavior.
func needsEmojiWidthFix() bool {
	needsEmojiFixCacheOnce.Do(func() {
		termProgram := os.Getenv("TERM_PROGRAM")
		lcTerminal := os.Getenv("LC_TERMINAL")
		term := os.Getenv("TERM")

		// Kitty handles emoji+VS width correctly, don't apply fix
		if strings.Contains(term, "kitty") || os.Getenv("KITTY_WINDOW_ID") != "" {
			needsEmojiFixCache = false
			return
		}

		// Default to true for common terminals that need the fix
		needsEmojiFixCache = termProgram == "iTerm.app" ||
			termProgram == "Apple_Terminal" ||
			termProgram == "vscode" ||
			lcTerminal == "iTerm2" ||
			os.Getenv("ALACRITTY_SOCKET") != ""
	})

	return needsEmojiFixCache
}

// tuiWidth returns the visual width of a string in TUI mode.
// For terminals that misrender emoji+VS sequences, this adjusts for variation
// selectors (U+FE00-U+FE0F) which render as 0 width but ansi.StringWidth
// reports as 1 width.
func tuiWidth(s string) int {
	width := ansi.StringWidth(s)
	if !needsEmojiWidthFix() {
		return width
	}

	// Subtract width for each variation selector
	for _, r := range []rune(s) {
		if r >= 0xFE00 && r <= 0xFE0F {
			width--
		}
	}
	return width
}

// normalizeEmojiWidth adds space after emoji variation selectors to compensate
// for terminals that render emoji+VS as 1 column. This ensures the visual width
// matches the calculated width (emoji+VS=1 + space=1 = 2 total).
func normalizeEmojiWidth(s string) string {
	if !needsEmojiWidthFix() {
		return s
	}

	runes := []rune(s)
	var result []rune

	for i := 0; i < len(runes); i++ {
		result = append(result, runes[i])

		// Add space after variation selector if preceded by likely emoji
		if runes[i] >= 0xFE00 && runes[i] <= 0xFE0F && i > 0 {
			prev := runes[i-1]
			if prev >= 0x2000 { // Most emoji are above U+2000
				result = append(result, ' ')
			}
		}
	}
	return string(result)
}
