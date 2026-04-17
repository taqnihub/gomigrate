package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// ConfirmModel is a yes/no dialog.
type ConfirmModel struct {
	question  string
	danger    bool // if true, "No" is the default
	answer    bool
	answered  bool
	cancelled bool
}

// NewConfirmModel creates a confirmation prompt.
// Set danger=true for destructive actions (defaults to "No").
func NewConfirmModel(question string, danger bool) ConfirmModel {
	return ConfirmModel{
		question: question,
		danger:   danger,
		answer:   !danger, // default: Yes for safe, No for dangerous
	}
}

func (m ConfirmModel) Init() tea.Cmd {
	return nil
}

func (m ConfirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "ctrl+c", "esc":
			m.cancelled = true
			return m, tea.Quit

		case "left", "h":
			m.answer = true // Yes

		case "right", "l":
			m.answer = false // No

		case "y", "Y":
			m.answer = true
			m.answered = true
			return m, tea.Quit

		case "n", "N":
			m.answer = false
			m.answered = true
			return m, tea.Quit

		case "enter", " ":
			m.answered = true
			return m, tea.Quit

		case "tab":
			m.answer = !m.answer
		}
	}

	return m, nil
}

func (m ConfirmModel) View() string {
	if m.answered || m.cancelled {
		return ""
	}

	prompt := fmt.Sprintf("  %s\n\n", SubtitleStyle.Render(m.question))

	var yes, no string
	if m.answer {
		yes = HighlightStyle.Render(" [ Yes ] ")
		no = MutedStyle.Render("   No   ")
	} else {
		yes = MutedStyle.Render("   Yes  ")
		no = HighlightStyle.Render(" [ No ] ")
	}

	prompt += fmt.Sprintf("  %s  %s\n\n", yes, no)
	prompt += MutedStyle.Render("  ←/→ or y/n to select · enter to confirm · esc to cancel\n")

	return prompt
}

// Confirmed returns true if the user selected Yes and confirmed.
func (m ConfirmModel) Confirmed() bool {
	return m.answered && m.answer
}

// Cancelled returns true if the user pressed escape.
func (m ConfirmModel) Cancelled() bool {
	return m.cancelled
}
