package pages

type ServerActionsPageMsg struct {
	ServerName string
}

func ServerActionItems() []MenuItem {
	return []MenuItem{
		{Label: "Send", Key: "send"},
		{Label: "Get", Key: "get"},
		{Label: "Clean", Key: "clean"},
	}
}
