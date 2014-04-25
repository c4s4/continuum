package main

/*
 * Main file for Gontinuum. Configuration file is in YAML format with following
 * format:
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
	"fmt"
	"os"
	"time"
)

// main function that iterate on configuration files passed on command line.
func main() {
	for i := 1; i < len(os.Args); i++ {
		start := time.Now()
		config := LoadConfig(os.Args[i])
		builds := BuildModules(config)
		duration := time.Since(start)
		fmt.Println("Done in", duration)
		fmt.Println(builds.String())
		SendEmail(builds, start, duration, config.Email)
	}
}
