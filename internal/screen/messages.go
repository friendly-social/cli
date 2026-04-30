package screen

// ChangeMsg singals that router must change the current screen.
type ChangeMsg struct {
	NewType Type
}

// ErrorMsg is a message that represents an error occured in program.
type ErrorMsg struct {
	Value error
}
