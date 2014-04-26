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

func TestMessage(t *testing.T) {
	expected := `From: from
To: to
Subject: subject

Build on 2001-02-03 12:34:

  foo: OK
  bar: ERROR
  baz: SKIPPED

Done in 1m2s
FAILURE

===================================
bar
-----------------------------------
Argh!
-----------------------------------

--
gontinuum`
	builds := []Build{
		Build{Module: ModuleConfig{Name: "foo"}, Success: true, Skipped: false},
		Build{Module: ModuleConfig{Name: "bar"}, Success: false, Skipped: false, Output: "Argh!"},
		Build{Module: ModuleConfig{Name: "baz"}, Success: true, Skipped: true},
	}
	start := time.Date(2001, 2, 3, 12, 34, 0, 0, time.Local)
	duration, _ := time.ParseDuration("1m2s")
	email := EmailConfig{Recipient: "to", Sender: "from"}
	actual := EmailMessage(builds, start, duration, email, "subject")
	if expected != actual {
		t.Error("Email mesage broken")
	}
}
