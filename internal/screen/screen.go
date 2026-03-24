package screen

import tea "github.com/charmbracelet/bubbletea"

// Type represents type of the current screen and serves as an identificator.
type Type string

const (
	TypeAuth Type = "auth"
	TypeHome Type = "home"
)

// Model represents Screen which is basically an extended tea.Model.
type Model interface {
	ID() Type

	Init() tea.Cmd
	Update(tea.Msg) (Model, tea.Cmd)
	View() string
}
