package main

import (
	"fmt"
	"os"
	"os/exec"

	"filepass/internal/services"
	"filepass/internal/tui"

	tea "charm.land/bubbletea/v2"
)

func main() {
	if _, err := exec.LookPath("rsync"); err != nil {
		fmt.Fprintln(os.Stderr, "error: rsync is required but was not found in PATH")
		fmt.Fprintln(os.Stderr, "install it with your package manager, e.g.:")
		fmt.Fprintln(os.Stderr, "  brew install rsync")
		fmt.Fprintln(os.Stderr, "  apt install rsync")
		os.Exit(1)
	}

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
