package cheek

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rs/zerolog/log"
)

const (
	listWidth        = 14
	headerHeight     = 2
	footerHeight     = 2
	focusBorderColor = "228"
)

var (
	serverPort   string
	yamlFile     string
	warningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500"))
	faintStyle   = lipgloss.NewStyle().Faint(true)
	titleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#49f770")).Bold(true)
)

type item struct {
	title, desc, jobName string
}

func (j *JobSpec) getTitle() string {
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
	httpPort      string
	viewport      viewport.Model
	notification  notification
}

type notificationType int32

const (
	Info notificationType = iota
	Error
)

type notification struct {
	content          string
	notificationType notificationType
}

func (j *JobSpec) runInfo() string {
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

func (j *JobSpec) view(maxWidth int) string {
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
	return refreshState
}

func refreshState() tea.Msg {
	schedule := &Schedule{}
	if err := schedule.getSchedule(serverPort, yamlFile); err != nil {
		log.Error().Err(err).Msg("Failed to load schedules")
		return notification{
			content:          "Can't refresh run info",
			notificationType: Error,
		}
	}

	for _, v := range schedule.Jobs {
		v.loadRuns()
	}
	return schedule
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case *Schedule:

		keys := make([]string, 0)
		for k := range msg.Jobs {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		items := []list.Item{}
		for _, k := range keys {
			j := msg.Jobs[k]
			item := item{title: j.getTitle(), jobName: j.Name}
			items = append(items, item)
		}

		m.list.SetItems(items)
		m.state = msg
		m.choice = m.list.SelectedItem().(item).jobName
		j := m.state.Jobs[m.choice]
		m.viewport.SetContent(j.view(m.viewport.Width - 2))

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
		case "r":
			return m, refreshState
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
			if ok && i.jobName != m.choice {
				m.choice = i.jobName
				j := m.state.Jobs[m.choice]
				m.viewport.SetContent(j.view(m.viewport.Width - 2))
			}
		}

	case notification:
		m.notification = msg
	}

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

	refresh := faintStyle.Align(lipgloss.Right).Render("(r)efresh")
	title := titleStyle.Width(m.width - lipgloss.Width(refresh)).Render("cheek |_|>")
	header := lipgloss.JoinHorizontal(lipgloss.Left, title, refresh)

	jobListStyle := lipgloss.NewStyle().Border(lipgloss.NormalBorder())

	if m.listFocus {
		jobListStyle = jobListStyle.BorderForeground(lipgloss.Color(focusBorderColor))
	}

	jobList := jobListStyle.Render(m.list.View())

	jobTitle := lipgloss.NewStyle().Foreground(lipgloss.Color("#49f770")).Bold(true).Render(j.Name)
	jobStatus := lipgloss.NewStyle().Faint(true).Align(lipgloss.Right).PaddingRight(1).Width(m.width - lipgloss.Width(jobTitle) - lipgloss.Width(jobList) - 4).Render(j.runInfo())

	logBoxHeaderBorder := lipgloss.Border{
		Bottom: "_.-.",
	}
	logBoxHeader := lipgloss.NewStyle().Border(logBoxHeaderBorder).BorderTop(false).MarginBottom(1).Render(lipgloss.JoinHorizontal(lipgloss.Left, jobTitle, jobStatus))

	var notification string
	switch m.notification.notificationType {
	case Info:
		notification = m.notification.content
	case Error:
		notification = lipgloss.
			NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#FF5F87")).
			Align(lipgloss.Right).
			Render(m.notification.content)
	}

	// job view
	vpBox := lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1).Render(m.viewport.View())

	// job box
	jobBoxStyle := lipgloss.NewStyle().Border(lipgloss.NormalBorder())

	if m.viewportFocus {
		jobBoxStyle.BorderForeground(lipgloss.Color(focusBorderColor))
	}

	jobBox := jobBoxStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, logBoxHeader, vpBox))

	mv := lipgloss.JoinVertical(lipgloss.Left, header, lipgloss.JoinHorizontal(lipgloss.Top, jobList, jobBox), notification)

	return mv
}

func (s *Schedule) getSchedule(httpPort string, scheduleFile string) error {
	// addr should be configurable
	r, server_err := http.Get(fmt.Sprintf("http://localhost:%s/schedule", httpPort))
	if server_err == nil {
		defer r.Body.Close()
		return json.NewDecoder(r.Body).Decode(s)
	}
	if scheduleFile != "" {
		schedule, err := loadSchedule(scheduleFile)
		if err != nil {
			return fmt.Errorf("%w\nError reading YAML: %v", server_err, err.Error())
		}
		*s = schedule
		return nil
	}
	return fmt.Errorf("Error connecting to cheek server and -s is not set: %w", server_err)
}

// TUI is the main entrypoint for the cheek ui.
func TUI(httpPort string, scheduleFile string) {
	serverPort = httpPort
	yamlFile = scheduleFile
	// init schedule schedule
	schedule := &Schedule{}
	if err := schedule.getSchedule(httpPort, scheduleFile); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	items := []list.Item{}
	id := list.NewDefaultDelegate()
	id.ShowDescription = false
	id.SetSpacing(0)

	l := list.NewModel(items, id, listWidth, 10)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowTitle(false)

	m := model{list: l, state: schedule, listFocus: true, httpPort: httpPort}

	if err := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion()).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
