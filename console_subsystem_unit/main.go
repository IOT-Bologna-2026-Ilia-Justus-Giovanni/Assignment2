package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type tickMsg struct{}

type Model struct {
	State   string
	Sectors [6]string
}

func NewModel() Model {
	model := Model{State: "IDLE"}
	for i := range model.Sectors {
		model.Sectors[i] = "empty"
	}
	return model
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m Model) Init() tea.Cmd {
	return tick()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tickMsg:
		if m.Sectors[0] == "empty" {
			m.Sectors[0] = "object"
		} else {
			m.Sectors[0] = "empty"
		}
		return m, tick()
	}

	return m, nil
}

func (m Model) View() string {
	var view strings.Builder
	fmt.Fprintln(&view, "Radar Console")
	fmt.Fprintf(&view, "State: %s\n\n", m.State)
	fmt.Fprintln(&view, "Sectors:")
	for i, sector := range m.Sectors {
		fmt.Fprintf(&view, "  %d: %s\n", i+1, sector)
	}
	fmt.Fprintln(&view, "\nPress q to quit.")
	return view.String()
}

func main() {
	program := tea.NewProgram(NewModel(), tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Console error:", err)
		os.Exit(1)
	}
}
