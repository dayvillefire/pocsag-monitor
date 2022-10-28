package main

import (
	"io/ioutil"

	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
)

var (
	config *Config
)

type Config struct {
	DbFile          string            `yaml:"db-file" default:"scan.db"`
	RtlFmBinary     string            `yaml:"rtlfm" default:"rtl_fm"`
	Frequency       string            `yaml:"frequency" default:"152.00750M"`
	MultiMonBinary  string            `yaml:"multimon" default:"multimon-ng"`
	PPM             int               `yaml:"ppm" default:"0"`
	DiscordToken    string            `yaml:"discord-token"`
	DiscordChannels map[string]string `yaml:"discord-channels"`
	ChannelMappings map[string]string `yaml:"channel-mappings"`
}

// UnmarshalYAML overrides default handling
func (s *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	defaults.Set(s)

	type plain Config
	if err := unmarshal((*plain)(s)); err != nil {
		return err
	}

	return nil
}

func LoadConfigWithDefaults(configPath string) (*Config, error) {
	c := &Config{}
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return c, err
	}
	err = yaml.Unmarshal([]byte(data), c)

	return c, err
}
