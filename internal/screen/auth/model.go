package auth

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/friendly-social/cli/internal/ui"
	sdk "github.com/friendly-social/golang-sdk"
)

// Screen is a model of authentication screen.
type Screen struct {
	client *sdk.Client

	cursor     int
	fields     []ui.TextField
	buttons    []ui.Button
	errorLabel ui.Label

	width  int
	height int
}

// New returns new initial model of authentication screen.
func New() Screen {
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

	screen := Screen{
		client: sdk.NewClient(),
		fields: []ui.TextField{
			ui.NewTextField(ni),
			ui.NewTextField(di),
			ui.NewTextField(ii),
			ui.NewTextField(si),
		},
	}

	screen.buttons = []ui.Button{
		ui.NewButton("Submit", screen.submitCmd()),
		ui.NewButton("Exit", tea.Quit),
	}

	return screen
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
