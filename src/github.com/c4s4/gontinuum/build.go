package main

import (
	"fmt"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

type Build struct {
	Module  ModuleConfig
	Success bool
	Output  string
}

type Builds []Build

type RepoStatus map[string]string

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

func BuildModule(module ModuleConfig, directory string) Build {
	fmt.Printf("Building '%s'... ", module.Name)
	moduleDir := path.Join(directory, module.Name)
	// go in build directory
	currentDir, err := os.Getwd()
	defer os.Chdir(currentDir)
	err = os.Chdir(directory)
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
	// get the module
	output, err := GetModule(module)
	if err != nil {
		fmt.Println("ERROR")
		return Build{
			Module:  module,
			Success: false,
			Output:  string(output),
		}
	} else {
		defer os.RemoveAll(moduleDir)
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

func BuildModules(config Config) Builds {
	builds := make(Builds, len(config.Modules))
	repoStatus := LoadRepoStatus(config.RepoStatus)
	for index, module := range config.Modules {
		if repoStatus[module.Name] == "" ||
			(repoStatus[module.Name] != GetRepoStatus(module)) {
			builds[index] = BuildModule(module, config.Directory)
		}
	}
	return builds
}

func GetModule(module ModuleConfig) (string, error) {
	cmd := exec.Command("git", "clone", module.Url)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func GetRepoStatus(module ModuleConfig) string {
	cmd := exec.Command("git", "ls-remote", module.Url)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		re, _ := regexp.Compile("(\\w+)\\s+HEAD")
		match := re.FindStringSubmatch(line)
		if len(match) > 0 {
			return match[1]
		}
	}
	return ""
}

func LoadRepoStatus(file string) RepoStatus {
	if file != "" {
		repoStatus := RepoStatus{}
		text, _ := ioutil.ReadFile(file)
		yaml.Unmarshal(text, &repoStatus)
		return repoStatus
	} else {
		return make(RepoStatus)
	}
}
