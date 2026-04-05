package tui

import (
	"sort"
	"strings"
	"time"

	"filepass/internal/pages"
	"filepass/internal/services"

	"charm.land/bubbles/v2/textinput"
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

type serverRemovedMsg struct {
	name    string
	servers map[string]services.Server
}

type clearFlashMsg struct{}

type storageFilesMsg struct {
	files []string
	err   error
}

type fileOpMsg struct {
	op      string // "get" or "delete"
	err     error
	success string
}

func getFileCmd(store *services.ServicesStore, serverName, filename, destDir string) tea.Cmd {
	return func() tea.Msg {
		storage, err := store.NewStorageService(serverName)
		if err != nil {
			return fileOpMsg{op: "get", err: err}
		}
		if err := storage.Get(filename, destDir); err != nil {
			return fileOpMsg{op: "get", err: err}
		}
		return fileOpMsg{op: "get", success: "✓  Downloaded \"" + filename + "\""}
	}
}

func deleteFileCmd(store *services.ServicesStore, serverName, filename string) tea.Cmd {
	return func() tea.Msg {
		storage, err := store.NewStorageService(serverName)
		if err != nil {
			return fileOpMsg{op: "delete", err: err}
		}
		if err := storage.Delete(filename); err != nil {
			return fileOpMsg{op: "delete", err: err}
		}
		return fileOpMsg{op: "delete", success: "✓  Deleted \"" + filename + "\""}
	}
}

type cleanAllMsg struct {
	err error
}

func cleanAllCmd(store *services.ServicesStore, serverName string) tea.Cmd {
	return func() tea.Msg {
		storage, err := store.NewStorageService(serverName)
		if err != nil {
			return cleanAllMsg{err: err}
		}
		return cleanAllMsg{err: storage.CleanAll()}
	}
}

func newCleanInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "yes"
	ti.Prompt = ""
	ti.CharLimit = 3
	ti.SetWidth(10)
	ti.Focus()
	return ti
}

func checkStorageCmd(store *services.ServicesStore, serverName string) tea.Cmd {
	return func() tea.Msg {
		storage, err := store.NewStorageService(serverName)
		if err != nil {
			return storageFilesMsg{err: err}
		}
		files, err := storage.Check()
		return storageFilesMsg{files: files, err: err}
	}
}

func clearFlashAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(time.Time) tea.Msg {
		return clearFlashMsg{}
	})
}

