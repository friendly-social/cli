package feed

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/friendly-social/cli/internal/session"
)

type Model struct {
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
					return session.ActiveModelChangedMsg{NewModel: session.ActiveModelRegistration}
				}
			}
	}

	return m, nil
}

func (m Model) View() string {
	return "work in progress..."
}
