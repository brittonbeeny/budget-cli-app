package main

import (
	"budget-cli/models"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(models.NewRootModel(), tea.WithAltScreen()) //needs an initial model
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error occurred while running budget-cli")
		os.Exit(1)
	}
}
