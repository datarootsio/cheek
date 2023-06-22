package cheek

import (
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestUIBasics(t *testing.T) {
	var tm tea.Model
	s, err := loadSchedule(zerolog.Logger{}, NewConfig(), "../testdata/jobs1.yaml")
	if err != nil {
		t.Fatal(err)
	}
	m := model{list: list.NewModel([]list.Item{}, list.NewDefaultDelegate(), listWidth, 10), state: &s, ready: false}

	// fetch initial model with dummy keymsg
	tm, _ = m.Update(tea.WindowSizeMsg{Width: 30, Height: 20})
	// just some simple assertions for now
	assert.Contains(t, tm.View(), "cheek")

	// try key presses
	_, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")})
	_, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("c")})
	_, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("left")})
	_, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("right")})
	_, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("enter")})

	// assert another always on screen value
	assert.Contains(t, tm.View(), "efresh")

	// test a job view
	j, ok := s.Jobs["bar"]
	assert.True(t, ok)
	assert.Contains(t, j.view(40), "no run history")

	// refresh state
	yamlFile = "../testdata/jobs1.yaml"
	n := refreshState()
	assert.IsType(t, &Schedule{}, n)

	m.Update(n)
}

func TestJobView(t *testing.T) {
	// add a bit of history
	j := JobSpec{Command: []string{"echo", "foo"}, Name: "blaat"}
	jr := j.execCommand("testrun")
	j.finalize(&jr)
	j.Runs = append(j.Runs, jr)
	assert.Contains(t, j.view(120), "foo")

	assert.Contains(t, j.getTitle(), j.Name)
}

func TestRefreshState(t *testing.T) {
	yamlFile = ""
	n := refreshState()
	assert.IsType(t, notification{}, n)
	assert.Contains(t, n.(notification).content, "Can't refresh")

	// provide path to schedule
	yamlFile = "../testdata/jobs1.yaml"
	n = refreshState()

	assert.IsType(t, &Schedule{}, n)
}

func TestUIEntrypoint(t *testing.T) {
	viper.Set("port", 9999)
	err := TUI(zerolog.Logger{}, "../testdata/jobs1.yaml")

	// should execute correctly until starts complaining about interface
	assert.Error(t, err, "device not configured")
}
