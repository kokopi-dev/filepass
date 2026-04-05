package pages

type ConfigPageMsg struct{}

func ConfigMenuItems() []MenuItem {
	return []MenuItem{
		{Label: "Add Server", Key: "add"},
		{Label: "Edit Server", Key: "edit", RequiresServers: true},
		{Label: "Remove Server", Key: "remove", RequiresServers: true},
		{Label: "Back", Key: "back"},
	}
}
