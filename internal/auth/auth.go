package auth

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/friendly-social/cli/internal/app"
	"github.com/friendly-social/cli/internal/ui"
	"github.com/friendly-social/cli/internal/vim"
	sdk "github.com/friendly-social/golang-sdk"
)

type AuthorizedMsg struct {
	Auth *sdk.Authorization
}

type Model struct {
	mode   vim.Mode
	client *sdk.Client

	cursor  int
	inputs  []*ui.TextInput
	buttons []*ui.Button
}

func New(client *sdk.Client) Model {
	ni := textinput.New()
	ni.Placeholder = "Nickname"
	ni.CharLimit = 256
	ni.Prompt = ""
	ni.Focus()

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

	m := Model{
		cursor: 0,
		client: client,
		inputs: []*ui.TextInput{
			ui.NewTextInput(ni),
			ui.NewTextInput(di),
			ui.NewTextInput(ii),
			ui.NewTextInput(si),
		},
	}

	m.buttons = []*ui.Button{
		ui.NewButton("Submit", m.authCmd()),
		ui.NewButton("Quit", tea.Quit),
	}

	return m
}

func (m Model) getFocusable(cursor int) ui.Focusable {
	if cursor < len(m.inputs) {
		return m.inputs[cursor]
	}

	return m.buttons[cursor-len(m.inputs)]
}

func (m Model) Init() tea.Cmd {
	return m.initCmd
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.mode != vim.ModeNormal {
			break
		}

		switch msg.String() {
		case "k", "up":
			if m.cursor > 0 {
				m.getFocusable(m.cursor).Blur()
				m.cursor--
				m.getFocusable(m.cursor).Focus()
				return m, nil
			}
		case "j", "down":
			if m.cursor < len(m.inputs)+len(m.buttons)-1 {
				m.getFocusable(m.cursor).Blur()
				m.cursor++
				m.getFocusable(m.cursor).Focus()
				return m, nil
			}
		}
	case AuthorizedMsg:
		return m, func() tea.Msg {
			return app.ScreenChangedMsg{NewScreen: app.ScreenFeed}
		}
	case app.ErrorMsg:
		switch msg.Error {
		case sdk.ErrNicknameLengthMustBeLessThan256:
			//...
		}
	case vim.ModeChangedMsg:
		m.mode = msg.NewMode
	}

	_, cmd := m.getFocusable(m.cursor).Update(msg)
	return m, cmd
}

func (m Model) View() string {
	inputViews := make([]string, len(m.inputs))
	for i, input := range m.inputs {
		cursor := ""
		if m.cursor == i {
			cursor = "-> "
		}

		inputViews[i] = lipgloss.JoinHorizontal(lipgloss.Center, cursor, input.View())
	}

	buttonViews := make([]string, len(m.buttons))
	for i, button := range m.buttons {
		cursor := ""
		if m.cursor == len(m.inputs)+i {
			cursor = "-> "
		}

		buttonViews[i] = lipgloss.JoinHorizontal(lipgloss.Center, cursor, button.View())
	}

	return lipgloss.JoinVertical(
		lipgloss.Center,
		"please introduce yourself!",
		"",
		lipgloss.JoinVertical(lipgloss.Center, inputViews...),
		"",
		lipgloss.JoinVertical(lipgloss.Center, buttonViews...),
	)
}
