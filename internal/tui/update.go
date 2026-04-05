package tui

import (
	"filepass/internal/pages"

	tea "charm.land/bubbletea/v2"
)

func (m TUIInterface) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case pages.HomePageMsg:
		return m, nil

	case configLoadedMsg:
		if msg.err != nil {
			m.InitErr = msg.err
			return m, nil
		}
		m.Servers = msg.servers
		m.NoServers = len(msg.servers) == 0
		return m, nil

	case tea.WindowSizeMsg:
		m.WindowWidth = msg.Width
		m.WindowHeight = msg.Height
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "up", "k":
			if m.Selected > 0 {
				m.Selected--
			}
		case "down", "j":
			if m.Selected < len(m.MenuItems)-1 {
				m.Selected++
			}
		case "enter":
			if m.MenuItems[m.Selected].Key == "exit" {
				m.Quitting = true
				return m, tea.Quit
			}
			// TODO: dispatch to server/config pages
		case "ctrl+c", "esc":
			m.Quitting = true
			return m, tea.Quit
		}
	}

	return m, nil
}
