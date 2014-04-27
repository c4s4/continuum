package main

/*
 * Main file for Gontinuum.
 */

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"time"
)

const Help = `Usage: goontinuum [configuration.yml]
Where configuration.yml is as follows:

  directory:   /tmp
  repo_hash:   /tmp/repo-hash.yml
  port:        6666
  email:
    smtp_host: smtp.orange.fr:25
    recipient: casa@sweetohm.net
    sender:    casa@sweetohm.net
    success:   true
  modules:
  - name:    module1
    url:     https://repository/url/module1.git
    command: |
      command to run tests
  - name:    module2
    url:     https://repository/url/module2.git
    command: |
      command to run tests

If configuration file is not passed on command line, it will be searched at
following locations:
- ~/.gontinuum.yml
- ~/etc/gontinuum.yml
- /etc/gontinuum.yml
`

// FileExists tells if a given file exists.
func FileExists(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	} else {
		return false
	}
}

// FindConfiguration looks for a configuration file as:
// - ~/.gontinuum.yml
// - ~/etc/gontinuum.yml
// - /etc/gontinuum.yml
// If none of these is found, it stops the program and prints help, else it returns
// the path of the configuration file.
func FindConfiguration() string {
	usr, _ := user.Current()
	home := usr.HomeDir
	config := path.Join(home, ".gontinuum.yml")
	if FileExists(config) {
		return config
	}
	config = path.Join(home, "etc", "gontinuum.yml")
	if FileExists(config) {
		return config
	}
	config = "/etc/gontinuum.yml"
	if FileExists(config) {
		return config
	}
	fmt.Println("ERROR: no configuration file found")
	fmt.Println(Help)
	os.Exit(1)
	return ""
}

// CheckArguments checks arguments passed on command line and prints help if
// -h or if an error occurs. Return the path of the configuration file.
func CheckArguments() string {
	if len(os.Args) > 2 {
		fmt.Println("ERROR: You must pass configuration file on command line")
		fmt.Println(Help)
		os.Exit(1)
	}
	if len(os.Args) == 2 {
		if os.Args[1] == "-h" || os.Args[1] == "--help" {
			fmt.Println(Help)
			os.Exit(0)
			return ""
		} else {
			return os.Args[1]
		}
	} else {
		return FindConfiguration()
	}
}

// main function that iterate on configuration files passed on command line.
func main() {
	configFile := CheckArguments()
	config := LoadConfig(configFile)
	if IsAnotherInstanceRunning(config.Port) {
		fmt.Println("Another instance is already running, aborting")
		os.Exit(0)
	} else {
		start := time.Now()
		builds := BuildModules(config)
		duration := time.Since(start)
		fmt.Println("Done in", duration)
		fmt.Println(builds.String())
		SendEmail(builds, start, duration, config.Email)
	}
}
