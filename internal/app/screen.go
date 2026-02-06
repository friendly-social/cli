package app

type Screen int

const (
	ScreenAuth Screen = iota
	ScreenFeed
)

type ScreenChangedMsg struct {
	NewScreen Screen
}
