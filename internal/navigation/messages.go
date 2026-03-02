package navigation

import tea "github.com/charmbracelet/bubbletea"

type Direction int

const (
	DirectionLeft Direction = iota
	DirectionRight
	DirectionDown
	DirectionUp
)

type MoveMsg struct {
	Direction Direction
}

type KeyMsg struct {
	Key tea.KeyType
}

type InteractMsg struct {
}

type FocusMsg struct {
}

type UnfocusMsg struct {
}

type SelectMsg struct {
}

type UnselectMsg struct {
}
