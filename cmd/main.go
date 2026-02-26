package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/friendly-social/cli/internal/auth"
	"github.com/friendly-social/cli/internal/feed"
	"github.com/friendly-social/cli/internal/navigation"
	"github.com/friendly-social/cli/internal/router"
	"github.com/friendly-social/cli/internal/screen"
)

func main() {
	screens := []screen.Model{
		auth.NewScreen(),
		feed.NewScreen(),
	}

	router := router.NewRouter(screens)
	wrapper := navigation.NewVimWrapper(router)

	p := tea.NewProgram(wrapper)
	if _, err := p.Run(); err != nil {
		panic("failed to run app router: " + err.Error())
	}
}
