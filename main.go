package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"friendly/cli/friendly"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			MarginBottom(1)

	focusedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4"))

	blurredStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#808080"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000"))

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			MarginTop(1)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2)

	selectedBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#00FF00")).
			Padding(1, 2)
)

type view int

const (
	viewRegister view = iota
	viewMenu
	viewFeed
	viewNetwork
	viewProfile
	viewAddFriend
	viewGenerateToken
)

type mode int

const (
	modeNormal mode = iota
	modeInsert
)

type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Help     key.Binding
	Quit     key.Binding
	Enter    key.Binding
	Back     key.Binding
	NextView key.Binding
	PrevView key.Binding
	Insert   key.Binding
	Normal   key.Binding
	Refresh  key.Binding
	Select   key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Enter, k.Back, k.Insert},
		{k.NextView, k.PrevView, k.Refresh},
		{k.Help, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("k/↑", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("j/↓", "down"),
	),
	Left: key.NewBinding(
		key.WithKeys("h", "left"),
		key.WithHelp("h/←", "left"),
	),
	Right: key.NewBinding(
		key.WithKeys("l", "right"),
		key.WithHelp("l/→", "right"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select/submit"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back/normal mode"),
	),
	Insert: key.NewBinding(
		key.WithKeys("i", "a"),
		key.WithHelp("i/a", "insert mode"),
	),
	Normal: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "normal mode"),
	),
	NextView: key.NewBinding(
		key.WithKeys("tab", "n"),
		key.WithHelp("tab/n", "next view"),
	),
	PrevView: key.NewBinding(
		key.WithKeys("shift+tab", "p"),
		key.WithHelp("shift+tab/p", "prev view"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
	Select: key.NewBinding(
		key.WithKeys(" ", "x"),
		key.WithHelp("space/x", "select"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type model struct {
	// SDK client
	client *friendly.Client

	// Auth
	auth *friendly.Authorization

	// Current view and mode
	currentView view
	currentMode mode

	// UI components
	help      help.Model
	inputs    []textinput.Model
	focusIdx  int
	list      list.Model

	// Data
	feedEntries    []friendly.FeedEntry
	feedSelection  int
	networkFriends []friendly.UserDetails
	profileData    map[string]string
	generatedToken string

	// State
	message      string
	messageStyle lipgloss.Style
	width        int
	height       int
}

func initialModel(client *friendly.Client) model {
	// Initialize text inputs for registration form
	inputs := make([]textinput.Model, 3)

	inputs[0] = textinput.New()
	inputs[0].Placeholder = "Enter nickname (max 256 chars)"
	inputs[0].CharLimit = 256
	inputs[0].Width = 50

	inputs[1] = textinput.New()
	inputs[1].Placeholder = "Enter description (max 1024 chars)"
	inputs[1].CharLimit = 1024
	inputs[1].Width = 50

	inputs[2] = textinput.New()
	inputs[2].Placeholder = "Enter interests, comma-separated (max 64 chars each)"
	inputs[2].CharLimit = 256
	inputs[2].Width = 50

	// Initialize list for menu
	items := []list.Item{
		menuItem{title: "Feed", desc: "Browse suggested connections"},
		menuItem{title: "Network", desc: "View your friends"},
		menuItem{title: "Profile", desc: "View your profile"},
		menuItem{title: "Generate Token", desc: "Generate friend token"},
		menuItem{title: "Add Friend", desc: "Add friend by token"},
		menuItem{title: "Quit", desc: "Exit the application"},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Friendly - Main Menu"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	return model{
		client:       client,
		currentView:  viewRegister,
		currentMode:  modeInsert, // Start in insert mode for registration
		help:         help.New(),
		inputs:       inputs,
		list:         l,
		profileData:  make(map[string]string),
		messageStyle: successStyle,
	}
}

type menuItem struct {
	title, desc string
}

func (i menuItem) Title() string       { return i.title }
func (i menuItem) Description() string { return i.desc }
func (i menuItem) FilterValue() string { return i.title }

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 10)
		return m, nil

	case tea.KeyMsg:
		// Handle mode switching first
		if m.currentMode == modeInsert {
			return m.handleInsertMode(msg)
		}
		return m.handleNormalMode(msg)
	}

	return m, nil
}

func (m model) handleInsertMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Exit insert mode
		m.currentMode = modeNormal
		// Blur all inputs
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		return m, nil

	case "enter":
		// Submit or move to next field
		if m.currentView == viewRegister {
			return m.handleRegistrationSubmit()
		}
		return m, nil

	case "tab", "shift+tab":
		// Navigate between fields in insert mode
		if msg.String() == "tab" {
			m.focusIdx++
		} else {
			m.focusIdx--
		}

		if m.focusIdx > len(m.inputs)-1 {
			m.focusIdx = 0
		} else if m.focusIdx < 0 {
			m.focusIdx = len(m.inputs) - 1
		}

		m.updateFocus()
		return m, nil
	}

	// Update text inputs
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m model) handleNormalMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Quit):
		if m.currentView == viewMenu {
			return m, tea.Quit
		}
		m.currentView = viewMenu
		m.message = ""
		return m, nil

	case key.Matches(msg, keys.Back):
		if m.currentView == viewRegister && m.auth == nil {
			// Can't go back from registration if not logged in
			return m, nil
		}
		m.currentView = viewMenu
		m.message = ""
		return m, nil

	case key.Matches(msg, keys.Help):
		m.help.ShowAll = !m.help.ShowAll
		return m, nil

	case key.Matches(msg, keys.Insert):
		// Enter insert mode
		m.currentMode = modeInsert
		if m.currentView == viewRegister {
			m.inputs[m.focusIdx].Focus()
		}
		return m, nil

	case key.Matches(msg, keys.Enter):
		return m.handleEnter()

	case key.Matches(msg, keys.NextView):
		return m.nextView()

	case key.Matches(msg, keys.PrevView):
		return m.prevView()

	case key.Matches(msg, keys.Refresh):
		return m.handleRefresh()
	}

	// Handle view-specific navigation in normal mode
	switch m.currentView {
	case viewMenu:
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd

	case viewFeed:
		return m.handleFeedNavigation(msg)

	case viewNetwork:
		// Future: navigation through network
		return m, nil

	case viewRegister:
		// In normal mode on register screen, allow navigation
		if key.Matches(msg, keys.Up) || key.Matches(msg, keys.Down) {
			if key.Matches(msg, keys.Down) {
				m.focusIdx++
			} else {
				m.focusIdx--
			}

			if m.focusIdx > len(m.inputs)-1 {
				m.focusIdx = 0
			} else if m.focusIdx < 0 {
				m.focusIdx = len(m.inputs) - 1
			}
		}
		return m, nil
	}

	return m, nil
}

