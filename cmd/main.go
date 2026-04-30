package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/friendly-social/cli/internal/navigation"
	"github.com/friendly-social/cli/internal/router"
	"github.com/friendly-social/cli/internal/screen"
	"github.com/friendly-social/cli/internal/screen/auth"
	"github.com/friendly-social/cli/internal/screen/home"
	sdk "github.com/friendly-social/golang-sdk"
)

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close() //nolint:errcheck

	screens := []screen.Model{
		auth.New(auth.NewService(sdk.NewClient())),
		home.New(),
	}

	router := router.NewRouter(screens)
	wrapper := navigation.NewVimWrapper(router)

	p := tea.NewProgram(wrapper)
	if _, err := p.Run(); err != nil {
		panic("failed to run app router: " + err.Error())
	}
}
