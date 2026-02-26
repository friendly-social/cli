package auth

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/friendly-social/cli/internal/app"
	"github.com/friendly-social/cli/internal/navigation"
)

type Screen struct {
	width  int
	height int

	temp tea.Msg
}

func NewScreen() Screen {
	return Screen{}
}

func (s Screen) ID() app.ScreenType {
	return app.ScreenAuth
}

func (s Screen) Init() tea.Cmd {
	return nil
}

func (s Screen) Update(msg tea.Msg) (app.Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
	}

	s.temp = msg
	return s, nil
}

func (s Screen) View() string {
	switch s.temp.(type) {
	case navigation.MovedMsg:
		return "moved"
	case navigation.FocusedMsg:
		return "focused"
	case navigation.UnfocusedMsg:
		return "unfocused"
	case navigation.InteractedMsg:
		return "interacted"
	default:
		return "other"
	}
}
