package main

import (
	"testing"
)

const port = 7777

func TestIsAnotherInstanceRunning(t *testing.T) {
	expected := false
	actual := IsAnotherInstanceRunning(port)
	if expected != actual {
		t.Error("Should be no other instance running")
	}
	expected = true
	actual = IsAnotherInstanceRunning(port)
	if expected != actual {
		t.Error("Should be another instance running")
	}
}
