package cheek

import (
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rs/zerolog"
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

	// try a key press
	tm, _ = m.Update(tea.Key{Type: tea.KeyRunes, Runes: []rune("r")})

	// assert another always on screen value
	assert.Contains(t, tm.View(), "efresh")
}
