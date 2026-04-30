package home

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/friendly-social/cli/internal/router"
	"github.com/friendly-social/cli/internal/screen"
	"github.com/friendly-social/cli/internal/ui"
)

// Screen is a model of home screen.
type Screen struct {
	content struct {
		list *ui.List

		buttons struct {
			exit *ui.Button
			back *ui.Button
		}
	}
}

// New returns new initial model of home screen.
func New() Screen {
	result := Screen{}

	result.content.buttons.exit = ui.NewButton("Exit", tea.Quit)
	result.content.buttons.back = ui.NewButton("Back", func() tea.Msg {
		return screen.ChangeMsg{NewType: screen.TypeAuth}
	})

	result.content.list = ui.NewList(
		result.content.buttons.back,
		result.content.buttons.exit)

	return result
}

func (s Screen) ID() screen.Type {
	return screen.TypeHome
}

func (s Screen) Init() tea.Cmd {
	return func() tea.Msg {
		return router.TargetMsg{Type: s.ID(), Inner: ui.SelectMsg{}}
	}
}

func (s Screen) Update(msg tea.Msg) (screen.Model, tea.Cmd) {
	_, cmd := s.content.list.Update(msg)
	return s, cmd
}

func (s Screen) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		"home screen placeholder...",
		"",
		s.content.list.View(),
	)
}
