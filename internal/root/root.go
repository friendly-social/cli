package root

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/friendly-social/cli/internal/registration"
	"github.com/friendly-social/cli/internal/session"
	sdk "github.com/friendly-social/golang-sdk"
)

type Model struct {
	mode        session.VimMode
	activeModel session.ActiveModel

	registration tea.Model
}

func NewModel(client *sdk.Client) Model {
	return Model{
		activeModel:  session.ActiveModelRegistration,
		registration: registration.NewModel(client),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.ClearScreen,
		m.registration.Init(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			if m.mode == session.VimModeInsert {
				m.mode = session.VimModeNormal
				return m, func() tea.Msg {
					return session.VimModeChangedMsg{NewMode: m.mode}
				}
			}
		case "i":
			if m.mode == session.VimModeNormal {
				m.mode = session.VimModeInsert
				return m, func() tea.Msg {
					return session.VimModeChangedMsg{NewMode: m.mode}
				}
			}
		}
	}

	switch m.activeModel {
	case session.ActiveModelRegistration:
		newModel, cmd := m.registration.Update(msg)
		m.registration = newModel
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	var content string
	switch m.activeModel {
	case session.ActiveModelRegistration:
		content = m.registration.View()
	}

	var s strings.Builder
	fmt.Fprintf(&s, "%s\n", content)

	mode := "NORMAL"
	if m.mode == session.VimModeInsert {
		mode = "INSERT"
	}

	fmt.Fprintf(&s, "\n--- %s ---\n", mode)
	return s.String()
}
