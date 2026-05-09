package profile

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/friendly-social/cli/internal/router"
	"github.com/friendly-social/cli/internal/screen"
	"github.com/friendly-social/cli/internal/screen/register"
	"github.com/friendly-social/cli/internal/ui"
)

type Screen struct {
	service *Service

	content struct {
		label *ui.Label
		list  *ui.List

		button struct {
			home *ui.Button
		}
	}
}

func New(service *Service) Screen {
	result := Screen{
		service: service,
	}

	result.content.label = ui.NewLabel("")
	result.content.button.home = ui.NewButton("Home", func() tea.Msg {
		return screen.ChangeMsg{NewType: screen.TypeHome}
	})

	result.content.list = ui.NewList(
		result.content.button.home)

	return result
}

func (Screen) ID() screen.Type {
	return screen.TypeProfile
}

func (s Screen) Init() tea.Cmd {
	return func() tea.Msg {
		return router.TargetMsg{Type: s.ID(), Inner: ui.SelectMsg{}}
	}
}

func (s Screen) Update(msg tea.Msg) (screen.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case register.AuthMsg:
		s.content.label.Set("loading...")
		return s, func() tea.Msg {
			self, err := s.service.GetDetails(msg.User)
			if err != nil {
				s.content.label.Set(fmt.Sprintf("error loading profile: %s", err.Error()))
				return screen.TickMsg{}
			}

			var interests strings.Builder
			interestsSlice := self.Interests.Value()
			for i, interest := range interestsSlice {
				interests.WriteString(interest.Value())
				if i != len(interestsSlice)-1 {
					interests.WriteString(", ")
				}
			}

			s.content.label.Set(fmt.Sprintf(
				"your logged in profile:\nnickname: %s\ndescription: %s\ninterests: %s\nsocial link: %s",
				self.Nickname.Value(), self.Description.Value(), interests.String(), self.SocialLink.Value()))

			return screen.TickMsg{}
		}
	}

	_, cmd := s.content.list.Update(msg)
	return s, cmd
}

func (s Screen) View() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		s.content.label.View(),
		"",
		s.content.list.View())
}
