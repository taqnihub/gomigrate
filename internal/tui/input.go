package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// InputModel is a simple text input prompt.
type InputModel struct {
	prompt    string
	input     textinput.Model
	submitted bool
	cancelled bool
}

// NewInputModel creates a text input with the given prompt.
func NewInputModel(prompt, placeholder string) InputModel {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 40

	return InputModel{
		prompt: prompt,
		input:  ti,
	}
}

func (m InputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m InputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "ctrl+c", "esc":
			m.cancelled = true
			return m, tea.Quit
		case "enter":
			if m.input.Value() != "" {
				m.submitted = true
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m InputModel) View() string {
	if m.submitted || m.cancelled {
		return ""
	}

	s := fmt.Sprintf("  %s\n\n", SubtitleStyle.Render(m.prompt))
	s += "  " + m.input.View() + "\n\n"
	s += MutedStyle.Render("  enter to submit · esc to cancel\n")
	return s
}

// Value returns the entered text.
func (m InputModel) Value() string {
	return m.input.Value()
}

// Submitted returns true if the user pressed enter.
func (m InputModel) Submitted() bool {
	return m.submitted
}

// Cancelled returns true if the user pressed escape.
func (m InputModel) Cancelled() bool {
	return m.cancelled
}
