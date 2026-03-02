package auth

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/friendly-social/cli/internal/navigation"
	"github.com/friendly-social/cli/internal/router"
	"github.com/friendly-social/cli/internal/screen"
	"github.com/friendly-social/cli/internal/ui"
)

type Screen struct {
	fields  []ui.TextField
	buttons []ui.Button
	cursor  int

	width  int
	height int
}

func NewScreen() Screen {
	ni := textinput.New()
	ni.Placeholder = "Nickname"
	ni.CharLimit = 256
	ni.Prompt = ""

	di := textinput.New()
	di.Placeholder = "Description"
	di.CharLimit = 1024
	di.Prompt = ""

	ii := textinput.New()
	ii.Placeholder = "Interests"
	ii.Prompt = ""

	si := textinput.New()
	si.Placeholder = "Social Link"
	si.CharLimit = 1024
	si.Prompt = ""

	return Screen{
		fields: []ui.TextField{
			ui.NewTextField(ni),
			ui.NewTextField(di),
			ui.NewTextField(ii),
			ui.NewTextField(si),
		},
		buttons: []ui.Button{
			// temp
			ui.NewButton("Submit", func() tea.Msg {
				return screen.ChangeMsg{NewType: screen.TypeMain}
			}),
			ui.NewButton("Exit", tea.Quit),
		},
	}
}

func (s Screen) ID() screen.Type {
	return screen.TypeAuth
}

func (s Screen) getSelected() tea.Model {
	if s.cursor < len(s.fields) {
		return s.fields[s.cursor]
	}

	return s.buttons[s.cursor-len(s.fields)]
}

func (s Screen) setSelected(m tea.Model) {
	if s.cursor < len(s.fields) {
		s.fields[s.cursor] = m.(ui.TextField)
		return
	}

	s.buttons[s.cursor-len(s.fields)] = m.(ui.Button)
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
			s.cursor = min(s.cursor+1, len(s.fields)+len(s.buttons)-1)
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
	)
}
