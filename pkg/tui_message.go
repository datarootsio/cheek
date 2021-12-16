package cheek

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type MessageType int32

const (
	Info MessageType = iota
	Error
)

// A model that shows a notification to the user
type notificationModel struct {
	// The contents of the notification
	msg string
	// The type of the notification
	msgType MessageType
	// The state to return to after the notification is shown
	returningState func() (tea.Model, tea.Cmd)
}

func (model notificationModel) Render() string {
	var content string
	switch model.msgType {
	case Error:
		content = fmt.Sprintf(`ಠ_ಠ That's an error: %v
    If you think this is a bug then let us know at https://github.com/datarootsio/cheek
    `, model.msg)
	case Info:
		content = model.msg
	}
	return fmt.Sprintf("%v\n\nPress any key to continue...", content)
}

func (model notificationModel) Init() tea.Cmd {
	return nil
}

func (model notificationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return model.returningState()
	}
	return model, nil
}

func (model notificationModel) View() string {
	return model.Render()
}