func (m model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.currentView {
	case viewMenu:
		selected := m.list.SelectedItem().(menuItem)
		switch selected.title {
		case "Feed":
			m.currentView = viewFeed
			m.loadFeed()
		case "Network":
			m.currentView = viewNetwork
			m.loadNetwork()
		case "Profile":
			m.currentView = viewProfile
			m.loadProfile()
		case "Generate Token":
			m.currentView = viewGenerateToken
			m.generateFriendToken()
		case "Add Friend":
			m.currentView = viewAddFriend
		case "Quit":
			return m, tea.Quit
		}

	case viewRegister:
		// In normal mode, enter switches to insert mode
		m.currentMode = modeInsert
		m.inputs[m.focusIdx].Focus()

	case viewFeed:
		return m.handleFeedSelect()
	}

	return m, nil
}

func (m model) handleRegistrationSubmit() (tea.Model, tea.Cmd) {
	// Validate all fields are filled
	if m.inputs[0].Value() == "" || m.inputs[1].Value() == "" {
		m.message = "Nickname and description are required"
		m.messageStyle = errorStyle
		return m, nil
	}

	// Create registration
	nickname, err := friendly.NewNickname(m.inputs[0].Value())
	if err != nil {
		m.message = fmt.Sprintf("Invalid nickname: %v", err)
		m.messageStyle = errorStyle
		return m, nil
	}

	description, err := friendly.NewUserDescription(m.inputs[1].Value())
	if err != nil {
		m.message = fmt.Sprintf("Invalid description: %v", err)
		m.messageStyle = errorStyle
		return m, nil
	}

	// Parse interests
	interests := []friendly.Interest{}
	if m.inputs[2].Value() != "" {
		interestStrs := strings.Split(m.inputs[2].Value(), ",")
		for _, s := range interestStrs {
			s = strings.TrimSpace(s)
			if s != "" {
				interest, err := friendly.NewInterest(s)
				if err != nil {
					m.message = fmt.Sprintf("Invalid interest '%s': %v", s, err)
					m.messageStyle = errorStyle
					return m, nil
				}
				interests = append(interests, interest)
			}
		}
	}

	// Call API
	auth, err := m.client.Generate(nickname, description, interests, nil)
	if err != nil {
		// Show full error message for debugging
		m.message = fmt.Sprintf("Registration failed: %v", err)
		m.messageStyle = errorStyle
		return m, nil
	}

	m.auth = auth
	m.message = fmt.Sprintf("Registration successful! Welcome, %s!", nickname)
	m.messageStyle = successStyle
	m.currentView = viewMenu
	m.currentMode = modeNormal

	// Clear inputs
	for i := range m.inputs {
		m.inputs[i].SetValue("")
		m.inputs[i].Blur()
	}

	return m, nil
}

