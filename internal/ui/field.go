package ui

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Field is an abstraction over textinput.Model for embedding it into interface.
type Field struct {
	input textinput.Model
}

// NewField creates new Field based on provided textinput.Model.
func NewField(input textinput.Model) *Field {
	input.Blur()
	return &Field{
		input: input,
	}
}

func (t *Field) Init() tea.Cmd {
	return nil
}

func (t *Field) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case FocusMsg:
		return t, tea.Batch(
			t.input.Focus(),
			t.input.Cursor.SetMode(cursor.CursorBlink),
		)
	case UnfocusMsg:
		t.input.Blur()
		return t, t.input.Cursor.SetMode(cursor.CursorStatic)
	}

	var cmd tea.Cmd
	t.input, cmd = t.input.Update(msg)
	return t, cmd
}

func (t *Field) View() string {
	return t.input.View()
}

// Value returns current filled string.
func (t *Field) Value() string {
	return t.input.Value()
}
