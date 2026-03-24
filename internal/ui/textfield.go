package ui

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/friendly-social/cli/internal/navigation"
)

// TextField is an abstraction over textinput.Model for embedding it to navigation interface.
type TextField struct {
	input textinput.Model
}

// NewTextField creates new TextField based on provided textinput.Model.
func NewTextField(input textinput.Model) TextField {
	input.Blur()
	return TextField{
		input: input,
	}
}

func (t TextField) Init() tea.Cmd {
	return nil
}

func (t TextField) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case navigation.FocusMsg:
		return t, tea.Batch(
			t.input.Focus(),
			t.input.Cursor.SetMode(cursor.CursorBlink),
		)
	case navigation.UnfocusMsg:
		t.input.Blur()
		return t, t.input.Cursor.SetMode(cursor.CursorStatic)
	}

	var cmd tea.Cmd
	t.input, cmd = t.input.Update(msg)
	return t, cmd
}

func (t TextField) View() string {
	return t.input.View()
}

// Value returns current filled string.
func (t TextField) Value() string {
	return t.input.Value()
}