func (m *model) updateFocus() {
	for i := range m.inputs {
		if i == m.focusIdx {
			m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
}

func (m model) handleFeedNavigation(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keys.Down):
		if m.feedSelection < len(m.feedEntries)-1 {
			m.feedSelection++
		}
	case key.Matches(msg, keys.Up):
		if m.feedSelection > 0 {
			m.feedSelection--
		}
	}
	return m, nil
}

func (m model) handleFeedSelect() (tea.Model, tea.Cmd) {
	if len(m.feedEntries) == 0 {
		return m, nil
	}

	entry := m.feedEntries[m.feedSelection]
	err := m.client.SendFriendRequest(m.auth, entry.Details.Id, entry.Details.AccessHash)
	if err != nil {
		m.message = fmt.Sprintf("Failed to send friend request: %v", err)
		m.messageStyle = errorStyle
	} else {
		m.message = fmt.Sprintf("Friend request sent to %s!", entry.Details.Nickname)
		m.messageStyle = successStyle
	}
	return m, nil
}

func (m model) handleRefresh() (tea.Model, tea.Cmd) {
	switch m.currentView {
	case viewFeed:
		m.loadFeed()
		m.message = "Feed refreshed"
		m.messageStyle = successStyle
	case viewNetwork:
		m.loadNetwork()
		m.message = "Network refreshed"
		m.messageStyle = successStyle
	case viewProfile:
		m.loadProfile()
		m.message = "Profile refreshed"
		m.messageStyle = successStyle
	}
	return m, nil
}

func (m *model) loadFeed() {
	if m.auth == nil {
		m.feedEntries = nil
		return
	}

	feed, err := m.client.GetFeedQueue(m.auth)
	if err != nil {
		m.message = fmt.Sprintf("Error loading feed: %v", err)
		m.messageStyle = errorStyle
		m.feedEntries = nil
		return
	}

	m.feedEntries = feed.Entries
	m.feedSelection = 0
}

func (m *model) loadNetwork() {
	if m.auth == nil {
		m.networkFriends = nil
		return
	}

	network, err := m.client.GetNetworkDetails(m.auth)
	if err != nil {
		m.message = fmt.Sprintf("Error loading network: %v", err)
		m.messageStyle = errorStyle
		m.networkFriends = nil
		return
	}

	m.networkFriends = network.Friends
}

func (m *model) loadProfile() {
	if m.auth == nil {
		m.profileData = map[string]string{"Status": "Not logged in"}
		return
	}

	details, err := m.client.GetSelfDetails(m.auth)
	if err != nil {
		m.profileData = map[string]string{"Error": fmt.Sprintf("%v", err)}
		return
	}

	network, _ := m.client.GetNetworkDetails(m.auth)
	friendCount := "0"
	if network != nil {
		friendCount = fmt.Sprintf("%d", len(network.Friends))
	}

	interests := convertInterestsToStrings(details.Interests)

	m.profileData = map[string]string{
		"Nickname":    string(details.Nickname),
		"Description": string(details.Description),
		"Interests":   strings.Join(interests, ", "),
		"Friends":     friendCount,
		"User ID":     fmt.Sprintf("%d", details.Id),
	}
}

func (m *model) generateFriendToken() {
	if m.auth == nil {
		m.message = "Please register first"
		m.messageStyle = errorStyle
		return
	}

	token, err := m.client.GenerateFriendToken(m.auth)
	if err != nil {
		m.message = fmt.Sprintf("Failed to generate token: %v", err)
		m.messageStyle = errorStyle
		return
	}

	m.generatedToken = string(token)
	m.message = "Token generated successfully!"
	m.messageStyle = successStyle
}

func (m *model) nextView() (tea.Model, tea.Cmd) {
	if m.auth == nil {
		return *m, nil
	}
	views := []view{viewMenu, viewFeed, viewNetwork, viewProfile}
	for i, v := range views {
		if v == m.currentView {
			m.currentView = views[(i+1)%len(views)]
			break
		}
	}
	return *m, nil
}

func (m *model) prevView() (tea.Model, tea.Cmd) {
	if m.auth == nil {
		return *m, nil
	}
	views := []view{viewMenu, viewFeed, viewNetwork, viewProfile}
	for i, v := range views {
		if v == m.currentView {
			m.currentView = views[(i-1+len(views))%len(views)]
			break
		}
	}
	return *m, nil
}

func convertInterestsToStrings(interests []friendly.Interest) []string {
	strs := make([]string, len(interests))
	for i, interest := range interests {
		strs[i] = string(interest)
	}
	return strs
}

