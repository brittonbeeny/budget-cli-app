package annual_overview

import (
	"budget-cli/shared"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestAnnualOverviewLeftPaneCursorUpdate(t *testing.T) {

	annualOverview := initiateModel()

	downPress := tea.KeyMsg{Type: tea.KeyDown}

	updatedModel, _ := sendTeaMsg(annualOverview, downPress)
	annualOveriew, ok := updatedModel.(AnnualOverview)

	if !ok {
		t.Fatal("Returned model is not an AnnualOverview")
	}

	currentLeftPaneCursorVal := annualOveriew.leftPane.cursor

	if currentLeftPaneCursorVal != 1 {
		t.Fatalf("Left Pane counter is not correct: %d", currentLeftPaneCursorVal)
	}
}

func TestAnnualOverviewLeftPaneMultiCursorUpdate(t *testing.T) {

	annualOverview := initiateModel().(AnnualOverview)
	downPress := tea.KeyMsg{Type: tea.KeyDown}
	upPress := tea.KeyMsg{Type: tea.KeyUp}

	_, _ = sendTeaMsg(annualOverview, downPress)
	updatedModel, _ := sendTeaMsg(annualOverview, upPress)

	annualOverview, ok := updatedModel.(AnnualOverview)

	if !ok {
		t.Fatal("Updated model is not of type AnnualOverview")
	}

	currentLeftPaneCursorVal := annualOverview.leftPane.cursor

	if currentLeftPaneCursorVal != 0 {
		t.Fatalf("Left Pane counter is not correct: %d", currentLeftPaneCursorVal)
	}
}

func initiateModel() tea.Model {
	annualOverview := NewAnnualOverview(&shared.WindowSize{Width: 80, Height: 80})
	annualOverview.leftPane.years = []int{2026, 2027}

	return annualOverview
}

func sendTeaMsg(model tea.Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	updatedModel, cmd := model.Update(msg)
	return updatedModel, cmd
}
