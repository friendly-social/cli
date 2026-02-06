package ui

import tea "github.com/charmbracelet/bubbletea"

type Focusable interface {
	Update(msg tea.Msg) (Focusable, tea.Cmd)
	View() string

	Focus() 
	Blur() 
}
