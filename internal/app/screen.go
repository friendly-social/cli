package app

import tea "github.com/charmbracelet/bubbletea"

type ScreenType int

const (
	ScreenAuth ScreenType = iota
	ScreenFeed
)

type Screen interface {
	ID() ScreenType

	Init() tea.Cmd
	Update(tea.Msg) (Screen, tea.Cmd)
	View() string
}
