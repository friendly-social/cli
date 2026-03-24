package navigation

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// VimMode represents possible modes for Vim motions.
type VimMode string

const (
	VimModeNormal VimMode = "NORMAL"
	VimModeInsert VimMode = "INSERT"
)

// VimWrapper translates raw tea.KeyMsgs to navigation messages using Vim motions driven logic.
type VimWrapper struct {
	mode  VimMode
	model tea.Model

	width  int
	height int
}

// NewVimWrapper creates new VimWrapper based on provided model.
func NewVimWrapper(model tea.Model) VimWrapper {
	return VimWrapper{
		model: model,
		mode:  VimModeNormal,
	}
}

func (w VimWrapper) Init() tea.Cmd {
	return w.model.Init()
}

func (w VimWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case VimModeNormal:
			switch msg.String() {
			case "i":
				w.mode = VimModeInsert
				return w, func() tea.Msg {
					return FocusMsg{}
				}
			case "h", "left":
				return w, func() tea.Msg {
					return MoveMsg{Direction: DirectionLeft}
				}
			case "j", "down":
				return w, func() tea.Msg {
					return MoveMsg{Direction: DirectionDown}
				}
			case "k", "up":
				return w, func() tea.Msg {
					return MoveMsg{Direction: DirectionUp}
				}
			case "l", "right":
				return w, func() tea.Msg {
					return MoveMsg{Direction: DirectionRight}
				}
			case "enter":
				return w, func() tea.Msg {
					return InteractMsg{}
				}
			default:
				return w, func() tea.Msg {
					return KeyMsg{Value: msg.String()}
				}
			}
		case VimModeInsert:
			switch msg.String() {
			case "esc", "ctrl+c":
				w.mode = VimModeNormal
				return w, func() tea.Msg {
					return UnfocusMsg{}
				}
			}
		}
	}

	var cmd tea.Cmd
	w.model, cmd = w.model.Update(msg)
	return w, cmd
}

func (w VimWrapper) footer() string {
	return lipgloss.NewStyle().
		Align(lipgloss.Right).
		Width(w.width).
		Border(lipgloss.InnerHalfBlockBorder(), true, false, false, false).
		Render(fmt.Sprintf("--- %s ---", w.mode))
}

func (w VimWrapper) View() string {
	footer := w.footer()

	content := lipgloss.NewStyle().
		Width(w.width).
		Height(w.height - lipgloss.Height(footer)).
		Render(w.model.View())

	return lipgloss.JoinVertical(lipgloss.Top, content, footer)
}
