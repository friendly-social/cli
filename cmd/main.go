package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/friendly-social/cli/internal/app"
	"github.com/friendly-social/cli/internal/auth"
	"github.com/friendly-social/cli/internal/vim"
)

func main() {
	screens := []app.Screen{
		auth.NewScreen(),
	}

	router := app.NewRouter(screens)
	wrapper := vim.NewWrapper(router)

	p := tea.NewProgram(wrapper)
	if _, err := p.Run(); err != nil {
		panic("failed to run app router: " + err.Error())
	}
}
