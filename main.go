package main

import (
	"budget-cli/db"
	"budget-cli/models"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Printf("Fatal error setting up logging: %v", err)
		os.Exit(1)
	}
	defer f.Close()

	if dbLoadErr := db.InitDB(); dbLoadErr != nil {
		fmt.Printf("Error occured: %v", err)
	}

	p := tea.NewProgram(models.NewRootModel(), tea.WithAltScreen()) //needs an initial model
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error occurred while running budget-cli")
		os.Exit(1)
	}
}
