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
