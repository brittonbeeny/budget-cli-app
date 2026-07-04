package tests

import (
	"budget-cli/models"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestRootModelUpdate_WindowSizeMsgUpdatesSharedWindow(t *testing.T) {
	m := models.NewRootModel()

	updatedModel, _ := m.Update(tea.WindowSizeMsg{Width: 123, Height: 456})
	root, ok := updatedModel.(models.RootModel)
	if !ok {
		t.Fatalf("expected RootModel, got %T", updatedModel)
	}

	if root.WindowSize() == nil {
		t.Fatal("expected root window pointer to be initialized")
	}

	if root.WindowSize() != root.LoadingModel().WindowSize() || root.WindowSize() != root.HomeModel().WindowSize() || root.WindowSize() != root.HomeModel().MenuModel().WindowSize() {
		t.Fatal("expected loading, home, and menu models to share the same window pointer")
	}

	if root.WindowSize().Width != 123 || root.WindowSize().Height != 456 {
		t.Fatalf("expected shared window size to be updated, got width=%d height=%d", root.WindowSize().Width, root.WindowSize().Height)
	}
}
