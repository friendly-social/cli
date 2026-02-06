package feed

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/friendly-social/cli/internal/app"
	"github.com/friendly-social/cli/internal/auth"
	sdk "github.com/friendly-social/golang-sdk"
)

type Model struct {
	auth *sdk.Authorization
}

func New() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				return m, func() tea.Msg {
					return app.ScreenChangedMsg{NewScreen: app.ScreenAuth}
				}
			}
		case auth.AuthorizedMsg:
			m.auth = msg.Auth
	}

	return m, nil
}

func (m Model) View() string {
	return fmt.Sprintf("id: %d\ntoken: %s\nhash: %s\n", m.auth.Id, m.auth.Token, m.auth.AccessHash)
}
