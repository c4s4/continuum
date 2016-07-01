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

const testRepoHashMap = `module1: dbe955d1d83ea4ec969656d1e002e25ca1382fd8
module2: c634c54781a89253167076ce102e588af8a60141
`

func TestLoadRepoHashMap(t *testing.T) {
	tempFile, err := ioutil.TempFile("/tmp", "go-test-")
	if err != nil {
		panic(errors.New("Could not open temp file"))
	}
	_, err = tempFile.WriteString(testRepoHashMap)
	if err != nil {
		panic(errors.New("Could not write temp file"))
	}
	defer os.Remove(tempFile.Name())
	repoHashMap := LoadRepoHashMap(tempFile.Name())
	if repoHashMap["module1"] != "dbe955d1d83ea4ec969656d1e002e25ca1382fd8" {
		t.Error("Bad repo hash")
	}
	if repoHashMap["module2"] != "c634c54781a89253167076ce102e588af8a60141" {
		t.Error("Bad repo hash")
	}
}

const TestRepoHashFile = "/tmp/test-repo-hash.yml"

func TestSaveRepoHash(t *testing.T) {
	repoHashMap := RepoHashMap{
		"module1": "dbe955d1d83ea4ec969656d1e002e25ca1382fd8",
		"module2": "c634c54781a89253167076ce102e588af8a60141",
	}
	SaveRepoHashMap(repoHashMap, TestRepoHashFile)
	defer os.Remove(TestRepoHashFile)
	actual, _ := ioutil.ReadFile(TestRepoHashFile)
	if string(actual) != testRepoHashMap {
		t.Error("Error writing repo file")
	}
}
