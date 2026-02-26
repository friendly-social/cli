package vim

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/friendly-social/cli/internal/navigation"
)

type Mode string

const (
	ModeNormal Mode = "NORMAL"
	ModeInsert Mode = "INSERT"
)

type Wrapper struct {
	mode  Mode
	model tea.Model

	width  int
	height int
}

func NewWrapper(model tea.Model) Wrapper {
	return Wrapper{
		model: model,
		mode:  ModeNormal,
	}
}

func (w Wrapper) Init() tea.Cmd {
	return nil
}

func (w Wrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		w.width = msg.Width
		w.height = msg.Height
		msg.Height -= lipgloss.Height(w.footer())

		var cmd tea.Cmd
		w.model, cmd = w.model.Update(msg)
		return w, cmd
	case tea.KeyMsg:
		switch w.mode {
		case ModeNormal:
			switch msg.String() {
			case "i":
				w.mode = ModeInsert
				return w, func() tea.Msg {
					return navigation.FocusedMsg{}
				}
			case "h":
				return w, func() tea.Msg {
					return navigation.MovedMsg{Direction: navigation.DirectionLeft}
				}
			case "j":
				return w, func() tea.Msg {
					return navigation.MovedMsg{Direction: navigation.DirectionDown}
				}
			case "k":
				return w, func() tea.Msg {
					return navigation.MovedMsg{Direction: navigation.DirectionUp}
				}
			case "l":
				return w, func() tea.Msg {
					return navigation.MovedMsg{Direction: navigation.DirectionRight}
				}
			case "enter":
				return w, func() tea.Msg {
					return navigation.InteractedMsg{}
				}
			default:
				return w, nil
			}
		case ModeInsert:
			switch msg.String() {
			case "esc", "ctrl+c":
				w.mode = ModeNormal
				return w, func() tea.Msg {
					return navigation.UnfocusedMsg{}
				}
			}
		}
	}

	var cmd tea.Cmd
	w.model, cmd = w.model.Update(msg)
	return w, cmd
}

func (w Wrapper) footer() string {
	return lipgloss.NewStyle().
		Align(lipgloss.Right).
		Width(w.width).
		Border(lipgloss.InnerHalfBlockBorder(), true, false, false, false).
		Render(fmt.Sprintf("--- %s ---", w.mode))
}

func (w Wrapper) View() string {
	footer := w.footer()

	content := lipgloss.NewStyle().
		Width(w.width).
		Height(w.height - lipgloss.Height(footer)).
		Render(w.model.View())

	return lipgloss.JoinVertical(lipgloss.Top, content, footer)
}
