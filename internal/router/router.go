package router

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/friendly-social/cli/internal/screen"
)

type Router struct {
	current screen.Type
	screens map[screen.Type]screen.Model

	width  int
	height int
}

func NewRouter(models []screen.Model) Router {
	screens := make(map[screen.Type]screen.Model)
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

func (r Router) broadcast(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, len(r.screens))
	for i := range r.screens {
		r.screens[i], cmds[i] = r.screens[i].Update(msg)
	}

	return r, tea.Batch(cmds...)
}

func (r Router) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		r.width = msg.Width
		r.height = msg.Height

		msg.Height -= lipgloss.Height(r.header())
		return r.broadcast(msg)
	case BroadcastMsg:
		return r.broadcast(msg.Inner)
	case screen.ChangeMsg:
		r.current = msg.NewType
		return r, nil
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
