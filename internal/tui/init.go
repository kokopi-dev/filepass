package tui

import (
	"filepass/internal/pages"
	"filepass/internal/services"

	tea "charm.land/bubbletea/v2"
)

type configLoadedMsg struct {
	servers map[string]services.Server
	err     error
}

func (m TUIInterface) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg { return pages.HomePageMsg{} },
		func() tea.Msg {
			return configLoadedMsg{servers: m.Services.Config.Servers()}
		},
	)
}
