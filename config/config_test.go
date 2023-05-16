package config

import (
	"testing"
)

func TestParseMissingFile(t *testing.T) {
	cfg := NewConfig("../binaryd.foobar.123.do.not.create.this.ini")
	err := cfg.LoadIni()
	if err == nil {
		t.Fatal("There has to be an error")
	}
}

func TestLoadExampleConfig(t *testing.T) {
	cfg := NewConfig("../binaryd.example.ini")
	err := cfg.LoadIni()
	if err != nil {
		t.Fatal(err)
	}
}

func TestParseExampleConfig(t *testing.T) {
	cfg := NewConfig("../binaryd.example.ini")
	err := cfg.LoadIni()
	if err != nil {
		t.Fatal(err)
	}

	err = cfg.ParseIni()
	if err != nil {
		t.Fatal(err)
	}

	if len(cfg.Commands) != 2 {
		t.Fatal("I expect exactly 2 commands")
	}

	if cfg.Commands[0].Command != "/usr/bin/ps -eaf" {
		t.Fatal("first command has to be '/usr/bin/ps -eaf'")
	}

	if cfg.Commands[1].Command != "whoami" {
		t.Fatal("first command has to be 'whoami'")
	}

}
