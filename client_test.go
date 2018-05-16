package routerosclient

import "testing"

func TestConfigValidate(t *testing.T) {
	conf := Config{
		Address:  "127.0.0.1:8728a",
		Username: "vagrant",
		Password: "vagrant",
	}

	if err := conf.validate(); err == nil {
		t.Errorf("invalid Config must raise error")
	}
}
