/*
 * Configuration management stuff. This configuration is built from YAML
 * configuration file:
 *
 *   directory:   /tmp
 *   repo_hash:   /tmp/repo-status.yml
 *   port:        6666
 *   email:
 *     smtp_host: smtp.example.com:25
 *     recipient: nobody@nowhere.com
 *     sender:    nobody@nowhere.com
 *     success:   true
 *     once:      true
 *   modules:
 *   - name:    continuum
 *     url:     git@github.com:c4s4/continuum.git
 *     branch:  develop
 *     command: |
 *       set -e
 *       make test
 */

package main

import (
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"path/filepath"
)

// EmailConfig is the configuration to send email report.
type EmailConfig struct {
	SmtpHost  string "smtp-host"
	Recipient string
	Sender    string
	Success   bool
	Once      bool
}

// ModuleConfig is the configuration for a given module.
type ModuleConfig struct {
	Name    string
	Url     string
	Branch  string
	Command string
}

// Config is the global configuration of the application.
type Config struct {
	Directory string
	Status    string
	Port      int
	Email     EmailConfig
	Modules   []ModuleConfig
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
