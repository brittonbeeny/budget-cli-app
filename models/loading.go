package models

import (
	"budget-cli/styles"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type doneLoading bool

type LoadingModel struct {
	window  *WindowSize
	spinner spinner.Model
	loading bool
}

var cmds []tea.Cmd

func NewLoadingModel(windowSize *WindowSize) LoadingModel {

	s := spinner.New()
	s.Spinner = spinner.Meter
	s.Style = styles.BaseStyle

	return LoadingModel{
		window:  windowSize,
		spinner: s,
		loading: true,
	}
}

func loadingCmd() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(1 * time.Second)
		return doneLoading(true)
	}
}

func tickCmd(spinner spinner.Model) tea.Cmd {
	return func() tea.Msg {
		return tea.Msg(spinner.Tick())
	}
}

func (m LoadingModel) WindowSize() *WindowSize {
	return m.window
}

func (m LoadingModel) Init() tea.Cmd {
	cmds = append(cmds, tickCmd(m.spinner), loadingCmd())
	return tea.Batch(cmds...)
}

func (m LoadingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case doneLoading:
		m.loading = false
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m LoadingModel) View() string {
	if m.loading {
		label := "Loading, please wait..."
		str := fmt.Sprintf("\n   %s %s\n\n", m.spinner.View(), label)

		spinnerStyle := m.spinner.Style.Width(m.window.Width).Height(m.window.Height)
		return spinnerStyle.Render(str)
	}

	return ""
}
