package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var listUnselectedStyle = lipgloss.NewStyle().PaddingLeft(3)

// List represents collection of elements that you can select and interact with.
type List struct {
	cursor int
	items  []tea.Model
}

// NewList creates new List based on the list of items.
func NewList(items ...tea.Model) *List {
	return &List{
		items: items,
	}
}

func (l *List) Init() tea.Cmd {
	return nil
}

func (l *List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case MoveMsg:
		cmds := make([]tea.Cmd, 2)
		l.items[l.cursor], cmds[0] = l.items[l.cursor].Update(UnselectMsg{})

		if msg.Direction == DirectionDown {
			l.cursor = min(l.cursor+1, len(l.items)-1)
		}

		if msg.Direction == DirectionUp {
			l.cursor = max(l.cursor-1, 0)
		}

		l.items[l.cursor], cmds[1] = l.items[l.cursor].Update(SelectMsg{})
		return l, tea.Batch(cmds...)
	}

	var cmd tea.Cmd
	l.items[l.cursor], cmd = l.items[l.cursor].Update(msg)
	return l, cmd
}

func (l *List) View() string {
	views := make([]string, len(l.items))
	for i, input := range l.items {
		if l.cursor == i {
			views[i] = lipgloss.JoinHorizontal(lipgloss.Left, "-> ", input.View())
			continue
		}

		views[i] = listUnselectedStyle.Render(input.View())
	}

	return lipgloss.JoinVertical(lipgloss.Left, views...)
}
