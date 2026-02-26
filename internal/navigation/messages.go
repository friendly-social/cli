package navigation

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

type FocusMsg struct {
}

type UnfocusMsg struct {
}

type InteractMsg struct {
}
