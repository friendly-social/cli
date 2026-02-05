package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/friendly-social/cli/internal/root"
	sdk "github.com/friendly-social/golang-sdk"
)

func main() {
	client := sdk.NewProductionClient()
	p := tea.NewProgram(root.NewModel(client))

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
