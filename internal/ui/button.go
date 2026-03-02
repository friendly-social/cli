package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/friendly-social/cli/internal/navigation"
)

type Button struct {
	cmd      tea.Cmd
	title    string
	selected bool

	selectedStyle   lipgloss.Style
	unselectedStyle lipgloss.Style
}

func NewButton(title string, cmd tea.Cmd) Button {
	return Button{
		cmd:      cmd,
		title:    title,
		selected: false,

		selectedStyle:   lipgloss.NewStyle().Background(lipgloss.Color("#7f00ff")),
		unselectedStyle: lipgloss.NewStyle().Background(lipgloss.Color("#808080")),
	}
}

func (b Button) Init() tea.Cmd {
	return nil
}

func (b Button) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case navigation.SelectMsg:
		b.selected = true
	case navigation.UnselectMsg:
		b.selected = false
	case navigation.InteractMsg:
		return b, b.cmd
	}

	return b, nil
}

func (b Button) View() string {
	if b.selected {
		return b.selectedStyle.Render(b.title)
	}

	return b.unselectedStyle.Render(b.title)
}
