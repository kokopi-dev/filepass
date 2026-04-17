package tui

import (
	"filepass/internal/pages"
	"filepass/internal/services"

	"charm.land/bubbles/v2/textinput"
)

type page int

const (
	pageHome page = iota
	pageConfig
	pageAddServer
	pageSelectServer
	pageServerActions
	pageFileAction
	pageSend
	pageCleanAll
	pageRemoveServer
	pageSelectEditServer
	pageEditServer
)

type TUIInterface struct {
	Services     *services.ServicesStore
	Page         page
	MenuItems    []pages.MenuItem
	Selected     int
	Servers      map[string]services.Server
	ServerNames  []string // sorted, stable order for list rendering
	NoServers    bool
	InitErr      error
	FlashMsg     string
	Form         addServerForm
	FormErr      string // inline field error (e.g. duplicate name)
	Quitting     bool
	WindowWidth  int
	WindowHeight int
	// server actions page
	ActiveServer    string
	LocalDir        string // user's cwd, destination for received files
	StorageFiles    []string
	StorageLoading  bool
	StorageErr      error
	FileSelected    int          // cursor within StorageFiles
	FileFocused     bool         // true = ↑↓ drives file list, false = action menu
	FileScrollOff   int          // first visible row in StorageFiles list
	FileViewHeight  int          // available visible rows for StorageFiles list
	FileMultiSelect map[int]bool // selected rows in StorageFiles
	// file action page
	ActiveFile    string
	ActiveFiles   []string
	FileOpLoading bool
	FileOpErr     error
	FileOpSuccess string
	// edit server page
	EditingServer string // original name of server being edited
	// clean all confirmation page
	CleanInput     textinput.Model
	CleanOpLoading bool
	CleanOpErr     error
	// send / file picker page
	Picker picker
}

func NewTUIInterface(store *services.ServicesStore, localDir string) TUIInterface {
	return TUIInterface{
		Services:        store,
		Page:            pageHome,
		MenuItems:       pages.HomeMenuItems(),
		LocalDir:        localDir,
		FileMultiSelect: make(map[int]bool),
	}
}

// fileListHeight calculates how many rows can fit in the server storage file list.
func (m TUIInterface) fileListHeight() int {
	h := m.WindowHeight
	if h == 0 {
		return 0 // unconstrained until first WindowSizeMsg
	}

	// Card chrome overhead (rounded border + compact inner vertical padding)
	const cardOverhead = 2
	// Header (title + subtitle)
	const headerLines = 2
	// Footer (top border + content)
	const footerLines = 2
	// Server actions menu rows (Send / Clean All)
	actionLines := len(m.MenuItems)
	if actionLines < 1 {
		actionLines = 2
	}
	// File section chrome: top border + local-dir label
	const fileSectionOverhead = 2

	used := cardOverhead + headerLines + footerLines + actionLines + fileSectionOverhead
	available := h - used
	if available < 1 {
		available = 1
	}
	return available
}
