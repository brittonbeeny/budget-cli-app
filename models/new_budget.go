package models

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

type CreateBudgetModel struct {
}

func NewCreateBudgetModel() CreateBudgetModel {
	return CreateBudgetModel{}
}

func (m CreateBudgetModel) Init() tea.Cmd {
	log.Println("Init called in CreateBudgetModel")
	return nil
}

func (m CreateBudgetModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m CreateBudgetModel) View() string {
	return "This is the add budget page"
}
