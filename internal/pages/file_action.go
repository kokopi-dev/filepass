package pages

type FileActionPageMsg struct {
	ServerName string
	Filename   string
}

func FileActionItems() []MenuItem {
	return []MenuItem{
		{Label: "Get", Key: "get"},
		{Label: "Delete", Key: "delete"},
	}
}