func (m model) View() string {
	var s strings.Builder

	// Header
	modeStr := "NORMAL"
	modeStyle := focusedStyle
	if m.currentMode == modeInsert {
		modeStr = "INSERT"
		modeStyle = successStyle
	}

	s.WriteString(titleStyle.Render("🤝 Friendly CLI"))
	s.WriteString("  ")
	s.WriteString(modeStyle.Render(fmt.Sprintf("[%s]", modeStr)))
	s.WriteString("\n\n")

	// View-specific content
	switch m.currentView {
	case viewRegister:
		s.WriteString(focusedStyle.Render("📝 Register New Account"))
		s.WriteString("\n\n")

		for i, input := range m.inputs {
			s.WriteString(input.View())
			if i == m.focusIdx && m.currentMode == modeNormal {
				s.WriteString(focusedStyle.Render(" ◀"))
			}
			s.WriteString("\n\n")
		}

		if m.currentMode == modeInsert {
			s.WriteString(helpStyle.Render("Tab: next field • Shift+Tab: prev field • Enter: submit • Esc: normal mode"))
		} else {
			s.WriteString(helpStyle.Render("j/k: navigate fields • i/a: insert mode • Enter: edit field"))
		}

	case viewMenu:
		s.WriteString(m.list.View())

	case viewFeed:
		s.WriteString(focusedStyle.Render("📱 Feed - Suggested Connections"))
		s.WriteString("\n\n")

		if len(m.feedEntries) == 0 {
			s.WriteString(blurredStyle.Render("No suggestions available. Press 'r' to refresh."))
		} else {
			for i, entry := range m.feedEntries {
				commonStr := ""
				if len(entry.CommonFriends) > 0 {
					commonStr = fmt.Sprintf("\n   💫 %d common friends", len(entry.CommonFriends))
				}
				interests := strings.Join(convertInterestsToStrings(entry.Details.Interests), ", ")
				entryStr := fmt.Sprintf("👤 %s\n   %s\n   🏷️  %s%s",
					entry.Details.Nickname,
					entry.Details.Description,
					interests,
					commonStr)

				if i == m.feedSelection {
					s.WriteString(selectedBoxStyle.Render(entryStr))
				} else {
					s.WriteString(boxStyle.Render(entryStr))
				}
				s.WriteString("\n")
			}
		}

		s.WriteString("\n")
		s.WriteString(helpStyle.Render("j/k: navigate • Enter: send friend request • r: refresh • Esc: back"))

	case viewNetwork:
		s.WriteString(focusedStyle.Render("👥 Your Network"))
		s.WriteString("\n\n")

		if len(m.networkFriends) == 0 {
			s.WriteString(blurredStyle.Render("No friends yet. Press 'r' to refresh."))
		} else {
			for _, friend := range m.networkFriends {
				interests := strings.Join(convertInterestsToStrings(friend.Interests), ", ")
				friendStr := fmt.Sprintf("👥 %s\n   %s\n   🏷️  %s",
					friend.Nickname,
					friend.Description,
					interests)
				s.WriteString(boxStyle.Render(friendStr))
				s.WriteString("\n")
			}
		}

		s.WriteString("\n")
		s.WriteString(helpStyle.Render("r: refresh • Esc: back"))

	case viewProfile:
		s.WriteString(focusedStyle.Render("👤 Your Profile"))
		s.WriteString("\n\n")

		for k, v := range m.profileData {
			s.WriteString(focusedStyle.Render(k + ": "))
			s.WriteString(v)
			s.WriteString("\n")
		}

		s.WriteString("\n")
		s.WriteString(helpStyle.Render("r: refresh • Esc: back"))

	case viewGenerateToken:
		s.WriteString(focusedStyle.Render("🔑 Friend Token"))
		s.WriteString("\n\n")

		if m.generatedToken != "" {
			s.WriteString("Your friend token (share with friends):\n\n")
			s.WriteString(boxStyle.Render(m.generatedToken))
			s.WriteString("\n\n")
			s.WriteString(blurredStyle.Render("Copy this token and share it with someone you want to add as a friend."))
		} else {
			s.WriteString(blurredStyle.Render("Token not generated yet."))
		}

		s.WriteString("\n\n")
		s.WriteString(helpStyle.Render("Esc: back"))

	case viewAddFriend:
		s.WriteString(focusedStyle.Render("➕ Add Friend"))
		s.WriteString("\n\n")
		s.WriteString(blurredStyle.Render("Enter friend token to connect (TODO: implement)"))
		s.WriteString("\n\n")
		s.WriteString(helpStyle.Render("Esc: back"))
	}

	// Message
	if m.message != "" {
		s.WriteString("\n\n")
		s.WriteString(m.messageStyle.Render(m.message))
	}

	// Help
	if m.currentMode == modeNormal {
		s.WriteString("\n\n")
		s.WriteString(m.help.View(keys))
	}

	return s.String()
}

func main() {
	client := friendly.NewMeetacyClient()

	p := tea.NewProgram(initialModel(client), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
