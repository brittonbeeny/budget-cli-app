package models

import (
	"budget-cli/styles"
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type menuState int

const (
	menuView menuState = iota
	createView
)

type MenuModel struct {
	state          menuState
	window         *WindowSize
	cursorPosition int
	choices        []choice
	style          lipgloss.Style
	createBudget   CreateBudgetModel
	menuShown      bool
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
		state:        menuView,
		window:       windowSize,
		choices:      menuChoices,
		style:        styles.BaseStyle,
		createBudget: NewCreateBudgetModel(),
		menuShown:    false,
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
	log.Printf("Menu: Incoming msg type=%T, Current State=%v", msg, m.state)

	if m.state == createView {
		log.Println("Menu: Routing directly to CreateBudget")
		var createCmd tea.Cmd
		cbm, createCmd := m.createBudget.Update(msg)
		m.createBudget = cbm.(CreateBudgetModel)
		return m, createCmd
	}

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

		case "enter":
			log.Println("Menu: Enter pressed")

			switch m.cursorPosition {
			case 0:
				log.Println("Menu: Selected first option (Create Budget)")
				m.state = createView
				return m, nil

			case 1:
				log.Println("Menu: Selected second option")
			}
		}
	}

	return m, nil
}

func (m MenuModel) View() string {
	switch m.state {
	case createView:
		return m.createBudget.View()
	default:
		return m.getMenuView()
	}

}

func (m *MenuModel) getMenuView() string {

	m.menuShown = true
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
