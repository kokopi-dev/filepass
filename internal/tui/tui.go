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
	ActiveServer   string
	LocalDir       string // user's cwd, destination for received files
	StorageFiles   []string
	StorageLoading bool
	StorageErr     error
	FileSelected   int  // cursor within StorageFiles
	FileFocused    bool // true = ↑↓ drives file list, false = action menu
	// file action page
	ActiveFile    string
	FileOpLoading bool
	FileOpErr     error
	FileOpSuccess string
	// clean all confirmation page
	CleanInput    textinput.Model
	CleanOpLoading bool
	CleanOpErr    error
	// send / file picker page
	Picker picker
}

func NewTUIInterface(store *services.ServicesStore, localDir string) TUIInterface {
	return TUIInterface{
		Services:  store,
		Page:      pageHome,
		MenuItems: pages.HomeMenuItems(),
		LocalDir:  localDir,
	}
}
