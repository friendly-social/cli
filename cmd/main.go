package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/friendly-social/cli/internal/navigation"
	"github.com/friendly-social/cli/internal/router"
	"github.com/friendly-social/cli/internal/screen"
	"github.com/friendly-social/cli/internal/screen/auth"
	"github.com/friendly-social/cli/internal/screen/home"
)

func main() {
	screens := []screen.Model{
		auth.New(),
		home.New(),
	}

	router := router.NewRouter(screens)
	wrapper := navigation.NewVimWrapper(router)

	p := tea.NewProgram(wrapper)
	if _, err := p.Run(); err != nil {
		panic("failed to run app router: " + err.Error())
	}
}
