package root

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/friendly-social/cli/internal/feed"
	"github.com/friendly-social/cli/internal/registration"
	"github.com/friendly-social/cli/internal/session"
	sdk "github.com/friendly-social/golang-sdk"
)

type Model struct {
	width  int
	height int

	mode   session.VimMode
	active session.ActiveModel
	models map[session.ActiveModel]tea.Model
}

func New(client *sdk.Client) Model {
	return Model{
		active: session.ActiveModelRegistration,
		models: map[session.ActiveModel]tea.Model{
			session.ActiveModelRegistration: registration.New(client),
			session.ActiveModelFeed: feed.New(),
		},
	}
}

func (m Model) Init() tea.Cmd {
	cmds := make([]tea.Cmd, len(m.models))
	for i, model := range m.models {
		cmds[i] = model.Init()
	}

	return tea.Batch(append(cmds, tea.ClearScreen)...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			if m.mode == session.VimModeInsert {
				return m, func() tea.Msg {
					return session.VimModeChangedMsg{NewMode: session.VimModeNormal}
				}
			}
		case "i":
			if m.mode == session.VimModeNormal {
				return m, func() tea.Msg {
					return session.VimModeChangedMsg{NewMode: session.VimModeInsert}
				}
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		cmds := make([]tea.Cmd, len(m.models))
		for i, model := range m.models {
			m.models[i], cmds[i] = model.Update(msg)
		}

		return m, tea.Batch(cmds...)
	case session.VimModeChangedMsg:
		m.mode = msg.NewMode

		cmds := make([]tea.Cmd, len(m.models))
		for i, model := range m.models {
			m.models[i], cmds[i] = model.Update(msg)
		}

		return m, tea.Batch(cmds...)
	case session.ActiveModelChangedMsg:
		m.active = msg.NewModel
		return m, nil
	}

	var cmd tea.Cmd
	m.models[m.active], cmd = m.models[m.active].Update(msg)
	return m, cmd
}

func (m Model) View() string {
	header := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width).
		Border(lipgloss.InnerHalfBlockBorder(), false, false, true, false).
		Render("Friendly CLI")

	mode := "NORMAL"
	if m.mode == session.VimModeInsert {
		mode = "INSERT"
	}

	footer := lipgloss.NewStyle().
		Align(lipgloss.Right).
		Width(m.width).
		Border(lipgloss.InnerHalfBlockBorder(), true, false, false, false).
		Render(fmt.Sprintf("--- %s ---", mode))

	content := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height-lipgloss.Height(header)-lipgloss.Height(footer)).
		Align(lipgloss.Center, lipgloss.Center).
		Render(m.models[m.active].View())

	return lipgloss.JoinVertical(lipgloss.Top, header, content, footer)
}
