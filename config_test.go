package main

import (
	"io/ioutil"
	"os"
	"testing"
)

const testConfigFile = "/tmp/test-config.yml"

func TestLoadConfig(t *testing.T) {
	config := `directory:   /tmp
repo_hash:   /tmp/repo-hash.yml
email:
  smtp_host: smtp.orange.fr:25
  recipient: casa@sweetohm.net
  sender:    casa@sweetohm.net
  success:   true
modules:
- name:    module1
  url:     https://repository/url/module1.git
  command: command to run tests
- name:    module2
  url:     https://repository/url/module2.git
  command: command to run tests`
	ioutil.WriteFile(testConfigFile, []byte(config), 0666)
	defer os.Remove(testConfigFile)
	expected := Config{
		Directory: "/tmp",
		RepoHash:  "/tmp/repo-hash.yml",
		Email: EmailConfig{
			SmtpHost:  "smtp.orange.fr:25",
			Recipient: "casa@sweetohm.net",
			Sender:    "casa@sweetohm.net",
			Success:   true,
		},
		Modules: []ModuleConfig{
			ModuleConfig{
				Name:    "module1",
				Url:     "https://repository/url/module1.git",
				Command: "command to run tests",
			},
			ModuleConfig{
				Name:    "module2",
				Url:     "https://repository/url/module2.git",
				Command: "command to run tests",
			},
		},
	}
	actual := LoadConfig(testConfigFile)
	if expected.Directory != actual.Directory ||
		expected.Email != actual.Email ||
		expected.RepoHash != actual.RepoHash {
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
