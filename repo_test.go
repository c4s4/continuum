package main

import (
	"errors"
	"io/ioutil"
	"os"
	"regexp"
	"testing"
)

func TestGetRepoHash(t *testing.T) {
	moduleConfig := ModuleConfig{
		Name:    "continuum",
		Url:     "git@github.com:c4s4/continuum.git",
		Branch:  "master",
		Command: "echo 'TEST'",
	}
	repoHash := GetRepoHash(moduleConfig)
	if match, _ := regexp.MatchString("^[0-9a-f]{40}$", repoHash); match != true {
		t.Errorf("GetRepoStatus() response '%s' doesn't look like a hash", repoHash)
	}
}

const testModulesInfo = `module1:
  repo-hash: dbe955d1d83ea4ec969656d1e002e25ca1382fd8
  build-ok: true
module2:
  repo-hash: c634c54781a89253167076ce102e588af8a60141
  build-ok: false
`

func TestLoadModulesInfo(t *testing.T) {
	tempFile, err := ioutil.TempFile("/tmp", "go-test-")
	if err != nil {
		panic(errors.New("Could not open temp file"))
	}
	_, err = tempFile.WriteString(testModulesInfo)
	if err != nil {
		panic(errors.New("Could not write temp file"))
	}
	defer os.Remove(tempFile.Name())
	modulesInfo := LoadModulesInfo(tempFile.Name())
	if modulesInfo["module1"].RepoHash != "dbe955d1d83ea4ec969656d1e002e25ca1382fd8" {
		t.Error("Bad repo hash")
	}
	if modulesInfo["module1"].BuildOK != true {
		t.Error("Bad build status")
	}
	if modulesInfo["module2"].RepoHash != "c634c54781a89253167076ce102e588af8a60141" {
		t.Error("Bad repo hash")
	}
	if modulesInfo["modules2"].BuildOK != false {
		t.Error("Bad build status")
	}
}

const TestModulesInfoFile = "/tmp/test-repo-hash.yml"

func TestSaveModulesInfo(t *testing.T) {
	modulesInfo := ModulesInfo{
		"module1": ModuleInfo{
			RepoHash: "dbe955d1d83ea4ec969656d1e002e25ca1382fd8",
			BuildOK:  true,
		},
		"module2": ModuleInfo{
			RepoHash: "c634c54781a89253167076ce102e588af8a60141",
			BuildOK:  false,
		},
	}
	SaveModulesInfo(modulesInfo, TestModulesInfoFile)
	defer os.Remove(TestModulesInfoFile)
	actual, _ := ioutil.ReadFile(TestModulesInfoFile)
	if string(actual) != testModulesInfo {
		t.Error("Error writing repo file")
	}
}
