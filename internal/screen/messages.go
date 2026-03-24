package screen

// ChangeMsg singals that router must change the current screen.
type ChangeMsg struct {
	NewType Type
}
