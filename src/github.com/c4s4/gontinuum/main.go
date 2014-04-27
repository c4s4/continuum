package main

/*
 * Main file for Gontinuum.
 */

import (
	"fmt"
	"os"
	"time"
)

const Help = `Usage: goontinuum configuration.yml
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
`

// CheckArguments checks arguments passed on command line and prints help if
// -h or if an error occurs.
func CheckArguments() {
	if len(os.Args) != 2 {
		fmt.Println("ERROR: You must pass configuration file on command line")
		fmt.Println(Help)
		os.Exit(1)
	} else if os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Println(Help)
		os.Exit(0)
	}
}

// main function that iterate on configuration files passed on command line.
func main() {
	CheckArguments()
	config := LoadConfig(os.Args[1])
	Highlander(config.Port)
	start := time.Now()
	builds := BuildModules(config)
	duration := time.Since(start)
	fmt.Println("Done in", duration)
	fmt.Println(builds.String())
	SendEmail(builds, start, duration, config.Email)
}
