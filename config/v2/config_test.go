package config

import "testing"

func TestParse(t *testing.T) {
	data := `
webauthn: true
mail: 127.0.0.1`

	model, err := ParseConfigContent(data)
	if err != nil {
		t.Errorf("parse config content error: %s", err)
		return
	}
	value, err := model.GetBool("webauthn")
	if err != nil {
		t.Errorf("get string error: %s", err)
		return
	}
	if !value {
		t.Errorf("get string error: %v", value)
		return
	}

	mailHost, err := model.GetString("mail")
	if err != nil {
		t.Errorf("get string error2: %s", err)
		return
	}
	if mailHost != "127.0.0.1" {
		t.Errorf("wrong mailHost: %v", mailHost)
		return
	}
}
