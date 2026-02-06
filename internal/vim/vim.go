package vim

type Mode int

const (
	ModeNormal Mode = iota
	ModeInsert
)

type ModeChangedMsg struct {
	NewMode Mode
}
