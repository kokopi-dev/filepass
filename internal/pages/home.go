package pages

type HomePageMsg struct{}

type MenuItem struct {
	Label string
	Key   string
}

func HomeMenuItems() []MenuItem {
	return []MenuItem{
		{Label: "Select Server", Key: "server"},
		{Label: "Config", Key: "config"},
		{Label: "Exit", Key: "exit"},
	}
}
