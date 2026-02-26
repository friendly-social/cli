package router

import tea "github.com/charmbracelet/bubbletea"

type BroadcastMsg struct {
	Inner tea.Msg
}
