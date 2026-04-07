package auth

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/friendly-social/cli/internal/navigation"
	"github.com/friendly-social/cli/internal/router"
	"github.com/friendly-social/cli/internal/screen"
	"github.com/friendly-social/cli/internal/system"
	"github.com/friendly-social/cli/internal/ui"
)

func (s Screen) ID() screen.Type {
	return screen.TypeAuth
}

func (s Screen) Init() tea.Cmd {
	return tea.Sequence(s.initCmd(), func() tea.Msg {
		return router.TargetMsg{Type: s.ID(), Inner: navigation.SelectMsg{}}
	})
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
			s.cursor = min(s.cursor+1, len(s.fields)+len(s.buttons)-1)
		case navigation.DirectionUp:
			s.cursor = max(s.cursor-1, 0)
		}

		second, cmds[1] = s.getSelected().Update(navigation.SelectMsg{})
		s.setSelected(second)

		return s, tea.Sequence(cmds...)
	case AuthMsg:
		return s, func() tea.Msg {
			return screen.ChangeMsg{NewType: screen.TypeHome}
		}
	case system.ErrorMsg:
		var newLabel tea.Model
		newLabel, cmd := s.errorLabel.Update(ui.LabelChangeMsg{Value: msg.Value.Error()})

		s.errorLabel = newLabel.(ui.Label)
		return s, cmd
	}

	model, cmd := s.getSelected().Update(msg)
	s.setSelected(model)
	return s, cmd
}

func (s Screen) View() string {
	inputViews := make([]string, len(s.fields))
	for i, input := range s.fields {
		cursor := ""
		if s.cursor == i {
			cursor = "-> "
		}

		inputViews[i] = lipgloss.JoinHorizontal(lipgloss.Left, cursor, input.View())
	}

	buttonViews := make([]string, len(s.buttons))
	for i, button := range s.buttons {
		cursor := ""
		if s.cursor == len(s.fields)+i {
			cursor = "-> "
		}

		buttonViews[i] = lipgloss.JoinHorizontal(lipgloss.Left, cursor, button.View())
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		"please introduce yourself!",
		"",
		lipgloss.JoinVertical(lipgloss.Center, inputViews...),
		"",
		lipgloss.JoinVertical(lipgloss.Center, buttonViews...),
		"",
		s.errorLabel.Value(),
	)
}
