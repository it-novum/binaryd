package config

import (
	"fmt"

	"gopkg.in/ini.v1"
)

type Config struct {
	Commands []CommandCfg
	Path     string
	cfg      *ini.File
}

type CommandCfg struct {
	CommandName string
	CommandLine string
}

func NewConfig(path string) *Config {
	return &Config{
		Path: path,
	}
}

func (c *Config) LoadIni() error {
	cfg, err := ini.Load(c.Path)
	if err != nil {
		return err
	}

	c.cfg = cfg

	return nil
}

func (c *Config) ParseIni() error {
	sections := c.cfg.SectionStrings()

	for _, section := range sections {
		if section != "DEFAULT" {
			val := c.cfg.Section(section).Key("command").String()
			if val != "" {
				// Append command from ini to commands array
				c.Commands = append(c.Commands, CommandCfg{
					CommandName: section,
					CommandLine: val,
				})
			}
		}
	}

	if len(c.Commands) == 0 {
		return fmt.Errorf("no commands defined. Please check your config.ini")
	}

	return nil
}
