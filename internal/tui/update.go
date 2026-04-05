package tui

import (
	"strings"
	"time"

	"filepass/internal/pages"
	"filepass/internal/services"

	tea "charm.land/bubbletea/v2"
)

type configLoadedMsg struct {
	servers map[string]services.Server
	err     error
}

type serverAddedMsg struct {
	name    string
	servers map[string]services.Server
}

type clearFlashMsg struct{}

func clearFlashAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(time.Time) tea.Msg {
		return clearFlashMsg{}
	})
}

// isDisabled reports whether a menu item is non-interactive given current state.
func (m TUIInterface) isDisabled(i int) bool {
	return m.MenuItems[i].RequiresServers && m.NoServers
}

// nextSelectable finds the next non-disabled index in direction (+1 or -1).
func (m TUIInterface) nextSelectable(from, dir int) int {
	i := from + dir
	for i >= 0 && i < len(m.MenuItems) {
		if !m.isDisabled(i) {
			return i
		}
		i += dir
	}
	return from
}

func (m TUIInterface) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case pages.HomePageMsg:
		m.Page = pageHome
		m.MenuItems = pages.HomeMenuItems()
		m.Selected = 0
		return m, nil

	case pages.ConfigPageMsg:
		m.Page = pageConfig
		m.MenuItems = pages.ConfigMenuItems()
		m.Selected = 0
		return m, nil

	case pages.AddServerPageMsg:
		m.Page = pageAddServer
		m.Form = newAddServerForm()
		m.FormErr = ""
		return m, nil

	case configLoadedMsg:
		if msg.err != nil {
			m.InitErr = msg.err
			return m, nil
		}
		m.Servers = msg.servers
		m.NoServers = len(msg.servers) == 0
		return m, nil

	case serverAddedMsg:
		m.Servers = msg.servers
		m.NoServers = len(msg.servers) == 0
		m.Page = pageConfig
		m.MenuItems = pages.ConfigMenuItems()
		m.Selected = 0
		m.FlashMsg = "✓  \"" + msg.name + "\" added successfully."
		return m, clearFlashAfter(2 * time.Second)

	case clearFlashMsg:
		m.FlashMsg = ""
		return m, nil

	case tea.WindowSizeMsg:
		m.WindowWidth = msg.Width
		m.WindowHeight = msg.Height
		return m, nil

	case tea.KeyPressMsg:
		// add server form has its own key handling
		if m.Page == pageAddServer {
			return m.updateAddServer(msg)
		}

		switch msg.String() {
		case "up", "k":
			m.Selected = m.nextSelectable(m.Selected, -1)
		case "down", "j":
			m.Selected = m.nextSelectable(m.Selected, +1)
		case "enter":
			if m.isDisabled(m.Selected) {
				return m, nil
			}
			switch m.MenuItems[m.Selected].Key {
			case "exit":
				m.Quitting = true
				return m, tea.Quit
			case "config":
				return m, func() tea.Msg { return pages.ConfigPageMsg{} }
			case "back":
				return m, func() tea.Msg { return pages.HomePageMsg{} }
			case "add":
				return m, func() tea.Msg { return pages.AddServerPageMsg{} }
			// TODO: "server", "edit", "remove"
			}
		case "ctrl+c":
			m.Quitting = true
			return m, tea.Quit
		case "esc":
			if m.Page == pageConfig {
				return m, func() tea.Msg { return pages.HomePageMsg{} }
			}
			m.Quitting = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m TUIInterface) updateAddServer(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	f := m.Form

	switch msg.String() {
	case "tab", "down":
		m.Form = f.focusNext()
		return m, nil

	case "shift+tab", "up":
		m.Form = f.focusPrev()
		return m, nil

	case "enter":
		// on an input field, advance to next
		if f.focused < fieldSave {
			m.Form = f.focusNext()
			return m, nil
		}
		// on save button
		if f.focused == fieldSave {
			return m.submitAddServer()
		}
		// on back button
		if f.focused == fieldBack {
			return m, func() tea.Msg { return pages.ConfigPageMsg{} }
		}

	case "ctrl+c":
		m.Quitting = true
		return m, tea.Quit

	case "esc":
		return m, func() tea.Msg { return pages.ConfigPageMsg{} }
	}

	// route keystrokes to the focused input
	if f.focused < len(f.inputs) {
		var cmd tea.Cmd
		f.inputs[f.focused], cmd = f.inputs[f.focused].Update(msg)
		// clear duplicate-name error when user edits the name field
		if f.focused == fieldName {
			m.FormErr = ""
		}
		m.Form = f
		return m, cmd
	}

	return m, nil
}

func (m TUIInterface) submitAddServer() (tea.Model, tea.Cmd) {
	f := m.Form
	name := strings.TrimSpace(f.inputs[fieldName].Value())

	if !f.canSave() {
		return m, nil
	}

	if m.Services.Config.HasServer(name) {
		m.FormErr = "✗  \"" + name + "\" already exists."
		return m, nil
	}

	s := services.Server{
		Host:       strings.TrimSpace(f.inputs[fieldHost].Value()),
		User:       strings.TrimSpace(f.inputs[fieldUser].Value()),
		PrivateKey: strings.TrimSpace(f.inputs[fieldPrivateKey].Value()),
		Port:       strings.TrimSpace(f.inputs[fieldPort].Value()),
	}

	if err := m.Services.Config.AddServer(name, s); err != nil {
		m.FormErr = "✗  " + err.Error()
		return m, nil
	}

	return m, func() tea.Msg {
		return serverAddedMsg{name: name, servers: m.Services.Config.Servers()}
	}
}
