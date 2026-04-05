package tui

import (
	"os"
	"path/filepath"
	"strings"
)

// entry is a single item in the file picker list.
type entry struct {
	name  string
	isDir bool
}

// picker is the state for the send file picker page.
type picker struct {
	dir      string    // current directory being browsed
	entries  []entry   // unfiltered entries in dir
	filtered []entry   // entries matching query
	query    string    // current filter string
	cursor   int       // index within filtered
}

func newPicker(startDir string) picker {
	p := picker{dir: startDir}
	p.entries = readDir(startDir)
	p.filtered = p.entries
	return p
}

// readDir lists the entries of a directory, dirs first then files.
func readDir(dir string) []entry {
	infos, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var dirs, files []entry
	for _, d := range infos {
		name := d.Name()
		if strings.HasPrefix(name, ".") {
			continue // skip hidden
		}
		if d.IsDir() {
			dirs = append(dirs, entry{name: name + "/", isDir: true})
		} else {
			files = append(files, entry{name: name, isDir: false})
		}
	}
	return append(dirs, files...)
}

// applyFilter rebuilds filtered from entries using query.
func (p picker) applyFilter() picker {
	if p.query == "" {
		p.filtered = p.entries
	} else {
		q := strings.ToLower(p.query)
		var out []entry
		for _, e := range p.entries {
			if strings.Contains(strings.ToLower(e.name), q) {
				out = append(out, e)
			}
		}
		p.filtered = out
	}
	p.cursor = 0
	return p
}

// descend enters a subdirectory.
func (p picker) descend(name string) picker {
	// strip trailing slash added for display
	name = strings.TrimSuffix(name, "/")
	next := filepath.Join(p.dir, name)
	p.dir = next
	p.entries = readDir(next)
	p.query = ""
	p.filtered = p.entries
	p.cursor = 0
	return p
}

// ascend goes up one directory level.
func (p picker) ascend() picker {
	parent := filepath.Dir(p.dir)
	if parent == p.dir {
		return p // already at root
	}
	p.dir = parent
	p.entries = readDir(parent)
	p.query = ""
	p.filtered = p.entries
	p.cursor = 0
	return p
}

// selectedPath returns the full path of the currently highlighted entry,
// or empty string if the list is empty.
func (p picker) selectedPath() string {
	if len(p.filtered) == 0 || p.cursor < 0 || p.cursor >= len(p.filtered) {
		return ""
	}
	e := p.filtered[p.cursor]
	name := strings.TrimSuffix(e.name, "/")
	return filepath.Join(p.dir, name)
}

// typeRune appends a rune to the query and re-filters.
func (p picker) typeRune(r rune) picker {
	p.query += string(r)
	return p.applyFilter()
}

// backspace removes the last rune from the query.
// If query is already empty, ascend instead.
func (p picker) backspace() (picker, bool) {
	if p.query == "" {
		return p.ascend(), false // false = did not consume (went up)
	}
	runes := []rune(p.query)
	p.query = string(runes[:len(runes)-1])
	return p.applyFilter(), true
}
