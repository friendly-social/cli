package home

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/friendly-social/cli/internal/navigation"
	"github.com/friendly-social/cli/internal/router"
	"github.com/friendly-social/cli/internal/screen"
	"github.com/friendly-social/cli/internal/ui"
)

// Screen is a model of home screen.
type Screen struct {
	cursor  int
	buttons []ui.Button

	width  int
	height int
}

// New returns new initial model of home screen.
func New() Screen {
	return Screen{
		buttons: []ui.Button{
			ui.NewButton("Back", func() tea.Msg {
				return screen.ChangeMsg{NewType: screen.TypeAuth}
			}),
			ui.NewButton("Quit", tea.Quit),
		},
	}
}

func (s Screen) ID() screen.Type {
	return screen.TypeHome
}

func (s Screen) getSelected() tea.Model {
	return s.buttons[s.cursor]
}

func (s Screen) setSelected(m tea.Model) {
	s.buttons[s.cursor] = m.(ui.Button)
}

func (s Screen) Init() tea.Cmd {
	return func() tea.Msg {
		return router.TargetMsg{Type: s.ID(), Inner: navigation.SelectMsg{}}
	}
}

func (s Screen) Update(msg tea.Msg) (screen.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
	case navigation.MoveMsg:
		var first, second tea.Model
		cmds := make([]tea.Cmd, 2)

		first, cmds[0] = s.getSelected().Update(navigation.UnselectMsg{})
		s.setSelected(first)

		switch msg.Direction {
		case navigation.DirectionDown:
			s.cursor = min(s.cursor+1, len(s.buttons)-1)
		case navigation.DirectionUp:
			s.cursor = max(s.cursor-1, 0)
		}

		second, cmds[1] = s.getSelected().Update(navigation.SelectMsg{})
		s.setSelected(second)

		return s, tea.Sequence(cmds...)
	}

	var cmd tea.Cmd
	var model tea.Model

	model, cmd = s.getSelected().Update(msg)
	s.setSelected(model)
	return s, cmd
}

func (s Screen) View() string {
	buttonViews := make([]string, len(s.buttons))
	for i, button := range s.buttons {
		cursor := ""
		if s.cursor == i {
			cursor = "-> "
		}

		buttonViews[i] = lipgloss.JoinHorizontal(lipgloss.Left, cursor, button.View())
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		"home screen placeholder...",
		"",
		lipgloss.JoinVertical(lipgloss.Center, buttonViews...),
	)
}
