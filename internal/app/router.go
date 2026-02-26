package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Router struct {
	current ScreenType
	screens map[ScreenType]Screen

	width  int
	height int
}

func NewRouter(models []Screen) Router {
	screens := make(map[ScreenType]Screen)
	for _, m := range models {
		screens[m.ID()] = m
	}

	return Router{
		current: models[0].ID(),
		screens: screens,
	}
}

func (r Router) Init() tea.Cmd {
	return tea.ClearScreen
}

func (r Router) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// TEMPORARY EXIT
		switch msg.String() {
		case "t":
			return r, tea.Quit
		}
	case tea.WindowSizeMsg:
		r.width = msg.Width
		r.height = msg.Height
		msg.Height -= lipgloss.Height(r.header())

		// TODO: consider a broadcast
		var cmd tea.Cmd
		r.screens[r.current], cmd = r.screens[r.current].Update(msg)
		return r, cmd
	}

	var cmd tea.Cmd
	r.screens[r.current], cmd = r.screens[r.current].Update(msg)
	return r, cmd
}

func (r Router) header() string {
	return lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(r.width).
		Border(lipgloss.InnerHalfBlockBorder(), false, false, true, false).
		Render("Friendly CLI")
}

func (r Router) View() string {
	header := r.header()

	content := lipgloss.NewStyle().
		Width(r.width).
		Height(r.height - lipgloss.Height(header)).
		Render(r.screens[r.current].View())

	return lipgloss.JoinVertical(lipgloss.Top, header, content)
}
