package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type status int

const divisor = 4

const (
	todo status = iota
	inProgress
	done
)

// Styling
var (
	columnStyle = lipgloss.NewStyle().
			Padding(1, 2)
	focusedStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

// Custom item
type Task struct {
	status      status
	title       string
	description string
}

// Implement the list.Item interface
func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}

// * MAIN MODEL *
type Model struct {
	focused  status
	lists    []list.Model
	err      error
	loaded   bool
	quitting bool
}

func New() *Model {
	return &Model{}
}

// TODO: call this on tea.WindowSizeMsg
func (m *Model) initLists(width, height int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width/divisor, height-divisor)
	defaultList.SetShowHelp(false)
	m.lists = []list.Model{defaultList, defaultList, defaultList}

	// Init todo
	m.lists[todo].Title = "To do"
	m.lists[todo].SetItems([]list.Item{
		Task{status: todo, title: "buy milk", description: "whole milk"},
		Task{status: todo, title: "fold laundry", description: "flat and tidy"},
		Task{status: todo, title: "play with cars", description: "zoom zoom"},
	})

	// Init in progress
	m.lists[inProgress].Title = "In progress"
	m.lists[inProgress].SetItems([]list.Item{
		Task{status: todo, title: "get good at computer", description: "gonna take a while"},
	})

	// Init done
	m.lists[done].Title = "Done"
	m.lists[done].SetItems([]list.Item{
		Task{status: todo, title: "dishes", description: "yuck!"},
	})
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.loaded {
			m.initLists(msg.Width, msg.Height)
			m.loaded = true
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "q", tea.KeyCtrlC.String():
			m.quitting = true
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}
	if m.loaded {
		todoView := m.lists[todo].View()
		inProgView := m.lists[inProgress].View()
		doneView := m.lists[done].View()

		switch m.focused {
		default:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				focusedStyle.Render(todoView),
				columnStyle.Render(inProgView),
				columnStyle.Render(doneView),
			)
		}
	}
	return "Loading..."
}

// do stuff
func main() {
	m := New()
	p := tea.NewProgram(m)
	_, err := p.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
