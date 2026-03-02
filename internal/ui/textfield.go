package ui

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/friendly-social/cli/internal/navigation"
)

type TextField struct {
	field textinput.Model
}

func NewTextField(input textinput.Model) TextField {
	return TextField{
		field: input,
	}
}

func (t TextField) Init() tea.Cmd {
	return t.field.Cursor.SetMode(cursor.CursorHide)
}

func (t TextField) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case navigation.FocusMsg:
		t.field.Focus()
		return t, tea.Batch(
			t.field.Cursor.SetMode(cursor.CursorBlink),
			t.field.Cursor.BlinkCmd(),
		)
	case navigation.UnfocusMsg:
		t.field.Blur()
		return t, t.field.Cursor.SetMode(cursor.CursorStatic)
	}

	var cmd tea.Cmd
	t.field, cmd = t.field.Update(msg)
	return t, cmd
}

func (t TextField) View() string {
	return t.field.View()
}

func (t TextField) Value() string {
	return t.field.Value()
}
