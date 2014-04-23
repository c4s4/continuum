// Sample configuration file:
//
//   directory:  /home/casa/tmp
//   email:
//     smtp_host: smtp.orange.fr:25
//     recipient: casa@sweetohm.net
//     sender:    casa@sweetohm.net
//     success:   true
//
//   modules:
//     - name:    continuum
//       url:     https://github.com/c4s4/continuum.git
//       command: |
//         set -e
//         export PATH=/opt/python/current/bin:$PATH
//         virtualenv env --no-site-packages
//         . env/bin/activate
//         pip install -r etc/requirements.txt
//         bee test

package main

import (
	"fmt"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"net/smtp"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"
)

///////////////////////////////////////////////////////////////////////////////
//                              CONFIG STUFF                                 //
///////////////////////////////////////////////////////////////////////////////

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
	Directory string
	Email     EmailConfig
	Modules   []ModuleConfig
}

func loadConfig(file string) Config {
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

///////////////////////////////////////////////////////////////////////////////
//                                BUILD STUFF                                //
///////////////////////////////////////////////////////////////////////////////

type Build struct {
	Module  ModuleConfig
	Success bool
	Output  string
}

type Builds []Build

func (build Build) String() string {
	if build.Success {
		return fmt.Sprintf("%s: OK", build.Module.Name)
	} else {
		return fmt.Sprintf("%s: ERROR", build.Module.Name)
	}
}

func (builds Builds) Success() bool {
	for _, build := range builds {
		if !build.Success {
			return false
		}
	}
	return true
}

func (builds Builds) String() string {
	if builds.Success() {
		return "SUCCESS"
	} else {
		return "FAILURE"
	}
}

func buildModule(module ModuleConfig, directory string) Build {
	fmt.Printf("Building '%s'... ", module.Name)
	moduleDir := path.Join(directory, module.Name)
	// go in build directory
	err := os.Chdir(directory)
	if err != nil {
		return Build{
			Module:  module,
			Success: false,
			Output:  err.Error(),
		}
	}
	// delete module directory if it already exists
	if _, err := os.Stat(moduleDir); err == nil {
		os.RemoveAll(moduleDir)
	}
	// git clone the module repository
	cmd := exec.Command("git", "clone", module.Url)
	output, err := cmd.CombinedOutput()
	defer os.RemoveAll(moduleDir)
	if err != nil {
		fmt.Println("ERROR")
		return Build{
			Module:  module,
			Success: false,
			Output:  string(output),
		}
	} else {
		os.Chdir(moduleDir)
		// run the build command
		cmd := exec.Command("bash", "-c", module.Command)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("ERROR")
			return Build{
				Module:  module,
				Success: false,
				Output:  strings.TrimSpace(string(output)),
			}
		} else {
			fmt.Println("OK")
			return Build{
				Module:  module,
				Success: true,
				Output:  string(output),
			}
		}
	}
}

func buildModules(config Config) Builds {
	builds := make(Builds, len(config.Modules))
	for index, module := range config.Modules {
		builds[index] = buildModule(module, config.Directory)
	}
	return builds
}

const timeFormat = "2006-01-02 15:04"

func sendEmail(builds Builds, start time.Time, duration time.Duration, config Config) {
	if !builds.Success() || (builds.Success() && config.Email.Success) {
		subject := fmt.Sprintf("Build on %s was a %s", start.Format(timeFormat), builds)
		message := fmt.Sprintf("From: %s\n", config.Email.Sender)
		message += fmt.Sprintf("To: %s\n", config.Email.Recipient)
		message += fmt.Sprintf("Subject: %s\n\n", subject)
		message += fmt.Sprintf("Build on %s:\n\n", start.Format(timeFormat))
		for _, build := range builds {
			message += fmt.Sprintf("  %s\n", build.String())
		}
		message += fmt.Sprintf("\nDone in %s\n", duration)
		message += builds.String()
		for _, build := range builds {
			if !build.Success {
				message += fmt.Sprintf("\n\n===================================\n")
				message += fmt.Sprintf(build.Module.Name)
				message += fmt.Sprintf("\n-----------------------------------\n")
				message += fmt.Sprintf(build.Output)
				message += fmt.Sprintf("\n-----------------------------------\n")
			}
		}
		message += "\n--\ngontinuum"
		err := smtp.SendMail(config.Email.SmtpHost, nil, config.Email.Sender,
			[]string{config.Email.Recipient}, []byte(message))
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	for i := 1; i < len(os.Args); i++ {
		start := time.Now()
		config := loadConfig(os.Args[i])
		builds := buildModules(config)
		duration := time.Since(start)
		fmt.Println("Done in", duration)
		fmt.Println(builds.String())
		sendEmail(builds, start, duration, config)
	}
}
