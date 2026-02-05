package session

import sdk "github.com/friendly-social/golang-sdk"

type VimMode int

const (
	VimModeNormal VimMode = iota
	VimModeInsert
)

type ActiveModel int

const (
	ActiveModelRegistration ActiveModel = iota
)

type Context struct {
	Auth *sdk.Authorization
	Mode VimMode
}

type VimModeChangedMsg struct {
	NewMode VimMode
}

type ActiveModelChangedMsg struct {
	NewModel ActiveModel
}
