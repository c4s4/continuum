package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

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

func GetModule(module ModuleConfig) (string, error) {
	cmd := exec.Command("git", "clone", module.Url)
	output, err := cmd.CombinedOutput()
	return string(output), err
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
	// git clone the module repository
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
	for index, module := range config.Modules {
		builds[index] = BuildModule(module, config.Directory)
	}
	return builds
}
