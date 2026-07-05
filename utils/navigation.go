package utils

import (
	tea "github.com/charmbracelet/bubbletea"
)

type BackToHome bool

func GoBackHomeCmd() tea.Cmd {
	return func() tea.Msg {
		return BackToHome(true)
	}
}
