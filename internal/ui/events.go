package ui

// Direction represents possible moving directions.
type Direction int

const (
	DirectionLeft Direction = iota
	DirectionRight
	DirectionDown
	DirectionUp
)

// MoveMsg shows that user wants to move to a different component in some direction.
type MoveMsg struct {
	Direction Direction
}

// InteractMsg shows that user wants to interact with some component.
type InteractMsg struct {
}

// FocusMsg shows that user wants to focus on some component for interacting with it further.
type FocusMsg struct {
}

// UnfocucMsg shows that user no longer wants to focus on current component.
type UnfocusMsg struct {
}

// SelectMsg shows that user wants to select some component for interacting with it further.
type SelectMsg struct {
}

// UnselectMsg shows that user no longer wants current component to be selected.
type UnselectMsg struct {
}
