package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"go.bug.st/serial"
)

type serialMsg struct {
	value int
	err   error
}

type Model struct {
	State   string
	Sectors [6]string
	reader  *bufio.Reader
	Err     error
}

func NewModel(reader *bufio.Reader) Model {
	model := Model{State: "IDLE", reader: reader}
	for i := range model.Sectors {
		model.Sectors[i] = "empty"
	}
	return model
}

func readSignal(reader *bufio.Reader) tea.Cmd {
	return func() tea.Msg {
		line, err := reader.ReadString('\n')
		if err != nil {
			return serialMsg{err: err}
		}

		value, err := strconv.Atoi(strings.TrimSpace(line))
		if err == nil && value != 0 && value != 1 {
			err = fmt.Errorf("invalid serial signal %q", strings.TrimSpace(line))
		}
		return serialMsg{value: value, err: err}
	}
}

func (m Model) Init() tea.Cmd {
	return readSignal(m.reader)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case serialMsg:
		if msg.err != nil {
			m.Err = msg.err
			return m, tea.Quit
		}
		if msg.value == 0 {
			m.Sectors[0] = "empty"
		} else {
			m.Sectors[0] = "object"
		}
		return m, readSignal(m.reader)
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
	portName := flag.String("port", "", "serial port, for example /dev/cu.usbmodem14101")
	baudRate := flag.Int("baud", 9600, "serial baud rate")
	flag.Parse()

	if *portName == "" {
		fmt.Fprintln(os.Stderr, "Usage: go run . -port <serial-port> [-baud 9600]")
		os.Exit(2)
	}

	port, err := serial.Open(*portName, &serial.Mode{BaudRate: *baudRate})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Serial error:", err)
		os.Exit(1)
	}
	defer port.Close()

	program := tea.NewProgram(NewModel(bufio.NewReader(port)), tea.WithAltScreen())
	finalModel, err := program.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Console error:", err)
		os.Exit(1)
	}
	if model, ok := finalModel.(Model); ok && model.Err != nil {
		fmt.Fprintln(os.Stderr, "Serial error:", model.Err)
		os.Exit(1)
	}
}
