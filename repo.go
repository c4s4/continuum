/*
 * Repository management stuff. Reads and writes a repository hash file that
 * holds the SHA1 for each module. This file is in YAML format and looks like:
 *
 *   module1: dbe955d1d83ea4ec969656d1e002e25ca1382fd8
 *   module2: c634c54781a89253167076ce102e588af8a60141
 *
 * To get this hash, we run command:
 *
 *   $ git ls-remote https://github.com/c4s4/gontinuum.git
 *   dbe955d1d83ea4ec969656d1e002e25ca1382fd8	HEAD
 *   dbe955d1d83ea4ec969656d1e002e25ca1382fd8	refs/heads/master
 *
 * We take hash for HEAD.
 */

package main

import (
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"
)

// RepoHashMap is a map that gives hash of the HEAD for a given repo.
type RepoHashMap map[string]string

// CloneRepo clones a given module repository in current directory.
func CloneRepo(module ModuleConfig) (string, error) {
	cmd := exec.Command("git", "clone", "-b", module.Branch, module.Url)
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
		re, _ := regexp.Compile("(\\w+)\\s+/refs/heads/" + module.Branch)
		match := re.FindStringSubmatch(line)
		if len(match) > 0 {
			return match[1]
		}
	}
	return ""
}

// LoadRepoHashMap loads the repo hash map in a given file.
func LoadRepoHashMap(file string) RepoHashMap {
	if file != "" {
		repoStatus := RepoHashMap{}
		text, _ := ioutil.ReadFile(file)
		yaml.Unmarshal(text, &repoStatus)
		return repoStatus
	} else {
		return make(RepoHashMap)
	}
}

// SaveRepoHashMap saves the repo hash in a given file.
func SaveRepoHashMap(repoHashMap RepoHashMap, file string) {
	contents, err := yaml.Marshal(repoHashMap)
	if err != nil {
		panic("Error writing repo hash map file: " + err.Error())
	}
	if err := ioutil.WriteFile(file, contents, 0644); err != nil {
		panic("Error writing repo hash map file: " + err.Error())
	}
}
