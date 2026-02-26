package feed

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/friendly-social/cli/internal/screen"
)

type Screen struct {
	width  int
	height int
}

func NewScreen() Screen {
	return Screen{}
}

func (s Screen) ID() screen.Type {
	return screen.TypeFeed
}

func (s Screen) Init() tea.Cmd {
	return nil
}

func (s Screen) Update(msg tea.Msg) (screen.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			return s, func() tea.Msg {
				return screen.ChangeMsg{NewType: screen.TypeAuth}
			}
		}
	}

	return s, nil
}

func (s Screen) View() string {
	return "\nfeed"
}
