package progress

import (
	"fmt"
	"os"
	"strings"
	"time"

	tprogress "github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding  = 2
	maxWidth = 80
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

func ProgressBar(progress func() float64) {
	prog := tprogress.New(tprogress.WithScaledGradient("#FF7CCB", "#FDFF8C"))

	if _, err := tea.NewProgram(model{progress: prog, getCP: progress}).Run(); err != nil {
		fmt.Println("Error with tea progress bar:", err)
		os.Exit(1)
	}
}

type tickMsg time.Time

type model struct {
	getCP    func() float64
	percent  float64
	progress tprogress.Model
}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.percent = m.getCP()
		return m, tick()

	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		m.progress.Width = min(maxWidth, m.progress.Width)
		return m, nil

	default:
		m.percent = m.getCP()
		return m, nil
	}

	return m, nil
}

func (m model) View() string {
	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + m.progress.ViewAs(m.percent) + "\n\n" +
		pad + helpStyle("Press q to quit")
}

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
