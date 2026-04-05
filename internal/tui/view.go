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

func (m TUIInterface) subtitle() string {
	switch m.Page {
	case pageConfig:
		return "Configuration"
	case pageAddServer:
		return "Add Server"
	default:
		return "Secure file transfer"
	}
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

	var body string
	switch m.Page {
	case pageAddServer:
		body = m.viewAddServer()
	default:
		body = m.viewMenu()
	}

	header := lipgloss.JoinVertical(lipgloss.Left,
		styles.CardTitleStyle.Render("✦  filepass"),
		styles.CardSubtitleStyle.Render(m.subtitle()),
	)

	topContent := styles.CardInnerStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, header, body),
	)

	var footerStr string
	if m.Page == pageAddServer {
		footerStr = footerHint("tab/↑↓", "navigate") +
			footerSep() +
			footerHint("enter", "confirm") +
			footerSep() +
			footerHint("esc", "back")
	} else {
		footerStr = footerHint("↑↓", "navigate") +
			footerSep() +
			footerHint("enter", "select") +
			footerSep() +
			footerHint("esc", "quit")
	}
	footer := styles.FooterStyle.Render(footerStr)

	card := styles.CardStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, topContent, footer),
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

func (m TUIInterface) viewMenu() string {
	var menuRows []string
	for i, item := range m.MenuItems {
		disabled := m.isDisabled(i)
		menuRows = append(menuRows, styles.MenuItemStyle(i == m.Selected, disabled).Render(item.Label))
	}
	menu := lipgloss.JoinVertical(lipgloss.Left, menuRows...)

	var statusLine string
	switch {
	case m.InitErr != nil:
		statusLine = styles.StatusErrStyle.Render("✗  " + m.InitErr.Error())
	case m.NoServers && m.Page == pageHome:
		statusLine = styles.StatusWarnStyle.Render("⚠  No servers configured. Select Config to add one.")
	case m.FlashMsg != "" && m.Page == pageConfig:
		statusLine = styles.StatusOKStyle.Render(m.FlashMsg)
	}

	if statusLine != "" {
		return lipgloss.JoinVertical(lipgloss.Left, menu, statusLine)
	}
	return menu
}

func (m TUIInterface) viewAddServer() string {
	f := m.Form
	labels := []string{"Name", "Host", "User", "Private Key Path", "Port"}
	required := []bool{true, true, true, true, false}

	var rows []string
	for i, label := range labels {
		lbl := styles.FieldLabelStyle(required[i]).Render(label)
		input := f.inputs[i].View()
		rows = append(rows, lipgloss.JoinVertical(lipgloss.Left, lbl, input))
	}
	form := lipgloss.JoinVertical(lipgloss.Left, rows...)

	// required legend
	legend := styles.FieldLegendStyle.Render("* required")

	// form error (duplicate name, etc.)
	var errLine string
	if m.FormErr != "" {
		errLine = styles.StatusErrStyle.Render(m.FormErr)
	}

	// save / back buttons
	saveBtn := styles.ButtonStyle(f.focused == fieldSave, f.canSave()).Render("Save")
	backBtn := styles.ButtonStyle(f.focused == fieldBack, true).Render("Back")
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, saveBtn, "  ", backBtn)

	parts := []string{form, legend}
	if errLine != "" {
		parts = append(parts, errLine)
	}
	parts = append(parts, buttons)

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}
