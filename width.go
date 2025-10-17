package lipgloss

import "github.com/charmbracelet/x/ansi"

// tuiWidth returns the visual width of a string in TUI mode, adjusting for
// variation selectors which render as 0 width but ansi.StringWidth reports as 1.
func tuiWidth(s string) int {
	width := ansi.StringWidth(s)
	for _, r := range []rune(s) {
		if r >= 0xFE00 && r <= 0xFE0F {
			width--
		}
	}
	return width
}

// normalizeEmojiWidth adds space after emoji variation selectors to match calculated width.
func normalizeEmojiWidth(s string) string {
	runes := []rune(s)
	var result []rune

	for i := 0; i < len(runes); i++ {
		result = append(result, runes[i])
		if runes[i] >= 0xFE00 && runes[i] <= 0xFE0F && i > 0 {
			prev := runes[i-1]
			if prev >= 0x2000 {
				result = append(result, ' ')
			}
		}
	}
	return string(result)
}
