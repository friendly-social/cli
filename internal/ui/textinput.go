package ui

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/friendly-social/cli/internal/vim"
)

type TextInput struct {
	mode  vim.Mode
	input textinput.Model
}

func NewTextInput(input textinput.Model) *TextInput {
	return &TextInput{
		mode:  vim.ModeNormal,
		input: input,
	}
}

func (t *TextInput) Update(msg tea.Msg) (Focusable, tea.Cmd) {
	switch msg := msg.(type) {
	case vim.ModeChangedMsg:
		t.mode = msg.NewMode
		if t.mode == vim.ModeInsert {
			cmd := t.input.Cursor.SetMode(cursor.CursorBlink)
			return t, tea.Batch(cmd, t.input.Cursor.BlinkCmd())
		}

		return t, t.input.Cursor.SetMode(cursor.CursorStatic)
	case tea.KeyMsg:
		if t.mode != vim.ModeInsert {
			return t, nil
		}

		var cmd tea.Cmd
		t.input, cmd = t.input.Update(msg)
		return t, cmd
	}

	return t, nil
}

func (t *TextInput) View() string {
	return t.input.View()
}

func (t *TextInput) Focus() {
	t.input.Focus()
}

func (t *TextInput) Blur() {
	t.input.Blur()
}

func (t *TextInput) Raw() *textinput.Model {
	return &t.input
}
