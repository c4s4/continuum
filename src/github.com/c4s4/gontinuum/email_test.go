package main

import (
	"testing"
	"time"
)

func TestSubjectSuccess(t *testing.T) {
	expected := "Build on 2001-02-03 12:34 was a SUCCESS"
	builds := []Build{
		Build{Success: true, Skipped: false},
		Build{Success: true, Skipped: true},
	}
	start := time.Date(2001, 2, 3, 12, 34, 0, 0, time.Local)
	actual := EmailSubject(builds, start)
	if expected != actual {
		t.Error("Email subject broken")
	}
}

func TestSubjectFailure(t *testing.T) {
	expected := "Build on 2001-02-03 12:34 was a FAILURE"
	builds := []Build{
		Build{Success: false, Skipped: false},
		Build{Success: true, Skipped: true},
	}
	start := time.Date(2001, 2, 3, 12, 34, 0, 0, time.Local)
	actual := EmailSubject(builds, start)
	if expected != actual {
		t.Error("Email subject broken")
	}
}
