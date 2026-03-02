package router

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/friendly-social/cli/internal/screen"
)

type BroadcastMsg struct {
	Inner tea.Msg
}

type TargetMsg struct {
	Type screen.Type
	Inner tea.Msg
}
