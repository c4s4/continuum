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

type Config struct {
	Directory string
	Email     struct {
		SmtpHost  string "smtp_host"
		Recipient string
		Sender    string
		Success   bool
	}
	Modules []struct {
		Name    string
		Url     string
		Command string
	}
}

type Build struct {
	Module  string
	Success bool
	Output  string
}

func (build Build) String() string {
	if build.Success {
		return fmt.Sprintf("%s: OK", build.Module)
	} else {
		return fmt.Sprintf("%s: ERROR", build.Module)
	}
}

type Builds []Build

func (builds Builds) Success() bool {
	for _, module := range builds {
		if !module.Success {
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

func buildModule(index int, config Config) Build {
	fmt.Printf("Building '%s'... ", config.Modules[index].Name)
	moduleDir := path.Join(config.Directory, config.Modules[index].Name)
	// go in build directory
	err := os.Chdir(config.Directory)
	if err != nil {
		return Build{
			Module:  config.Modules[index].Name,
			Success: false,
			Output:  err.Error(),
		}
	}
	// delete module directory if it already exists
	if _, err := os.Stat(moduleDir); err == nil {
		os.RemoveAll(moduleDir)
	}
	// git clone the module repository
	cmd := exec.Command("git", "clone", config.Modules[index].Url)
	output, err := cmd.CombinedOutput()
	defer os.RemoveAll(moduleDir)
	if err != nil {
		fmt.Println("ERROR")
		return Build{
			Module:  config.Modules[index].Name,
			Success: false,
			Output:  string(output),
		}
	} else {
		os.Chdir(moduleDir)
		// run the build command
		cmd := exec.Command("bash", "-c", config.Modules[index].Command)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("ERROR")
			return Build{
				Module:  config.Modules[index].Name,
				Success: false,
				Output:  strings.TrimSpace(string(output)),
			}
		} else {
			fmt.Println("OK")
			return Build{
				Module:  config.Modules[index].Name,
				Success: true,
				Output:  string(output),
			}
		}
	}
}

func buildModules(config Config) Builds {
	builds := make(Builds, len(config.Modules))
	for index, _ := range config.Modules {
		builds[index] = buildModule(index, config)
	}
	return builds
}

const timeFormat = "2006-01-02 15:04"

func sendReport(builds Builds, start time.Time, duration time.Duration, config Config) {
	if !builds.Success() || (builds.Success() && config.Email.Success) {
		subject := fmt.Sprintf("Build on %s was a %s", start.Format(timeFormat), builds)
		message := fmt.Sprintf("From: %s\n", config.Email.Sender)
		message += fmt.Sprintf("To: %s\n", config.Email.Recipient)
		message += fmt.Sprintf("Subject: %s\n\n", subject)
		message += fmt.Sprintf("Build on %s:\n\n", start.Format(timeFormat))
		for module := range builds {
			message += fmt.Sprintf("  %s\n", builds[module].String())
		}
		message += fmt.Sprintf("\nDone in %s\n", duration)
		message += builds.String()
		for _, build := range builds {
			if !build.Success {
				message += fmt.Sprintf("\n\n===================================\n")
				message += fmt.Sprintf(build.Module)
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
		sendReport(builds, start, duration, config)
	}
}
