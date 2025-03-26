package main

import (
	"os"
	"testing"

	"github.com/dayvillefire/pocsag-monitor/config"
	"gopkg.in/yaml.v2"
)

func Test_Config(t *testing.T) {
	d := config.DynamicConfig{}
	data, err := os.ReadFile("dynamic.yaml")
	if err != nil {
		t.Error(err)
	}
	err = yaml.Unmarshal([]byte(data), &d)
	if err != nil {
		t.Error(err)
	}
}

func Test_LoadConfig(t *testing.T) {
	c, err := config.LoadConfigWithDefaults("config.yaml", "dynamic.yaml")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	t.Logf("%#v\n", c)
}
