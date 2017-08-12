/*
 * Repository management stuff. Reads and writes modules info file that holds
 * information on each module. This file is in YAML format and looks like:
 *
 *   module1:
 *     repo-hash: dbe955d1d83ea4ec969656d1e002e25ca1382fd8
 *     build-ok:  true
 *   module2:
 *     repo-hash: c634c54781a89253167076ce102e588af8a60141
 *     build-ok: false
 *
 * To get this hash, we run command:
 *
 *   $ git ls-remote https://github.com/c4s4/continuum.git
 *   dbe955d1d83ea4ec969656d1e002e25ca1382fd8	HEAD
 *   dbe955d1d83ea4ec969656d1e002e25ca1382fd8	refs/heads/master
 *
 * We take hash for given branch.
 */

package main

import (
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"
)

// ModuleInfo stores information about a given module
type ModuleInfo struct {
	RepoHash string `yaml:"repo-hash"`
	BuildOK  bool   `yaml:"build-ok"`
}

// ModulesInfo is a map that stores info on a given module
type ModulesInfo map[string]ModuleInfo

// CloneRepo clones a given module repository in current directory.
func CloneRepo(module ModuleConfig) (string, error) {
	cmd := exec.Command("git", "clone", "-b", module.Branch, module.Url, module.Name)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// GetRepoHash return the hash of the branch for a given repository.
func GetRepoHash(module ModuleConfig) string {
	cmd := exec.Command("git", "ls-remote", module.Url)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		re, _ := regexp.Compile("(\\w+)\\s+/?refs/heads/" + module.Branch)
		match := re.FindStringSubmatch(line)
		if len(match) > 0 {
			return match[1]
		}
	}
	return ""
}

// LoadModulesInfo loads the modules info in a given file.
func LoadModulesInfo(file string) ModulesInfo {
	if file != "" {
		modulesInfo := ModulesInfo{}
		text, _ := ioutil.ReadFile(file)
		yaml.Unmarshal(text, &modulesInfo)
		return modulesInfo
	}
	return make(ModulesInfo)
}

// SaveModulesInfo saves modules info in a given file.
func SaveModulesInfo(modulesInfo ModulesInfo, file string) {
	contents, err := yaml.Marshal(modulesInfo)
	if err != nil {
		panic("Error writing modules info file: " + err.Error())
	}
	if err := ioutil.WriteFile(file, contents, 0644); err != nil {
		panic("Error writing modules info file: " + err.Error())
	}
}
