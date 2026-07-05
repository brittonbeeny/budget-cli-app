package annual_overview

import (
	"budget-cli/shared"
	"budget-cli/styles"
	"budget-cli/utils"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AnnualOverview struct {
	leftPane  leftPane
	rightPane rightPane
	window    *shared.WindowSize
}

type leftPane struct {
	cursor int
	years  []int
	style  lipgloss.Style
}

type rightPane struct {
	totalIncome   income
	totalExpenses expenses
	top5Expenses  top5
	goals         goals
	style         lipgloss.Style
}

type income struct {
	totalIncome float64
	style       lipgloss.Style
}

type expenses struct {
	totalExpense float64
	style        lipgloss.Style
}

type top5 struct {
	top5 [5]float64
}

type goals struct {
	goals []string
}

type backToHome bool

func NewAnnualOverview(windowSize *shared.WindowSize) AnnualOverview {
	leftPane := leftPane{
		years: []int{2026, 2027, 2028}, //this will come from DB later based on existing budgets
	}

	income := income{}
	expenses := expenses{}
	top5Expenses := top5{}
	goals := goals{}

	rightPane := rightPane{
		totalIncome:   income,
		totalExpenses: expenses,
		top5Expenses:  top5Expenses,
		goals:         goals,
	}

	return AnnualOverview{
		leftPane:  leftPane,
		rightPane: rightPane,
		window:    windowSize,
	}

}

func (m AnnualOverview) Init() tea.Cmd {
	return nil
}

func (m AnnualOverview) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.leftPane.cursor > 0 {
				m.leftPane.cursor--
			}
		case "down":
			if m.leftPane.cursor < len(m.leftPane.years)-1 {
				m.leftPane.cursor++
			}
		case "esc":
			m.leftPane.cursor = 0
			return m, utils.GoBackHomeCmd()
		}
	}

	return m, nil
}

func (m AnnualOverview) View() string {

	leftWidth := (m.window.Width / 4) - 2
	rightWidth := (3 * m.window.Width / 4) - 1
	paneHeight := m.window.Height - 2

	leftPaneView := m.leftPane.getLeftPaneView(leftWidth, paneHeight)
	rightPaneView := m.rightPane.getRightPaneView(rightWidth, paneHeight)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftPaneView, rightPaneView)

}

func (lp leftPane) getLeftPaneView(width, height int) string {
	lines := []string{}
	for i, year := range lp.years {
		line := ""
		if lp.cursor == i {
			line = styles.ActiveChoiceStyle.Render(strconv.Itoa(year))
		} else {
			line = styles.InactiveChoiceStyle.Render(strconv.Itoa(year))
		}
		lines = append(lines, line)
	}

	return lp.style.
		Border(lipgloss.NormalBorder()).
		Width(width).
		Height(height).
		Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
}

func (rp rightPane) getRightPaneView(width, height int) string {
	content := "This is the overview pane"

	return rp.style.
		Border(lipgloss.NormalBorder()).
		Width(width).
		Height(height).
		Align(lipgloss.Top, lipgloss.Center).
		Render(content)
}

func backCmd() tea.Cmd {
	return func() tea.Msg {
		return backToHome(true)
	}
}
