package ui

// Label represents simple string that is need to be a part of UI.
type Label struct {
	title string
}

// NewLabel returns new Label from string.
func NewLabel(title string) *Label {
	return &Label{
		title: title,
	}
}

func (l *Label) View() string {
	return l.title
}

// Set sets label's value.
func (l *Label) Set(value string) {
	l.title = value
}

// Value returns the content of Label.
func (l *Label) Value() string {
	return l.title
}
