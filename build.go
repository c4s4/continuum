/*
 * Build management stuff. This is used to build modules.
 */

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

// Build is the result of a build.
type Build struct {
	Module   ModuleConfig
	Previous bool
	Success  bool
	Skipped  bool
	Output   string
}

// Builds is a list of builds of the configuration.
type Builds []Build

func (build Build) String() string {
	if build.Skipped {
		return "SKIPPED"
	}
	if build.Success {
		return "SUCCESS"
	}
	return "FAILURE"
}

// SendMail tells if should send an email:
// - config: email configuration.
// Return a bool that tells if should send email.
func (build Build) SendEmail(config EmailConfig) bool {
	if build.Skipped {
		return false
	}
	if config.Once {
		return build.Success != build.Previous
	}
	return !build.Success || (build.Success && config.Success)
}

// BuildModule is called to build a module, that is:
// - get the repository clone.
// - run command to build the module.
// If build command returns 0 (as of Unix standard), the build is a success, else
// this is a failure.
func BuildModule(module ModuleConfig, directory string) Build {
	moduleDir := path.Join(directory, module.Name)
	// go in build directory
	currentDir, err := os.Getwd()
	defer os.Chdir(currentDir)
	err = os.Chdir(directory)
	if err != nil {
		return Build{
			Module:  module,
			Success: false,
			Skipped: false,
			Output:  err.Error(),
		}
	}
	// delete module directory if it already exists
	if _, err := os.Stat(moduleDir); err == nil {
		os.RemoveAll(moduleDir)
	}
	// get the module
	output, err := CloneRepo(module)
	if err != nil {
		return Build{
			Module:  module,
			Success: false,
			Skipped: false,
			Output:  string(output),
		}
	}
	defer os.RemoveAll(moduleDir)
	os.Chdir(moduleDir)
	// run the build command
	cmd := exec.Command("bash", "-c", module.Command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return Build{
			Module:  module,
			Success: false,
			Skipped: false,
			Output:  strings.TrimSpace(string(out)),
		}
	}
	return Build{
		Module:  module,
		Success: true,
		Skipped: false,
		Output:  string(out),
	}
}

// BuildModules builds the list of modules in the configuration (in the exact same
// order).
func BuildModules(config Config) Builds {
	builds := make(Builds, len(config.Modules))
	modulesInfo := LoadModulesInfo(config.Status)
	for _, module := range config.Modules {
		fmt.Printf("Building '%s'... ", module.Name)
		start := time.Now()
		repoHash := GetRepoHash(module)
		var build Build
		if modulesInfo[module.Name].RepoHash == "" {
			build = BuildModule(module, config.Directory)
			build.Previous = true
		} else if modulesInfo[module.Name].RepoHash != repoHash {
			build = BuildModule(module, config.Directory)
			build.Previous = modulesInfo[module.Name].BuildOK
		} else {
			build = Build{
				Module:   module,
				Previous: modulesInfo[module.Name].BuildOK,
				Success:  modulesInfo[module.Name].BuildOK,
				Skipped:  true,
				Output:   "",
			}
		}
		duration := time.Since(start)
		fmt.Println(build.String())
		SendEmail(build, start, duration, config.Email)
		if build.Success {
			modulesInfo[module.Name] = ModuleInfo{
				RepoHash: repoHash,
				BuildOK:  true,
			}
		} else {
			modulesInfo[module.Name] = ModuleInfo{
				RepoHash: repoHash,
				BuildOK:  false,
			}
		}
	}
	SaveModulesInfo(modulesInfo, config.Status)
	return builds
}
