package styles

import lipgloss "charm.land/lipgloss/v2"

var (
	// Card / box
	CardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Width(52)

	CardInnerStyle = lipgloss.NewStyle().
			Padding(1, 3)

	CardTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true).
			MarginBottom(1)

	CardSubtitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("245")).
				MarginBottom(1)

	// Menu items
	menuItemBase = lipgloss.NewStyle().
			PaddingLeft(2).
			Width(44)

	menuItemActive = menuItemBase.
			Foreground(lipgloss.Color("75")).
			Bold(true).
			SetString("▸ ")

	menuItemInactive = menuItemBase.
				Foreground(lipgloss.Color("245"))

	// Status lines
	StatusWarnStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("221")).
			MarginTop(1)

	StatusErrStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("203")).
			MarginTop(1)

	// Footer
	FooterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			BorderTop(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("237")).
			Padding(0, 1).
			Width(50)

	FooterKeyStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true)
	FooterSepStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("237"))
	FooterDescStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("243"))
)

// MenuItemStyle returns the appropriate style for a menu row.
func MenuItemStyle(active bool) lipgloss.Style {
	if active {
		return menuItemActive
	}
	return menuItemInactive
}
