package tests

import (
	"budget-cli/models"
	"strings"
	"testing"
)

func TestMenuViewSelectionDoesNotChangeLayoutHeight(t *testing.T) {
	menu := models.NewMenuModel(&models.WindowSize{Width: 80, Height: 24})

	menu.SetCursor(0)
	selected := menu.View()

	menu.SetCursor(1)
	unselected := menu.View()

	if strings.Count(selected, "\n") != strings.Count(unselected, "\n") {
		t.Fatalf("expected selection to keep the same layout height, got %d lines vs %d lines", strings.Count(selected, "\n"), strings.Count(unselected, "\n"))
	}
}
