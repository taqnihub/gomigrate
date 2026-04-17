package tui

import "github.com/charmbracelet/lipgloss"

// Color palette — consistent across the whole app.
// Using adaptive colors so terminals with light/dark backgrounds both look good.
var (
	ColorSuccess   = lipgloss.AdaptiveColor{Light: "#059669", Dark: "#10B981"}
	ColorError     = lipgloss.AdaptiveColor{Light: "#DC2626", Dark: "#EF4444"}
	ColorWarning   = lipgloss.AdaptiveColor{Light: "#D97706", Dark: "#F59E0B"}
	ColorInfo      = lipgloss.AdaptiveColor{Light: "#2563EB", Dark: "#3B82F6"}
	ColorBrand     = lipgloss.AdaptiveColor{Light: "#7C3AED", Dark: "#8B5CF6"}
	ColorMuted     = lipgloss.AdaptiveColor{Light: "#6B7280", Dark: "#9CA3AF"}
	ColorSubtle    = lipgloss.AdaptiveColor{Light: "#9CA3AF", Dark: "#6B7280"}
	ColorText      = lipgloss.AdaptiveColor{Light: "#111827", Dark: "#F9FAFB"}
	ColorHighlight = lipgloss.AdaptiveColor{Light: "#F3F4F6", Dark: "#1F2937"}
)

// Reusable styles — define once, use everywhere.
var (
	// Status indicators
	SuccessStyle = lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true)
	ErrorStyle   = lipgloss.NewStyle().Foreground(ColorError).Bold(true)
	WarningStyle = lipgloss.NewStyle().Foreground(ColorWarning).Bold(true)
	InfoStyle    = lipgloss.NewStyle().Foreground(ColorInfo).Bold(true)

	// Text styles
	TitleStyle = lipgloss.NewStyle().
			Foreground(ColorBrand).
			Bold(true).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorText).
			Bold(true)

	MutedStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	SubtleStyle = lipgloss.NewStyle().
			Foreground(ColorSubtle)

	HighlightStyle = lipgloss.NewStyle().
			Foreground(ColorBrand).
			Bold(true)

	// Box styles for headers/callouts
	HeaderBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBrand).
			Padding(0, 2).
			MarginBottom(1)

	WarningBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorWarning).
			Padding(0, 2).
			MarginBottom(1)

	ErrorBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorError).
			Padding(0, 2).
			MarginBottom(1)

	// Indent for nested content
	IndentStyle = lipgloss.NewStyle().MarginLeft(2)

	// Code/path formatting
	CodeStyle = lipgloss.NewStyle().
			Foreground(ColorBrand).
			Background(ColorHighlight).
			Padding(0, 1)

	PathStyle = lipgloss.NewStyle().
			Foreground(ColorInfo).
			Underline(true)
)

// Icons — Unicode glyphs used throughout the UI.
// Using these constants keeps icons consistent across commands.
const (
	IconSuccess = "✓"
	IconError   = "✗"
	IconWarning = "⚠"
	IconInfo    = "ℹ"
	IconArrow   = "→"
	IconBullet  = "•"
	IconStar    = "★"
	IconRocket  = "🚀"
	IconSpinner = "⣾"
	IconCheck   = "✔"
	IconCross   = "✘"
)
