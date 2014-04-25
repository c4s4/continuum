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
		Name:    "test",
		Url:     "ssh://casa@sweetohm.net/home/git/gontinuum.git",
		Command: "echo 'TEST'",
	}
	hash := GetRepoHash(moduleConfig)
	if match, _ := regexp.MatchString("^[0-9a-f]{40}$", hash); match != true {
		t.Errorf("GetRepoStatus() response '%s' doesn't look like a hash", hash)
	}
}

const testRepoHash = "module1: dbe955d1d83ea4ec969656d1e002e25ca1382fd8\n" +
	"module2: c634c54781a89253167076ce102e588af8a60141\n"

func TestLoadRepoHash(t *testing.T) {
	tempFile, err := ioutil.TempFile("/tmp", "go-test-")
	if err != nil {
		panic(errors.New("Could not open temp file"))
	}
	_, err = tempFile.WriteString(testRepoHash)
	if err != nil {
		panic(errors.New("Could not write temp file"))
	}
	defer os.Remove(tempFile.Name())
	repoHash := LoadRepoHash(tempFile.Name())
	if repoHash["module1"] != "dbe955d1d83ea4ec969656d1e002e25ca1382fd8" {
		t.Error("Bad repo hash")
	}
	if repoHash["module2"] != "c634c54781a89253167076ce102e588af8a60141" {
		t.Error("Bad repo hash")
	}
}

func TestSaveRepoHash(t *testing.T) {
	repoHash := RepoHash{
		"module1": "dbe955d1d83ea4ec969656d1e002e25ca1382fd8",
		"module2": "c634c54781a89253167076ce102e588af8a60141",
	}
	SaveRepoHash(repoHash, "/tmp/repo-hash.yml")
	defer os.Remove("/tmp/repo-hash.yml")
	actual, _ := ioutil.ReadFile("/tmp/repo-hash.yml")
	if string(actual) != testRepoHash {
		t.Error("Error writing repo file")
	}
}
