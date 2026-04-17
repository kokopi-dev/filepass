package pages

type FileActionPageMsg struct {
	ServerName string
	Filename   string   // single-file mode
	Filenames  []string // multi-file mode
}

func FileActionItems() []MenuItem {
	return []MenuItem{
		{Label: "Get", Key: "get"},
		{Label: "Delete", Key: "delete"},
	}
}
