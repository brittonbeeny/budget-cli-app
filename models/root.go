package models

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

type sessionState int

const (
	loadingView sessionState = iota
	homeView
	quit
)

type WindowSize struct {
	Height int
	Width  int
}

type RootModel struct {
	state   sessionState
	window  *WindowSize
	loading LoadingModel
	home    HomeModel
}

func NewRootModel() RootModel {

	sharedWindowSize := &WindowSize{Width: 0, Height: 0}

	return RootModel{
		state:   loadingView,
		window:  sharedWindowSize,
		loading: NewLoadingModel(sharedWindowSize),
		home:    NewHomeModel(sharedWindowSize),
	}
}

func (m RootModel) WindowSize() *WindowSize {
	return m.window
}

func (m RootModel) LoadingModel() LoadingModel {
	return m.loading
}

func (m RootModel) HomeModel() HomeModel {
	return m.home
}

func (m RootModel) Init() tea.Cmd {
	return m.loading.Init()
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q": //Global Quit
			log.Println("Quit detected")
			m.state = quit
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.window.Width = msg.Width
		m.window.Height = msg.Height

	}

	switch m.state {
	case loadingView:
		if _, ok := msg.(doneLoading); ok {
			m.state = homeView
			return m, nil
		}

		newLoading, cmd := m.loading.Update(msg)
		m.loading = newLoading.(LoadingModel)
		return m, cmd

	case homeView:
		var cmd tea.Cmd
		home, cmd := m.home.Update(msg)
		m.home = home.(HomeModel)
		return m, cmd
	}

	return m, nil
}

func (m RootModel) View() string {
	switch m.state {
	case loadingView:
		return m.loading.View()
	case homeView:
		return m.home.View()
	default:
		return ""
	}
}
