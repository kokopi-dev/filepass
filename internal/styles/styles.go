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

	menuItemDisabled = menuItemBase.
				Foreground(lipgloss.Color("240")).
				PaddingLeft(4)

	// Form fields
	fieldLabelRequired = lipgloss.NewStyle().
				Foreground(lipgloss.Color("75")).
				Bold(true).
				MarginTop(1)

	fieldLabelOptional = lipgloss.NewStyle().
				Foreground(lipgloss.Color("245")).
				MarginTop(1)

	FieldLegendStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Italic(true).
				MarginTop(1)

	// Buttons
	buttonActive = lipgloss.NewStyle().
			Foreground(lipgloss.Color("232")).
			Background(lipgloss.Color("75")).
			Bold(true).
			Padding(0, 2).
			MarginTop(1)

	buttonInactive = lipgloss.NewStyle().
			Foreground(lipgloss.Color("232")).
			Background(lipgloss.Color("240")).
			Padding(0, 2).
			MarginTop(1)

	buttonLocked = lipgloss.NewStyle().
			Foreground(lipgloss.Color("238")).
			Background(lipgloss.Color("235")).
			Padding(0, 2).
			MarginTop(1)

	// Status lines
	StatusOKStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			MarginTop(1)

	StatusWarnStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("221")).
			MarginTop(1)

	StatusErrStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("203")).
			MarginTop(1)

	CleanWarningStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("203")).
				Bold(true).
				MarginBottom(1).
				Width(44)

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

var (
	serverRowBase = lipgloss.NewStyle().
			PaddingLeft(2).
			PaddingTop(0).
			Width(44)

	serverRowBaseActive = lipgloss.NewStyle().
				PaddingLeft(2).
				Width(44)

	serverRowNameStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("255"))

	serverRowNameActiveStyle = lipgloss.NewStyle().
					Bold(true).
					Foreground(lipgloss.Color("75"))

	// Storage file list
	StorageFileSectionStyle = lipgloss.NewStyle().
					BorderTop(true).
					BorderStyle(lipgloss.NormalBorder()).
					BorderForeground(lipgloss.Color("237")).
					MarginTop(1).
					PaddingTop(1).
					Width(44)

	StorageEmptyStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("243")).
				Italic(true)

	fileItemInactive = lipgloss.NewStyle().
				PaddingLeft(4).
				Foreground(lipgloss.Color("252")).
				Width(44)

	fileItemActive = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("75")).
				Bold(true).
				Width(44).
				SetString("▸ ")

	FilenameLabelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("243")).
				MarginBottom(1)

	// Local directory label (above file list and in picker breadcrumb)
	LocalDirStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("243")).
				Italic(true).
				MarginBottom(1)

	// File picker
	PickerQueryStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("75")).
				MarginBottom(1)

	pickerItemBase = lipgloss.NewStyle().
				PaddingLeft(4).
				Width(44)

	pickerItemActive = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("75")).
				Bold(true).
				Width(44).
				SetString("▸ ")

	pickerDirColor  = lipgloss.Color("75")
	pickerFileColor = lipgloss.Color("252")
)

func MenuItemStyle(active, disabled bool) lipgloss.Style {
	switch {
	case disabled:
		return menuItemDisabled
	case active:
		return menuItemActive
	default:
		return menuItemInactive
	}
}

func FieldLabelStyle(required bool) lipgloss.Style {
	if required {
		return fieldLabelRequired
	}
	return fieldLabelOptional
}

// ButtonStyle returns the style for a button.
// focused: cursor is on this button. enabled: button is interactive.
func ButtonStyle(focused, enabled bool) lipgloss.Style {
	switch {
	case !enabled:
		return buttonLocked
	case focused:
		return buttonActive
	default:
		return buttonInactive
	}
}

// FileItemStyle returns the style for a file list row.
func FileItemStyle(active bool) lipgloss.Style {
	if active {
		return fileItemActive
	}
	return fileItemInactive
}

// PickerItemStyle returns the style for a file picker entry.
// Directories are coloured differently from files.
func PickerItemStyle(active, isDir bool) lipgloss.Style {
	if active {
		if isDir {
			return pickerItemActive.Foreground(pickerDirColor)
		}
		return pickerItemActive.Foreground(lipgloss.Color("255"))
	}
	if isDir {
		return pickerItemBase.Foreground(pickerDirColor)
	}
	return pickerItemBase.Foreground(pickerFileColor)
}

// ServerRowStyle renders a single-line server list entry showing only the server name.
func ServerRowStyle(active bool, name string) string {
	if active {
		return serverRowBaseActive.Render(serverRowNameActiveStyle.Render("▸ " + name))
	}
	return serverRowBase.Render(serverRowNameStyle.Render("  " + name))
}
