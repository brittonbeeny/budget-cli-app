package models

import (
	"budget-cli/styles"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MenuModel struct {
	window         *WindowSize
	cursorPosition int
	choices        []choice
	style          lipgloss.Style
}

type choice struct {
	text string
}

type backToHome bool

func NewMenuModel(windowSize *WindowSize) MenuModel {

	menuChoices := []choice{
		{"Create A New Budget"},
		{"View Budgets"},
	}
	return MenuModel{
		window:  windowSize,
		choices: menuChoices,
		style:   styles.BaseStyle,
	}
}

func (m MenuModel) WindowSize() *WindowSize {
	return m.window
}

func (m *MenuModel) SetCursor(cursor int) {
	m.cursorPosition = cursor
}

func (m MenuModel) Init() tea.Cmd {
	return nil
}

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursorPosition > 0 {
				m.cursorPosition--
			}
		case "down", "j":
			if m.cursorPosition < len(m.choices)-1 {
				m.cursorPosition++
			}
		case "esc":
			m.cursorPosition = 0
			return m, backCmd()
		}

	}
	return m, nil
}

func (m MenuModel) View() string {
	content := "What would you like to do?\n\n"

	lines := []string{content}
	for i, choice := range m.choices {
		line := ""
		if m.cursorPosition == i {
			line = styles.ActiveChoiceStyle.Render(choice.text)
		} else {
			line = styles.InactiveChoiceStyle.Render(choice.text)
		}
		lines = append(lines, line)
	}

	lines = append(lines, fmt.Sprintln("Press ENTER to continue, Press ESC to go back"))

	return m.style.
		Width(m.window.Width).
		Height(m.window.Height).
		Render(lipgloss.JoinVertical(lipgloss.Center, lines...))
}

func backCmd() tea.Cmd {
	return func() tea.Msg {
		return backToHome(true)
	}
}
