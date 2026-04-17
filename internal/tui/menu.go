package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/taqnihub/gomigrate/config"
)

// MenuItem represents one option in the main menu.
type MenuItem struct {
	Icon  string
	Label string
	Desc  string
	Key   string // action identifier returned when selected
}

// MenuModel holds the state of the interactive menu.
type MenuModel struct {
	choices  []MenuItem
	cursor   int
	selected string // set to the Key of the chosen item when user confirms
	quitting bool
	config   *config.Config
	width    int
	height   int
}

// NewMenuModel creates a new menu with default choices.
func NewMenuModel(cfg *config.Config) MenuModel {
	return MenuModel{
		config: cfg,
		choices: []MenuItem{
			{Icon: "▶", Label: "Apply migrations", Desc: "Run all pending migrations", Key: "up"},
			{Icon: "◀", Label: "Revert last migration", Desc: "Undo the most recent migration", Key: "down"},
			{Icon: "✚", Label: "Create new migration", Desc: "Generate up/down SQL files", Key: "create"},
			{Icon: "📋", Label: "View status", Desc: "Show all migrations and their state", Key: "status"},
			{Icon: "🔧", Label: "Force version", Desc: "Fix a dirty migration state", Key: "force"},
			{Icon: "q", Label: "Quit", Desc: "Exit gomigrate", Key: "quit"},
		},
	}
}

// Init is called once when the program starts. We don't need any initial
// commands, so we return nil.
func (m MenuModel) Init() tea.Cmd {
	return nil
}

// Update handles events (keypresses, window resizes) and returns an
// updated model and any commands to run.
func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		// Remember the terminal size for layout
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			} else {
				// wrap to bottom
				m.cursor = len(m.choices) - 1
			}

		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			} else {
				// wrap to top
				m.cursor = 0
			}

		case "enter", " ":
			choice := m.choices[m.cursor]
			m.selected = choice.Key
			if choice.Key == "quit" {
				m.quitting = true
			}
			return m, tea.Quit
		}
	}

	return m, nil
}

// View renders the menu as a string.
func (m MenuModel) View() string {
	if m.quitting {
		return ""
	}

	var s string

	// Header
	header := TitleStyle.Render("⚡ GoMigrate")
	s += header + "\n"

	if m.config != nil {
		dbInfo := fmt.Sprintf("%s · %s:%d · %s",
			m.config.Driver, m.config.Host, m.config.Port, m.config.Database)
		s += MutedStyle.Render("  "+dbInfo) + "\n"
	}
	s += "\n"

	// Prompt
	s += SubtitleStyle.Render("  What would you like to do?") + "\n\n"

	// Menu items
	for i, choice := range m.choices {
		cursor := "  "
		if m.cursor == i {
			cursor = HighlightStyle.Render(" >")
		}

		icon := choice.Icon
		label := choice.Label

		if m.cursor == i {
			// Highlighted item: bold + colored
			icon = HighlightStyle.Render(icon)
			label = HighlightStyle.Render(label)
			desc := MutedStyle.Render("  " + choice.Desc)
			s += fmt.Sprintf("%s %s  %s\n%s\n", cursor, icon, label, desc)
		} else {
			// Normal item
			s += fmt.Sprintf("%s %s  %s\n", cursor, icon, label)
		}
	}

	// Footer with keybinds
	s += "\n"
	footer := lipgloss.NewStyle().
		Foreground(ColorSubtle).
		Render("  ↑/↓ navigate · enter select · q quit")
	s += footer + "\n"

	return s
}

// Selected returns the key of the chosen menu item, or empty string if
// the user quit without choosing.
func (m MenuModel) Selected() string {
	if m.quitting && m.selected != "quit" {
		return ""
	}
	return m.selected
}
