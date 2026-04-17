package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Success prints a success message with a green checkmark.
func Success(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	icon := SuccessStyle.Render(IconSuccess)
	fmt.Printf("  %s %s\n", icon, msg)
}

// Error prints an error message with a red X.
func Error(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	icon := ErrorStyle.Render(IconError)
	fmt.Printf("  %s %s\n", icon, msg)
}

// Warning prints a warning with a yellow exclamation.
func Warning(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	icon := WarningStyle.Render(IconWarning)
	fmt.Printf("  %s %s\n", icon, msg)
}

// Info prints an info message with a blue icon.
func Info(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	icon := InfoStyle.Render(IconInfo)
	fmt.Printf("  %s %s\n", icon, msg)
}

// Muted prints text in a subtle gray color (for secondary info).
func Muted(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println(MutedStyle.Render("  " + msg))
}

// Title prints a bold, colored title with spacing around it.
func Title(text string) {
	fmt.Println()
	fmt.Println(TitleStyle.Render("  " + text))
}

// Header prints a boxed header (used at the start of commands).
func Header(title, subtitle string) {
	content := TitleStyle.Render(title)
	if subtitle != "" {
		content += "\n" + MutedStyle.Render(subtitle)
	}
	fmt.Println(HeaderBoxStyle.Render(content))
}

// Section prints a section separator.
func Section(text string) {
	fmt.Println()
	fmt.Println(SubtitleStyle.Render("  " + text))
	fmt.Println(SubtleStyle.Render("  " + strings.Repeat("─", len(text))))
}

// KeyValue prints a "key: value" pair with aligned formatting.
// Keys are muted, values are highlighted.
func KeyValue(key, value string) {
	fmt.Printf("  %s %s\n",
		MutedStyle.Render(fmt.Sprintf("%-16s", key+":")),
		value,
	)
}

// Code formats a string as inline code (highlighted).
func Code(text string) string {
	return CodeStyle.Render(text)
}

// Path formats a filesystem path with underline.
func Path(p string) string {
	return PathStyle.Render(p)
}

// Bold makes text bold.
func Bold(text string) string {
	return lipgloss.NewStyle().Bold(true).Render(text)
}

// Dim makes text faded.
func Dim(text string) string {
	return MutedStyle.Render(text)
}

// Box wraps content in a bordered box.
func Box(content string) string {
	return HeaderBoxStyle.Render(content)
}

// WarningBox wraps content in a yellow-bordered warning box.
func WarningBox(content string) string {
	return WarningBoxStyle.Render(content)
}

// ErrorBox wraps content in a red-bordered error box.
func ErrorBox(content string) string {
	return ErrorBoxStyle.Render(content)
}

// Banner prints the gomigrate logo/banner at startup.
func Banner() {
	logo := `
   ┌─────────────────────────┐
   │   ⚡ GoMigrate  v0.1     │
   │   Database migrations   │
   └─────────────────────────┘
`
	fmt.Println(TitleStyle.Render(logo))
}

// Newline prints a blank line (for readability).
func Newline() {
	fmt.Println()
}
