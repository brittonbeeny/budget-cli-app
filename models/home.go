package models

import (
	"budget-cli/models/annual_overview"
	"budget-cli/shared"
	"budget-cli/styles"
	"budget-cli/utils"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const bannerContent = `
    ____  __  ______  __________________    ________    ____
   / __ )/ / / / __ \/ ____/ ____/_  __/   / ____/ /   /  _/
  / __  / / / / / / / / __/ __/   / ______/ /   / /    / /  
 / /_/ / /_/ / /_/ / /_/ / /___  / /_____/ /___/ /____/ /   
/_____/\____/_____/\____/_____/ /_/      \____/_____/___/   
                                                            
`

type HomeModel struct {
	window       *shared.WindowSize
	Banner       Banner
	menu         MenuModel
	overview     annual_overview.AnnualOverview
	gotoMenu     bool
	gotoOverview bool
}

type Banner struct {
	content string
	style   lipgloss.Style
}

func NewHomeModel(windowSize *shared.WindowSize) HomeModel {
	banner := Banner{
		content: bannerContent,
		style:   styles.BaseStyle,
	}

	return HomeModel{
		window:   windowSize,
		Banner:   banner,
		overview: annual_overview.NewAnnualOverview(windowSize),
	}
}

func (m HomeModel) WindowSize() *shared.WindowSize {
	return m.window
}

func (m HomeModel) MenuModel() MenuModel {
	return m.menu
}

func (m HomeModel) Init() tea.Cmd {
	return nil
}

func (m HomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			log.Println("Home: pressed Enter")
			if !m.gotoOverview {
				log.Println("Home: pressed Enter to transition to Menu")
				m.gotoOverview = true
				return m, nil
			}
		}
	}
	if _, ok := msg.(utils.BackToHome); ok {
		m.gotoOverview = false
		return m, nil
	} else if m.gotoOverview { //Pass down to Menu
		overviewModel, overviewCmd := m.overview.Update(msg)
		m.overview = overviewModel.(annual_overview.AnnualOverview)
		return m, overviewCmd
	}
	return m, nil
}

func (m HomeModel) View() string {
	if m.gotoOverview {
		return m.overview.View()
	}

	content := "\n\nPress Enter to continue\n\n"

	bannerStyle := m.Banner.style
	if m.window != nil {
		bannerStyle = bannerStyle.Width(m.window.Width).Height(m.window.Height)
	}

	return bannerStyle.Render(m.Banner.content, content)
}
