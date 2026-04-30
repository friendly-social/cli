package auth

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/friendly-social/cli/internal/router"
	"github.com/friendly-social/cli/internal/screen"
	"github.com/friendly-social/cli/internal/ui"
)

// Screen is a model of authentication screen.
type Screen struct {
	service *Service

	content struct {
		list   *ui.List
		status *ui.Label

		fields struct {
			nickname    *ui.Field
			description *ui.Field
			interests   *ui.Field
			social      *ui.Field
		}

		buttons struct {
			submit *ui.Button
			exit   *ui.Button
		}
	}
}

func field(label string, limit int) *ui.Field {
	field := textinput.New()
	field.Placeholder = label
	field.CharLimit = limit
	field.Prompt = ""
	return ui.NewField(field)
}

// New returns new initial model of authentication screen.
func New(service *Service) Screen {
	result := Screen{
		service: service,
	}

	result.content.fields.nickname = field("Nickname", 256)
	result.content.fields.description = field("Description", 1024)
	result.content.fields.interests = field("Interests", 0)
	result.content.fields.social = field("Social Link", 1024)

	result.content.buttons.exit = ui.NewButton("Exit", tea.Quit)
	result.content.buttons.submit = ui.NewButton("Submit", func() tea.Msg {
		result.content.status.Set("authenticating...")
		return service.auth(
			result.content.fields.nickname.Value(),
			result.content.fields.description.Value(),
			result.content.fields.interests.Value(),
			result.content.fields.social.Value())
	})

	result.content.status = ui.NewLabel("")
	result.content.list = ui.NewList(
		result.content.fields.nickname,
		result.content.fields.description,
		result.content.fields.interests,
		result.content.fields.social,
		result.content.buttons.submit,
		result.content.buttons.exit)

	return result
}

func (s Screen) ID() screen.Type {
	return screen.TypeAuth
}

func (s Screen) Init() tea.Cmd {
	return tea.Sequence(
		func() tea.Msg {
			return s.service.loadAuth()
		},
		func() tea.Msg {
			return router.TargetMsg{Type: s.ID(), Inner: ui.SelectMsg{}}
		})
}

func (s Screen) Update(msg tea.Msg) (screen.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case AuthMsg:
		return s, func() tea.Msg {
			return screen.ChangeMsg{NewType: screen.TypeHome}
		}
	case screen.ErrorMsg:
		s.content.status.Set(msg.Value.Error())
		return s, nil
	}

	_, cmd := s.content.list.Update(msg)
	return s, cmd
}

func (s Screen) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		"please introduce yourself!",
		"",
		s.content.list.View(),
		"",
		s.content.status.View(),
	)
}