// sortedServerNames returns the keys of servers sorted alphabetically.
func sortedServerNames(servers map[string]services.Server) []string {
	names := make([]string, 0, len(servers))
	for name := range servers {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
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

	case pages.SelectServerPageMsg:
		m.Page = pageSelectServer
		m.Selected = 0
		return m, nil

	case pages.ServerActionsPageMsg:
		m.Page = pageServerActions
		m.ActiveServer = msg.ServerName
		m.MenuItems = pages.ServerActionItems()
		m.Selected = 0
		m.FileSelected = 0
		m.FileFocused = true
		m.StorageFiles = nil
		m.StorageErr = nil
		m.StorageLoading = true
		return m, checkStorageCmd(m.Services, msg.ServerName)

	case pages.FileActionPageMsg:
		m.Page = pageFileAction
		m.ActiveFile = msg.Filename
		m.MenuItems = pages.FileActionItems()
		m.Selected = 0
		m.FileOpLoading = false
		m.FileOpErr = nil
		m.FileOpSuccess = ""
		return m, nil

	case pages.RemoveServerPageMsg:
		m.Page = pageRemoveServer
		m.Selected = 0
		return m, nil

	case serverRemovedMsg:
		m.Servers = msg.servers
		m.ServerNames = sortedServerNames(msg.servers)
		m.NoServers = len(msg.servers) == 0
		m.Page = pageConfig
		m.MenuItems = pages.ConfigMenuItems()
		m.Selected = 0
		m.FlashMsg = "✓  \"" + msg.name + "\" removed."
		return m, clearFlashAfter(2 * time.Second)

	case pages.CleanAllPageMsg:
		m.Page = pageCleanAll
		m.CleanInput = newCleanInput()
		m.CleanOpLoading = false
		m.CleanOpErr = nil
		return m, nil

	case cleanAllMsg:
		m.CleanOpLoading = false
		if msg.err != nil {
			m.CleanOpErr = msg.err
			return m, nil
		}
		server := m.ActiveServer
		return m, func() tea.Msg { return pages.ServerActionsPageMsg{ServerName: server} }

	case pages.SendPageMsg:
		m.Page = pageSend
		m.Picker = newPicker(m.LocalDir)
		return m, nil

	case fileOpMsg:
		m.FileOpLoading = false
		if msg.err != nil {
			m.FileOpErr = msg.err
			return m, nil
		}
		// on success: return to server actions and refresh file list
		server := m.ActiveServer
		return m, tea.Batch(
			func() tea.Msg { return pages.ServerActionsPageMsg{ServerName: server} },
			clearFlashAfter(0), // immediate clear of any stale flash
		)

	case storageFilesMsg:
		m.StorageLoading = false
		m.StorageFiles = msg.files
		m.StorageErr = msg.err
		m.FileSelected = 0
		return m, nil

	case configLoadedMsg:
		if msg.err != nil {
			m.InitErr = msg.err
			return m, nil
		}
		m.Servers = msg.servers
		m.ServerNames = sortedServerNames(msg.servers)
		m.NoServers = len(msg.servers) == 0
		return m, nil

	case serverAddedMsg:
		m.Servers = msg.servers
		m.ServerNames = sortedServerNames(msg.servers)
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
		if m.Page == pageAddServer {
			return m.updateAddServer(msg)
		}
		if m.Page == pageSelectServer {
			return m.updateSelectServer(msg)
		}
		if m.Page == pageServerActions {
			return m.updateServerActions(msg)
		}
		if m.Page == pageFileAction {
			return m.updateFileAction(msg)
		}
		if m.Page == pageSend {
			return m.updateSend(msg)
		}
		if m.Page == pageCleanAll {
			return m.updateCleanAll(msg)
		}
		if m.Page == pageRemoveServer {
			return m.updateRemoveServer(msg)
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
			case "server":
				return m, func() tea.Msg { return pages.SelectServerPageMsg{} }
			case "remove":
				return m, func() tea.Msg { return pages.RemoveServerPageMsg{} }
			// TODO: "edit"
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

	case tea.PasteMsg:
		if m.Page == pageAddServer {
			return m.updateAddServerPaste(msg.Content)
		}

	case tea.ClipboardMsg:
		if m.Page == pageAddServer {
			return m.updateAddServerPaste(msg.Content)
		}
	}

	return m, nil
}

func (m TUIInterface) updateServerActions(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab":
		if len(m.StorageFiles) > 0 || !m.StorageLoading {
			m.FileFocused = !m.FileFocused
		}

	case "up", "k":
		if m.FileFocused {
			if m.FileSelected > 0 {
				m.FileSelected--
			}
		} else {
			if m.Selected > 0 {
				m.Selected--
			}
		}

	case "down", "j":
		if m.FileFocused {
			if m.FileSelected < len(m.StorageFiles)-1 {
				m.FileSelected++
			}
		} else {
			if m.Selected < len(m.MenuItems)-1 {
				m.Selected++
			}
		}

	case "enter":
		if m.FileFocused && len(m.StorageFiles) > 0 {
			file := m.StorageFiles[m.FileSelected]
			server := m.ActiveServer
			return m, func() tea.Msg {
				return pages.FileActionPageMsg{ServerName: server, Filename: file}
			}
		}
		server := m.ActiveServer
		switch m.MenuItems[m.Selected].Key {
		case "send":
			return m, func() tea.Msg { return pages.SendPageMsg{ServerName: server} }
		case "clean":
			return m, func() tea.Msg { return pages.CleanAllPageMsg{ServerName: server} }
		}

	case "ctrl+c":
		m.Quitting = true
		return m, tea.Quit

	case "esc":
		return m, func() tea.Msg { return pages.SelectServerPageMsg{} }
	}

	return m, nil
}

func (m TUIInterface) updateFileAction(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	// block input while an operation is running
	if m.FileOpLoading {
		if msg.String() == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
		return m, nil
	}

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
		server := m.ActiveServer
		file := m.ActiveFile
		switch m.MenuItems[m.Selected].Key {
		case "get":
			m.FileOpLoading = true
			m.FileOpErr = nil
			m.FileOpSuccess = ""
			return m, getFileCmd(m.Services, server, file, m.LocalDir)
		case "delete":
			m.FileOpLoading = true
			m.FileOpErr = nil
			m.FileOpSuccess = ""
			return m, deleteFileCmd(m.Services, server, file)
		}
	case "ctrl+c":
		m.Quitting = true
		return m, tea.Quit
	case "esc":
		server := m.ActiveServer
		return m, func() tea.Msg { return pages.ServerActionsPageMsg{ServerName: server} }
	}
	return m, nil
}

