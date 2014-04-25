package main

/*
 * Configuration management stuff. This configuration is built from YAML
 * configuration file:
 *
 *   directory:   /tmp
 *   repo_status: /tmp/repo-status.yml
 *   email:
 *     smtp_host: smtp.orange.fr:25
 *     recipient: casa@sweetohm.net
 *     sender:    casa@sweetohm.net
 *     success:   true
 *   modules:
 *   - name:    module1
 *     url:     https://repository/url/module1.git
 *     command: |
 *       command to run tests
 *   - name:    module2
 *     url:     https://repository/url/module2.git
 *     command: |
 *       command to run tests
 */

import (
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"path/filepath"
)

// EmailConfig is the configuration to send email report.
type EmailConfig struct {
	SmtpHost  string "smtp_host"
	Recipient string
	Sender    string
	Success   bool
}

// ModuleConfig is the configuration for a given module.
type ModuleConfig struct {
	Name    string
	Url     string
	Command string
}

// Config is the global configuration of the application.
type Config struct {
	Directory  string
	RepoStatus string "repo_status"
	Email      EmailConfig
	Modules    []ModuleConfig
}

// LoadConfig loads configuration from a given YAML configuration file.
func LoadConfig(file string) Config {
	config := Config{}
	text, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err.Error())
	}
	err = yaml.Unmarshal(text, &config)
	if err != nil {
		panic(err.Error())
	}
	// make directory absolute path
	absdir, err := filepath.Abs(config.Directory)
	if err != nil {
		panic(err.Error())
	}
	config.Directory = absdir
	return config
}
