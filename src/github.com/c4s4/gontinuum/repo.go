package main

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

import (
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"
)

type RepoHash map[string]string

func GetModule(module ModuleConfig) (string, error) {
	cmd := exec.Command("git", "clone", module.Url)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func GetRepoHash(module ModuleConfig) string {
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

func LoadRepoHash(file string) RepoHash {
	if file != "" {
		repoStatus := RepoHash{}
		text, _ := ioutil.ReadFile(file)
		yaml.Unmarshal(text, &repoStatus)
		return repoStatus
	} else {
		return make(RepoHash)
	}
}

func SaveRepoHash(repoHash RepoHash, file string) {
	contents, err := yaml.Marshal(repoHash)
	if err != nil {
		panic("Error writing repo hash file: " + err.Error())
	}
	ioutil.WriteFile(file, contents, 0644)
}
