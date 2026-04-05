package tui

import (
	"strings"

	"filepass/internal/services"

	"charm.land/bubbles/v2/textinput"
)

const (
	fieldName = iota
	fieldHost
	fieldUser
	fieldPrivateKey
	fieldPort
	fieldSave
	fieldBack
	fieldCount
)

type addServerForm struct {
	inputs  [5]textinput.Model
	focused int // 0–6: inputs 0-4, save 5, back 6
}

func newAddServerForm() addServerForm {
	mkInput := func(placeholder string, limit int) textinput.Model {
		ti := textinput.New()
		ti.Prompt = ""
		ti.CharLimit = limit
		ti.SetWidth(40)
		ti.Placeholder = placeholder
		return ti
	}

	f := addServerForm{}
	f.inputs[fieldName] = mkInput("production-web", 64)
	f.inputs[fieldHost] = mkInput("192.168.1.1 or example.com", 253)
	f.inputs[fieldUser] = mkInput("deploy", 64)
	f.inputs[fieldPrivateKey] = mkInput("~/.ssh/id_rsa", 512)
	f.inputs[fieldPort] = mkInput("22  (optional)", 5)
	f.inputs[fieldName].Focus()
	return f
}

func (f *addServerForm) focusField(i int) {
	for j := range f.inputs {
		f.inputs[j].Blur()
	}
	if i < len(f.inputs) {
		f.inputs[i].Focus()
	}
	f.focused = i
}

func (f addServerForm) canSave() bool {
	return strings.TrimSpace(f.inputs[fieldName].Value()) != "" &&
		strings.TrimSpace(f.inputs[fieldHost].Value()) != "" &&
		strings.TrimSpace(f.inputs[fieldUser].Value()) != "" &&
		strings.TrimSpace(f.inputs[fieldPrivateKey].Value()) != ""
}

func (f addServerForm) focusNext() addServerForm {
	next := f.focused + 1
	if next >= fieldCount {
		next = fieldCount - 1
	}
	f.focusField(next)
	return f
}

func (f addServerForm) focusPrev() addServerForm {
	prev := f.focused - 1
	if prev < 0 {
		prev = 0
	}
	f.focusField(prev)
	return f
}

// newEditServerForm builds a pre-filled form for editing an existing server.
// The name field is pre-populated with the server key; all other fields with
// the existing server values.
func newEditServerForm(name string, s services.Server) addServerForm {
	mkInput := func(placeholder string, limit int, value string) textinput.Model {
		ti := textinput.New()
		ti.Prompt = ""
		ti.CharLimit = limit
		ti.SetWidth(40)
		ti.Placeholder = placeholder
		ti.SetValue(value)
		return ti
	}

	port := s.Port
	if port == "22" {
		port = ""
	}

	f := addServerForm{}
	f.inputs[fieldName] = mkInput("production-web", 64, name)
	f.inputs[fieldHost] = mkInput("192.168.1.1 or example.com", 253, s.Host)
	f.inputs[fieldUser] = mkInput("deploy", 64, s.User)
	f.inputs[fieldPrivateKey] = mkInput("~/.ssh/id_rsa", 512, s.PrivateKey)
	f.inputs[fieldPort] = mkInput("22  (optional)", 5, port)
	f.inputs[fieldName].Focus()
	return f
}
