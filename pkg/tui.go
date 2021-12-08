package jdi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/davecgh/go-spew/spew"
)

const listHeight = 14

var (
	docStyle          = lipgloss.NewStyle().Background(lipgloss.Color("#000000")).Foreground(lipgloss.Color("#33FF33"))
	titleStyle        = lipgloss.NewStyle().MarginLeft(2).Bold(true).Underline(true)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item struct {
	title, desc, jobName string
}

func (j *JobSpec) GetTitle() string {
	if len(j.runs) > 0 && j.runs[0].Status != 0 {
		return j.Name + " ⛔️"
	}
	return j.Name
}

func (j *JobSpec) GetStatusDescription() string {
	spew.Dump(999, j)
	if len(j.runs) == 0 {
		return ""
	}

	lastRun := j.runs[0]
	var sb strings.Builder

	since := time.Since(lastRun.TriggeredAt).String()
	sb.WriteString("ran " + since + " ago")

	if lastRun.Status != 0 {
		sb.WriteString(" | ERROR")
	}

	return sb.String()

}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list     list.Model
	state    *Schedule
	choice   string
	quitting bool
}

func (j *JobSpec) View() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Job: %s\n", j.Name))
	// sb.WriteString(fmt.Sprintf("> ran %v times\n", len(j.Statuses)))

	// sum := 0
	// for range j.Statuses {
	// 	sum += 1
	// }

	// sb.WriteString(fmt.Sprintf("> %v%% sucessful\n\n", (float64(sum) / float64(len(j.Statuses)) * 100)))
	// sb.WriteString("## Log tail\n\n")
	// sb.WriteString(j.LogTail)

	return sb.String()

}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		// case "enter":
		// 	i, ok := m.list.SelectedItem().(item)
		// 	if ok {
		// 		m.choice = i.title
		// 	}
		// 	return m, tea.Quit
		default:
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = i.jobName
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var jobViewContent string
	if _, ok := m.state.Jobs[m.choice]; ok {
		jobViewContent = m.state.Jobs[m.choice].View()
	} else {
		jobViewContent = "Please select a job."
	}

	// job view
	vp := viewport.Model{Width: 78, Height: 20}
	renderer, _ := glamour.NewTermRenderer(glamour.WithStylePath("notty"), glamour.WithWordWrap(40))
	str, _ := renderer.Render(jobViewContent)
	vp.SetContent(str)

	grid := lipgloss.JoinVertical(lipgloss.Left, titleStyle.Render("Just Do It!\n"),
		lipgloss.JoinHorizontal(lipgloss.Top, m.list.View(), vp.View()))
	return grid
}

func (s *Schedule) GetSchedule() error {
	// addr should be configurable
	r, err := http.Get("http://localhost:8081/schedule")
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(s)
}

func TUI() {
	// init schedule schedule
	schedule := &Schedule{}
	if err := schedule.GetSchedule(); err != nil {
		fmt.Printf("Error connecting with JDI server: %v\n", err.Error())
		os.Exit(1)
	}

	items := []list.Item{}
	for _, v := range schedule.Jobs {
		v.LoadRuns()
		item := item{title: v.GetTitle(), desc: v.GetStatusDescription(), jobName: v.Name}
		items = append(items, item)
		// get run history for each job
	}

	const defaultWidth = 20

	l := list.NewModel(items, list.NewDefaultDelegate(), defaultWidth, listHeight)
	l.Title = "Jobs"
	// l.SetShowStatusBar(false)
	// l.SetFilteringEnabled(false)
	// l.Styles.Title = titleStyle
	// l.Styles.PaginationStyle = paginationStyle
	// l.Styles.HelpStyle = helpStyle

	m := model{list: l, state: schedule}
	if len(items) > 0 {
		m.choice = items[len(items)-1].(item).jobName
	}

	if err := tea.NewProgram(m, tea.WithAltScreen()).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
