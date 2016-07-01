package main

import (
	"io/ioutil"
	"os"
	"testing"
)

const testConfigFile = "/tmp/test-config.yml"

func TestLoadConfig(t *testing.T) {
	config := `directory: /tmp
status:    /tmp/repo-hash.yml
email:
  smtp-host: smtp.example.com:25
  recipient: nobody@nowhere.com
  sender:    nobody@nowhere.com
  success:   true
  once:      true
modules:
- name:    continuum
  url:     git@github.com:c4s4/continuum.git
  command: |
    set -e
    make test`
	ioutil.WriteFile(testConfigFile, []byte(config), 0666)
	defer os.Remove(testConfigFile)
	expected := Config{
		Directory: "/tmp",
		Status:    "/tmp/repo-hash.yml",
		Email: EmailConfig{
			SmtpHost:  "smtp.example.com:25",
			Recipient: "nobody@nowhere.com",
			Sender:    "nobody@nowhere.com",
			Success:   true,
			Once:      true,
		},
		Modules: []ModuleConfig{
			ModuleConfig{
				Name:    "continuum",
				Url:     "git@github.com:c4s4/continuum.git",
				Command: "set -e\nmake test",
			},
		},
	}
	actual := LoadConfig(testConfigFile)
	if expected.Directory != actual.Directory ||
		expected.Email != actual.Email ||
		expected.Status != actual.Status {
		t.Error("Broken configuration loader")
	}
	if len(expected.Modules) != len(actual.Modules) {
		t.Error("Broken configuration loader")
	}
	for i := 0; i < len(expected.Modules); i++ {
		if expected.Modules[i] != actual.Modules[i] {
			t.Error("Broken configuration loader")
		}
	}
}
