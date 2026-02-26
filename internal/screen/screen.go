package screen

import tea "github.com/charmbracelet/bubbletea"

type Type int

const (
	TypeAuth Type = iota
	TypeFeed
)

type Model interface {
	ID() Type

	Init() tea.Cmd
	Update(tea.Msg) (Model, tea.Cmd)
	View() string
}

type ChangeMsg struct {
	NewType Type
}
