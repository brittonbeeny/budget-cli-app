package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	BaseColor             = lipgloss.Color("40")
	BaseStyle             = lipgloss.NewStyle().Foreground(BaseColor).Align(lipgloss.Center, lipgloss.Center)
	activeChoicePadding   = 1
	inactiveChoicePadding = activeChoicePadding + 1
	ActiveChoiceStyle     = lipgloss.NewStyle().
				Background(lipgloss.Color(BaseColor)).
				Foreground(lipgloss.Color("0")).
				Border(lipgloss.BlockBorder()).
				BorderForeground(lipgloss.Color("0")).
				Padding(activeChoicePadding)
	InactiveChoiceStyle = lipgloss.NewStyle().
		// PaddingTop and PaddingBottom compensate for the missing top/bottom borders
		Padding(inactiveChoicePadding)
)
