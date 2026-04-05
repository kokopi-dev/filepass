package tui

import (
	"filepass/internal/pages"

	tea "charm.land/bubbletea/v2"
)

func (m TUIInterface) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg { return pages.HomePageMsg{} },
		func() tea.Msg {
			return configLoadedMsg{servers: m.Services.Config.Servers()}
		},
	)
}
