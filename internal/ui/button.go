package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Button struct {
	focused  bool
	title    string
	enterCmd tea.Cmd

	focusedStyle  lipgloss.Style
	inactiveStyle lipgloss.Style
}

func NewButton(title string, enterCmd tea.Cmd) *Button {
	return &Button{
		focused:       false,
		title:         title,
		enterCmd:      enterCmd,
		focusedStyle:  lipgloss.NewStyle().Background(lipgloss.Color("#7f00ff")),
		inactiveStyle: lipgloss.NewStyle().Background(lipgloss.Color("#808080")),
	}
}

func (b *Button) Update(msg tea.Msg) (Focusable, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok && msg.String() == "enter" {
		return b, b.enterCmd
	}

	return b, nil
}

func (b *Button) View() string {
	if b.focused {
		return b.focusedStyle.Render(b.title)
	}

	return b.inactiveStyle.Render(b.title)
}

func (b *Button) Focus() {
	b.focused = true
}

func (b *Button) Blur() {
	b.focused = false
}
