package registration

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/friendly-social/cli/internal/session"
	sdk "github.com/friendly-social/golang-sdk"
)

type Model struct {
	client      *sdk.Client
	finalOutput string

	cursor int
	mode   session.VimMode
	inputs []textinput.Model
}

func New(client *sdk.Client) Model {
	ni := textinput.New()
	ni.Placeholder = "Nickname"
	ni.CharLimit = 256
	ni.Prompt = ""
	ni.Focus()

	di := textinput.New()
	di.Placeholder = "Description"
	di.CharLimit = 1024
	di.Prompt = ""

	ii := textinput.New()
	ii.Placeholder = "Interests"
	ii.Prompt = ""

	si := textinput.New()
	si.Placeholder = "Social Link"
	si.CharLimit = 1024
	si.Prompt = ""

	return Model{
		cursor: 0,
		client: client,
		inputs: []textinput.Model{ni, di, ii, si},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "k":
			if m.mode != session.VimModeNormal {
				break
			}

			previousCursor := m.cursor
			if msg.String() == "j" && m.cursor+1 < len(m.inputs)+2 {
				m.cursor++
			}

			if msg.String() == "k" && m.cursor > 0 {
				m.cursor--
			}

			if previousCursor < len(m.inputs) {
				m.inputs[previousCursor].Blur()
			}

			if m.cursor < len(m.inputs) {
				return m, m.inputs[m.cursor].Focus()
			}

			return m, nil
		case "enter":
			switch m.cursor {
			case len(m.inputs):
				/*
					nickname, _ := sdk.NewNickname(m.inputs[0].Value())
					description, _ := sdk.NewUserDescription(m.inputs[1].Value())
					social, _ := sdk.NewSocialLink(m.inputs[3].Value())

					interests := make([]sdk.Interest, 0)
					for interestStr := range strings.SplitSeq(m.inputs[2].Value(), ",") {
						interest, _ := sdk.NewInterest(strings.TrimSpace(interestStr))
						interests = append(interests, interest)
					}

					auth, _ := m.client.Generate(context.Background(), nickname, description, interests, nil, social)
					m.finalOutput = fmt.Sprintf("Id: %d\nToken: %s\nHash: %s", auth.Id, auth.Token, auth.AccessHash)
				*/

				return m, func() tea.Msg {
					return session.ActiveModelChangedMsg{NewModel: session.ActiveModelFeed}
				}
			case len(m.inputs) + 1:
				return m, tea.Quit
			}
		}
	case session.VimModeChangedMsg:
		m.mode = msg.NewMode
		if m.cursor >= len(m.inputs) {
			break
		}

		if m.mode == session.VimModeInsert {
			cmd := m.inputs[m.cursor].Cursor.SetMode(cursor.CursorBlink)
			return m, tea.Batch(cmd, m.inputs[m.cursor].Cursor.BlinkCmd())
		}

		return m, m.inputs[m.cursor].Cursor.SetMode(cursor.CursorStatic)
	}

	if m.mode != session.VimModeInsert || m.cursor >= len(m.inputs) {
		return m, nil
	}

	var cmd tea.Cmd
	m.inputs[m.cursor], cmd = m.inputs[m.cursor].Update(msg)
	return m, cmd
}

func (m Model) View() string {
	var s strings.Builder
	fmt.Fprintf(&s, "%s\n\n", "please introduce yourself!")

	for i, input := range m.inputs {
		cursor := "  "
		if m.cursor == i {
			cursor = "->"
		}

		fmt.Fprintf(&s, "%s %s\n", cursor, input.View())
	}

	cursor := "  "
	if m.cursor == len(m.inputs) {
		cursor = "->"
	}
	fmt.Fprintf(&s, "\n%s [ %s ]\n", cursor, "Submit")

	cursor = "  "
	if m.cursor == len(m.inputs)+1 {
		cursor = "->"
	}
	fmt.Fprintf(&s, "%s [ %s ]\n", cursor, "Quit")

	if m.finalOutput != "" {
		fmt.Fprintf(&s, "\n\n%s\n\n", m.finalOutput)
	}

	return s.String()
}
