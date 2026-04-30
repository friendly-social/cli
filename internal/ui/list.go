package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type List struct {
	cursor int
	fields []tea.Model
}

func NewList(fields ...tea.Model) *List {
	return &List{
		fields: fields,
	}
}

func (l *List) Init() tea.Cmd {
	return nil
}

func (l *List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case MoveMsg:
		cmds := make([]tea.Cmd, 2)
		l.fields[l.cursor], cmds[0] = l.fields[l.cursor].Update(UnselectMsg{})

		if msg.Direction == DirectionDown {
			l.cursor = min(l.cursor+1, len(l.fields)-1)
		}

		if msg.Direction == DirectionUp {
			l.cursor = max(l.cursor-1, 0)
		}

		l.fields[l.cursor], cmds[1] = l.fields[l.cursor].Update(SelectMsg{})
		return l, tea.Batch(cmds...)
	}

	var cmd tea.Cmd
	l.fields[l.cursor], cmd = l.fields[l.cursor].Update(msg)
	return l, cmd
}

func (l *List) View() string {
	views := make([]string, len(l.fields))
	for i, input := range l.fields {
		cursor := ""
		if l.cursor == i {
			cursor = "-> "
		}

		views[i] = lipgloss.JoinHorizontal(lipgloss.Left, cursor, input.View())
	}

	return lipgloss.JoinVertical(lipgloss.Center, views...)
}
