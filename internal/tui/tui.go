package tui

import (
	"filepass/internal/pages"
	"filepass/internal/services"
)

type page int

const (
	pageHome page = iota
	pageConfig
	pageAddServer
	pageSelectServer
	pageServerActions
	pageFileAction
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
	StorageFiles   []string
	StorageLoading bool
	StorageErr     error
	FileSelected   int  // cursor within StorageFiles
	FileFocused    bool // true = ↑↓ drives file list, false = action menu
	// file action page
	ActiveFile string
}

func NewTUIInterface(store *services.ServicesStore) TUIInterface {
	return TUIInterface{
		Services:  store,
		Page:      pageHome,
		MenuItems: pages.HomeMenuItems(),
	}
}
