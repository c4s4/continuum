package main

import (
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"path/filepath"
)

type EmailConfig struct {
	SmtpHost  string "smtp_host"
	Recipient string
	Sender    string
	Success   bool
}

type ModuleConfig struct {
	Name    string
	Url     string
	Command string
}

type Config struct {
	Directory  string
	RepoStatus string "repo_status"
	Email      EmailConfig
	Modules    []ModuleConfig
}

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
