package tui

import (
	"filepass/internal/pages"
	"filepass/internal/services"
)

type TUIInterface struct {
	Services     *services.ServicesStore
	MenuItems    []pages.MenuItem
	Selected     int
	Servers      map[string]services.Server
	NoServers    bool
	InitErr      error
	Quitting     bool
	WindowWidth  int
	WindowHeight int
}

func NewTUIInterface(store *services.ServicesStore) TUIInterface {
	return TUIInterface{
		Services:  store,
		MenuItems: pages.HomeMenuItems(),
	}
}
