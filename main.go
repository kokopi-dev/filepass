package main

import (
	"fmt"
	"os"

	"filepass/internal/services"
	"filepass/internal/tui"

	tea "charm.land/bubbletea/v2"
)

func main() {
	store, err := services.NewServicesStore()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to initialise config:", err)
		os.Exit(1)
	}

	m := tui.NewTUIInterface(store)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
