package tui

import (
	"filepass/internal/styles"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
)

func footerHint(key, desc string) string {
	return styles.FooterKeyStyle.Render(key) +
		" " +
		styles.FooterDescStyle.Render(desc)
}

func footerSep() string {
	return styles.FooterSepStyle.Render(" · ")
}

func (m TUIInterface) View() tea.View {
	if m.Quitting {
		return tea.NewView("")
	}

	w := m.WindowWidth
	h := m.WindowHeight
	if w == 0 {
		w = 80
	}
	if h == 0 {
		h = 24
	}

	// menu rows
	var menuRows []string
	for i, item := range m.MenuItems {
		menuRows = append(menuRows, styles.MenuItemStyle(i == m.Selected).Render(item.Label))
	}
	menu := lipgloss.JoinVertical(lipgloss.Left, menuRows...)

	// status line — error takes priority over no-servers hint
	var statusLine string
	switch {
	case m.InitErr != nil:
		statusLine = styles.StatusErrStyle.Render("✗  " + m.InitErr.Error())
	case m.NoServers:
		statusLine = styles.StatusWarnStyle.Render("⚠  No servers configured. Select Config to add one.")
	}

	// top content
	innerRows := []string{
		styles.CardTitleStyle.Render("✦  filepass"),
		styles.CardSubtitleStyle.Render("Secure file transfer"),
		menu,
	}
	if statusLine != "" {
		innerRows = append(innerRows, statusLine)
	}
	topContent := styles.CardInnerStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, innerRows...),
	)

	// footer
	hints := footerHint("↑↓", "navigate") +
		footerSep() +
		footerHint("enter", "select") +
		footerSep() +
		footerHint("esc", "quit")
	footer := styles.FooterStyle.Render(hints)

	// card
	card := styles.CardStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			topContent,
			footer,
		),
	)

	cardHeight := lipgloss.Height(card)
	topPad := max((h-cardHeight)/2, 0)

	centeredCard := lipgloss.NewStyle().
		Width(w).
		Align(lipgloss.Center).
		PaddingTop(topPad).
		Render(card)

	v := tea.NewView(centeredCard)
	v.AltScreen = true
	return v
}
