package session

type VimMode int

const (
	VimModeNormal VimMode = iota
	VimModeInsert
)

type ActiveModel int

const (
	ActiveModelRegistration ActiveModel = iota
	ActiveModelFeed 
)

type VimModeChangedMsg struct {
	NewMode VimMode
}

type ActiveModelChangedMsg struct {
	NewModel ActiveModel
}
