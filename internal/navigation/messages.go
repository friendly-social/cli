package navigation

type Direction int

const (
	DirectionLeft Direction = iota
	DirectionRight
	DirectionDown
	DirectionUp
)

type MovedMsg struct {
	Direction Direction
}

type FocusedMsg struct {
}

type UnfocusedMsg struct {
}

type InteractedMsg struct {
}
