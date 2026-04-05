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
)

type TUIInterface struct {
	Services       *services.ServicesStore
	Page           page
	MenuItems      []pages.MenuItem
	Selected       int
	Servers        map[string]services.Server
	ServerNames    []string // sorted, stable order for list rendering
	NoServers      bool
	InitErr        error
	FlashMsg       string
	Form           addServerForm
	FormErr        string // inline field error (e.g. duplicate name)
	Quitting       bool
	WindowWidth    int
	WindowHeight   int
}

func NewTUIInterface(store *services.ServicesStore) TUIInterface {
	return TUIInterface{
		Services:  store,
		Page:      pageHome,
		MenuItems: pages.HomeMenuItems(),
	}
}
