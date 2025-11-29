package main

import (
	"testing"

	"github.com/dayvillefire/pocsag-monitor/config"
)

func Test_LoadConfig(t *testing.T) {
	c, err := config.LoadConfigWithDefaults("config.yaml", "dynamic.yaml")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	t.Logf("%#v\n", c)
}
