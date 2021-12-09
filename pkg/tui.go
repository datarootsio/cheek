package butt

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
	"github.com/charmbracelet/lipgloss"
)

const (
	listWidth        = 14
	headerHeight     = 2
	footerHeight     = 2
	focusBorderColor = "228"
)

var (
	warningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500"))
	faintStyle   = lipgloss.NewStyle().Faint(true)
	titleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#49f770")).Bold(true)
)

type item struct {
	title, desc, jobName string
}

func (j *JobSpec) GetTitle() string {
	if len(j.runs) > 0 && j.runs[0].Status != 0 {
		return j.Name + " " + warningStyle.Bold(true).Render("!")
	}
	return j.Name
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	list          list.Model
	state         *Schedule
	choice        string
	quitting      bool
	width         int
	height        int
	ready         bool
	listFocus     bool
	viewportFocus bool
	hx            string
	httpPort      string
	viewport      viewport.Model
}

func (j *JobSpec) RunInfo() string {
	var runInfo string
	if len(j.runs) == 0 {
		runInfo = "no run history"
	} else if j.runs[0].Status == 0 {
		since := time.Since(j.runs[0].TriggeredAt).String()
		runInfo = "ran " + since + " ago"
	} else {
		runInfo += warningStyle.Render("error'd")
	}

	return runInfo

}

func (j *JobSpec) View(maxWidth int) string {

	var sb strings.Builder

	if len(j.runs) == 0 {
		sb.WriteString("no run history")
		return sb.String()
	}

	for _, jr := range j.runs {
		sb.WriteString(faintStyle.Render(jr.TriggeredAt.String()))
		sb.WriteString("\n")
		sb.WriteString(hardWrap(jr.Log, maxWidth))
		sb.WriteString("\n\n")

	}

	return sb.String()

}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

		m.viewport = viewport.Model{Width: msg.Width - listWidth, Height: msg.Height - headerHeight - footerHeight - 3}
		if !m.ready {
			m.viewport.SetContent("")
			m.ready = true
		}

		m.list.SetHeight(msg.Height - footerHeight - headerHeight)

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "left":
			m.listFocus = true
			m.viewportFocus = !m.listFocus
		case "right":
			m.viewportFocus = true
			m.listFocus = !m.viewportFocus
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				if i.jobName != m.choice {
					// m.ready = false
					m.choice = i.jobName
					j := m.state.Jobs[m.choice]
					m.viewport.SetContent(j.View(m.viewport.Width - 2))
				}

			}
		}
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd

	if m.viewportFocus {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}
	if m.listFocus {
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)

}

func (m model) View() string {
	var j JobSpec
	if _, ok := m.state.Jobs[m.choice]; ok {
		j = *m.state.Jobs[m.choice]
	} else {
		j = JobSpec{}
	}

	title := titleStyle.Width(m.width).Render("butt: Better Unified Time-Driven Triggers")

	jobListStyle := lipgloss.NewStyle().Border(lipgloss.NormalBorder())

	if m.listFocus {
		jobListStyle = jobListStyle.BorderForeground(lipgloss.Color(focusBorderColor))
	}

	jobList := jobListStyle.Render(m.list.View())

	jobTitle := lipgloss.NewStyle().Foreground(lipgloss.Color("#49f770")).Bold(true).Render(j.Name)
	jobStatus := lipgloss.NewStyle().Faint(true).Align(lipgloss.Right).PaddingRight(1).Width(m.width - lipgloss.Width(jobTitle) - lipgloss.Width(jobList) - 4).Render(j.RunInfo())

	headerBorder := lipgloss.Border{
		Bottom: "_.-.",
	}
	header := lipgloss.NewStyle().Border(headerBorder).BorderTop(false).MarginBottom(1).Render(lipgloss.JoinHorizontal(lipgloss.Left, jobTitle, jobStatus))

	hx := faintStyle.Align(lipgloss.Right).Render(m.hx)

	// job view
	vpBox := lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1).Render(m.viewport.View())

	// job box
	jobBoxStyle := lipgloss.NewStyle().Border(lipgloss.NormalBorder())

	if m.viewportFocus {
		jobBoxStyle.BorderForeground(lipgloss.Color(focusBorderColor))
	}

	jobBox := jobBoxStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, header, vpBox))

	mv := lipgloss.JoinVertical(lipgloss.Left, title, lipgloss.JoinHorizontal(lipgloss.Top, jobList, jobBox), hx)

	return mv
}

func (s *Schedule) GetSchedule(httpPort string) error {
	// addr should be configurable
	r, err := http.Get(fmt.Sprintf("http://localhost:%s/schedule", httpPort))
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(s)
}

func TUI(httpPort string) {
	// init schedule schedule
	schedule := &Schedule{}
	if err := schedule.GetSchedule(httpPort); err != nil {
		fmt.Printf("Error connecting with butt server: %v\n", err.Error())
		os.Exit(1)
	}

	items := []list.Item{}
	for _, v := range schedule.Jobs {
		v.LoadRuns()
		item := item{title: v.GetTitle(), jobName: v.Name}
		items = append(items, item)
		// get run history for each job
	}

	id := list.NewDefaultDelegate()
	id.ShowDescription = false
	id.SetSpacing(0)

	l := list.NewModel(items, id, listWidth, 10)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowTitle(false)

	m := model{list: l, state: schedule, listFocus: true, hx: Hex.Poke(), httpPort: httpPort}

	if err := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion()).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
