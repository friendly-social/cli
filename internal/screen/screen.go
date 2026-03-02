package screen

import tea "github.com/charmbracelet/bubbletea"

type Type string

const (
	TypeAuth Type = "auth"
	TypeMain Type = "main"
)

type Model interface {
	ID() Type

	Init() tea.Cmd
	Update(tea.Msg) (Model, tea.Cmd)
	View() string
}
