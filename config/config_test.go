package config

import "testing"

func Test_parsfile(t *testing.T) {
	cfg := New()
	err := cfg.Init("config.toml")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(*cfg)
}
