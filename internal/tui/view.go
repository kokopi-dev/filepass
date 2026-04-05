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
	case pageSelectServer:
		return "Select Server"
	case pageServerActions:
		if m.ActiveServer != "" {
			return m.ActiveServer
		}
		return "Server"
	case pageFileAction:
		return m.ActiveFile
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
	case pageSelectServer:
		body = m.viewSelectServer()
	case pageServerActions:
		body = m.viewServerActions()
	case pageFileAction:
		body = m.viewFileAction()
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
	switch m.Page {
	case pageAddServer:
		footerStr = footerHint("tab/↑↓", "navigate") +
			footerSep() +
			footerHint("enter", "confirm") +
			footerSep() +
			footerHint("ctrl+v", "paste") +
			footerSep() +
			footerHint("esc", "back")
	case pageSelectServer:
		footerStr = footerHint("↑↓", "navigate") +
			footerSep() +
			footerHint("enter", "connect") +
			footerSep() +
			footerHint("esc", "back")
	case pageServerActions:
		footerStr = footerHint("tab", "switch pane") +
			footerSep() +
			footerHint("↑↓", "navigate") +
			footerSep() +
			footerHint("enter", "select") +
			footerSep() +
			footerHint("esc", "back")
	case pageFileAction:
		footerStr = footerHint("↑↓", "navigate") +
			footerSep() +
			footerHint("enter", "confirm") +
			footerSep() +
			footerHint("esc", "back")
	default:
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

func (m TUIInterface) viewServerActions() string {
	// action menu — single column, unfocused when file pane is active
	var actionRows []string
	for i, item := range m.MenuItems {
		active := !m.FileFocused && i == m.Selected
		actionRows = append(actionRows, styles.MenuItemStyle(active, false).Render(item.Label))
	}
	actions := lipgloss.JoinVertical(lipgloss.Left, actionRows...)

	// file list section below, separated by a top border
	var fileRows []string
	switch {
	case m.StorageLoading:
		fileRows = append(fileRows, styles.StatusWarnStyle.Render("  loading…"))
	case m.StorageErr != nil:
		fileRows = append(fileRows, styles.StatusErrStyle.Render("✗  "+m.StorageErr.Error()))
	case len(m.StorageFiles) == 0:
		fileRows = append(fileRows, styles.StorageEmptyStyle.Render("  no files in storage"))
	default:
		for i, f := range m.StorageFiles {
			active := m.FileFocused && i == m.FileSelected
			fileRows = append(fileRows, styles.FileItemStyle(active).Render(f))
		}
	}
	fileList := lipgloss.JoinVertical(lipgloss.Left, fileRows...)
	fileSection := styles.StorageFileSectionStyle.Render(fileList)

	return lipgloss.JoinVertical(lipgloss.Left, actions, fileSection)
}

func (m TUIInterface) viewFileAction() string {
	// filename shown as a dim label above the menu
	filenameLabel := styles.FilenameLabelStyle.Render(m.ActiveFile)

	var menuRows []string
	for i, item := range m.MenuItems {
		menuRows = append(menuRows, styles.MenuItemStyle(i == m.Selected, false).Render(item.Label))
	}
	menu := lipgloss.JoinVertical(lipgloss.Left, menuRows...)

	return lipgloss.JoinVertical(lipgloss.Left, filenameLabel, menu)
}

func (m TUIInterface) viewSelectServer() string {
	if len(m.ServerNames) == 0 {
		return styles.StatusWarnStyle.Render("⚠  No servers configured.")
	}

	var rows []string
	for i, name := range m.ServerNames {
		rows = append(rows, styles.ServerRowStyle(i == m.Selected, name))
	}
	return lipgloss.JoinVertical(lipgloss.Left, rows...)
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

	legend := styles.FieldLegendStyle.Render("* required")

	var errLine string
	if m.FormErr != "" {
		errLine = styles.StatusErrStyle.Render(m.FormErr)
	}

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
