package ui

import tea "github.com/charmbracelet/bubbletea"

// Label represents simple string that is need to be a part of UI.
type Label struct {
	title string
}

// LabelChangeMsg signals Label to change its content.
type LabelChangeMsg struct {
	Value string
}

// NewLabel returns new Label from string.
func NewLabel(title string) Label {
	return Label{
		title: title,
	}
}

func (l Label) Init() tea.Cmd {
	return nil
}

func (l Label) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case LabelChangeMsg:
		l.title = msg.Value
	}

	return l, nil
}

func (l Label) View() string {
	return l.title
}

// Value returns content of Label.
func (l Label) Value() string {
	return l.title
}
