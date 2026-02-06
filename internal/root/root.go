package root

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/friendly-social/cli/internal/app"
	"github.com/friendly-social/cli/internal/auth"
	"github.com/friendly-social/cli/internal/feed"
	"github.com/friendly-social/cli/internal/vim"
	sdk "github.com/friendly-social/golang-sdk"
)

type Model struct {
	width  int
	height int

	mode   vim.Mode
	screen app.Screen
	models map[app.Screen]tea.Model
}

func New(client *sdk.Client) Model {
	return Model{
		screen: app.ScreenAuth,
		models: map[app.Screen]tea.Model{
			app.ScreenAuth: auth.New(client),
			app.ScreenFeed: feed.New(),
		},
	}
}

func (m Model) broadcast(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.models))
	for i, model := range m.models {
		m.models[i], cmds[i] = model.Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m Model) Init() tea.Cmd {
	cmds := make([]tea.Cmd, len(m.models))
	for i, model := range m.models {
		cmds[i] = model.Init()
	}

	return tea.Sequence(tea.ClearScreen, tea.Batch(cmds...))
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			if m.mode == vim.ModeInsert {
				return m, func() tea.Msg {
					return vim.ModeChangedMsg{NewMode: vim.ModeNormal}
				}
			}
		case "i":
			if m.mode == vim.ModeNormal {
				return m, func() tea.Msg {
					return vim.ModeChangedMsg{NewMode: vim.ModeInsert}
				}
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, m.broadcast(msg)
	case vim.ModeChangedMsg:
		m.mode = msg.NewMode
		return m, m.broadcast(msg)
	case app.ScreenChangedMsg:
		m.screen = msg.NewScreen
		return m, nil
	case auth.AuthorizedMsg:
		return m, m.broadcast(msg)
	}

	var cmd tea.Cmd
	m.models[m.screen], cmd = m.models[m.screen].Update(msg)
	return m, cmd
}

func (m Model) View() string {
	header := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(m.width).
		Border(lipgloss.InnerHalfBlockBorder(), false, false, true, false).
		Render("Friendly CLI")

	mode := "NORMAL"
	if m.mode == vim.ModeInsert {
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
		Render(m.models[m.screen].View())

	return lipgloss.JoinVertical(lipgloss.Top, header, content, footer)
}