func (m TUIInterface) updateSend(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	p := m.Picker

	switch msg.String() {
	case "up", "k":
		if p.cursor > 0 {
			p.cursor--
		}
		m.Picker = p
		return m, nil

	case "down", "j":
		if p.cursor < len(p.filtered)-1 {
			p.cursor++
		}
		m.Picker = p
		return m, nil

	case "enter":
		if len(p.filtered) == 0 {
			return m, nil
		}
		selected := p.filtered[p.cursor]
		if selected.isDir {
			m.Picker = p.descend(selected.name)
			return m, nil
		}
		// file selected — send it
		path := p.selectedPath()
		_ = path // TODO: wire to send service
		server := m.ActiveServer
		return m, func() tea.Msg { return pages.ServerActionsPageMsg{ServerName: server} }

	case "backspace":
		newP, consumed := p.backspace()
		m.Picker = newP
		_ = consumed
		return m, nil

	case "ctrl+c":
		m.Quitting = true
		return m, tea.Quit

	case "esc":
		server := m.ActiveServer
		return m, func() tea.Msg { return pages.ServerActionsPageMsg{ServerName: server} }

	default:
		// printable single rune → append to query
		if msg.Text != "" {
			m.Picker = p.typeRune([]rune(msg.Text)[0])
			return m, nil
		}
	}

	return m, nil
}

func (m TUIInterface) updateRemoveServer(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	last := len(m.ServerNames) - 1
	switch msg.String() {
	case "up", "k":
		if m.Selected > 0 {
			m.Selected--
		}
	case "down", "j":
		if m.Selected < last {
			m.Selected++
		}
	case "enter":
		if m.Selected >= 0 && m.Selected < len(m.ServerNames) {
			name := m.ServerNames[m.Selected]
			if err := m.Services.Config.RemoveServer(name); err != nil {
				// surface error via flash on config page
				m.Page = pageConfig
				m.MenuItems = pages.ConfigMenuItems()
				m.Selected = 0
				m.FlashMsg = "✗  " + err.Error()
				return m, clearFlashAfter(3 * time.Second)
			}
			servers := m.Services.Config.Servers()
			return m, func() tea.Msg {
				return serverRemovedMsg{name: name, servers: servers}
			}
		}
	case "ctrl+c":
		m.Quitting = true
		return m, tea.Quit
	case "esc":
		return m, func() tea.Msg { return pages.ConfigPageMsg{} }
	}
	return m, nil
}

func (m TUIInterface) updateCleanAll(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	if m.CleanOpLoading {
		if msg.String() == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
		return m, nil
	}

	switch msg.String() {
	case "enter":
		if m.CleanInput.Value() == "yes" {
			m.CleanOpLoading = true
			m.CleanOpErr = nil
			return m, cleanAllCmd(m.Services, m.ActiveServer)
		}
	case "ctrl+c":
		m.Quitting = true
		return m, tea.Quit
	case "esc":
		server := m.ActiveServer
		return m, func() tea.Msg { return pages.ServerActionsPageMsg{ServerName: server} }
	}

	// route all other keys to the text input
	var cmd tea.Cmd
	m.CleanInput, cmd = m.CleanInput.Update(msg)
	return m, cmd
}

func (m TUIInterface) updateSelectServer(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	last := len(m.ServerNames) - 1
	switch msg.String() {
	case "up", "k":
		if m.Selected > 0 {
			m.Selected--
		}
	case "down", "j":
		if m.Selected < last {
			m.Selected++
		}
	case "enter":
		if m.Selected >= 0 && m.Selected < len(m.ServerNames) {
			name := m.ServerNames[m.Selected]
			return m, func() tea.Msg { return pages.ServerActionsPageMsg{ServerName: name} }
		}
	case "ctrl+c":
		m.Quitting = true
		return m, tea.Quit
	case "esc":
		return m, func() tea.Msg { return pages.HomePageMsg{} }
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
		if f.focused < fieldSave {
			m.Form = f.focusNext()
			return m, nil
		}
		if f.focused == fieldSave {
			return m.submitAddServer()
		}
		if f.focused == fieldBack {
			return m, func() tea.Msg { return pages.ConfigPageMsg{} }
		}

	case "ctrl+c":
		m.Quitting = true
		return m, tea.Quit

	case "ctrl+v":
		if f.focused < len(f.inputs) {
			return m, tea.ReadClipboard
		}
		return m, nil

	case "esc":
		return m, func() tea.Msg { return pages.ConfigPageMsg{} }
	}

	if f.focused < len(f.inputs) {
		var cmd tea.Cmd
		f.inputs[f.focused], cmd = f.inputs[f.focused].Update(msg)
		if f.focused == fieldName {
			m.FormErr = ""
		}
		m.Form = f
		return m, cmd
	}

	return m, nil
}

func (m TUIInterface) updateAddServerPaste(text string) (tea.Model, tea.Cmd) {
	f := m.Form
	if f.focused >= len(f.inputs) {
		return m, nil
	}
	var cmd tea.Cmd
	f.inputs[f.focused], cmd = f.inputs[f.focused].Update(tea.PasteMsg{Content: text})
	if f.focused == fieldName {
		m.FormErr = ""
	}
	m.Form = f
	return m, cmd
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
