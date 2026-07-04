package models

import (
	"budget-cli/styles"

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
	window   *WindowSize
	Banner   Banner
	menu     MenuModel
	gotoMenu bool
}

type Banner struct {
	content string
	style   lipgloss.Style
}

func NewHomeModel(windowSize *WindowSize) HomeModel {
	banner := Banner{
		content: bannerContent,
		style:   styles.BaseStyle,
	}

	return HomeModel{
		window: windowSize,
		Banner: banner,
		menu:   NewMenuModel(windowSize),
	}
}

func (m HomeModel) WindowSize() *WindowSize {
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
			m.gotoMenu = true
			menuModel, menuCmd := m.menu.Update(msg)
			m.menu = menuModel.(MenuModel)
			return m, menuCmd
		}
	}
	if _, ok := msg.(backToHome); ok {
		m.gotoMenu = false
		return m, nil
	} else if m.gotoMenu {
		menuModel, menuCmd := m.menu.Update(msg)
		m.menu = menuModel.(MenuModel)
		return m, menuCmd
	}
	return m, nil
}

func (m HomeModel) View() string {
	if m.gotoMenu {
		return m.menu.View()
	}

	content := "\n\nPress Enter to continue\n\nPress CTRL+C, q to Quit"

	bannerStyle := m.Banner.style
	if m.window != nil {
		bannerStyle = bannerStyle.Width(m.window.Width).Height(m.window.Height)
	}

	return bannerStyle.Render(m.Banner.content, content)
}
