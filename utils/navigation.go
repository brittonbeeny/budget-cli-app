package utils

import (
	tea "github.com/charmbracelet/bubbletea"
)

type BackToHome bool
type BackToAnnual bool

func GoBackHomeCmd() tea.Cmd {
	return func() tea.Msg {
		return BackToHome(true)
	}
}

func GoBackToAnnualCmd() tea.Cmd {
	return func() tea.Msg {
		return BackToAnnual(true)
	}
}
