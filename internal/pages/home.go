package pages

type HomePageMsg struct{}

type MenuItem struct {
	Label           string
	Key             string
	RequiresServers bool
}

func HomeMenuItems() []MenuItem {
	return []MenuItem{
		{Label: "Select Server", Key: "server", RequiresServers: true},
		{Label: "Config", Key: "config"},
		{Label: "Exit", Key: "exit"},
	}
}
