package main

import (
	"testing"
	"time"
)

func TestSubjectSuccess(t *testing.T) {
	expected := "Foo build on 2001-02-03 12:34 was a SUCCESS"
	build := Build{Module: ModuleConfig{Name: "Foo"}, Success: true, Skipped: false, Output: "Argh!"}
	start := time.Date(2001, 2, 3, 12, 34, 0, 0, time.Local)
	actual := EmailSubject(build, start)
	if expected != actual {
		t.Error("Email subject broken")
	}
}

func TestSubjectFailure(t *testing.T) {
	expected := "Foo build on 2001-02-03 12:34 was a FAILURE"
	build := Build{Module: ModuleConfig{Name: "Foo"}, Success: false, Skipped: false, Output: "Argh!"}
	start := time.Date(2001, 2, 3, 12, 34, 0, 0, time.Local)
	actual := EmailSubject(build, start)
	if expected != actual {
		t.Error("Email subject broken")
	}
}

func TestMessage(t *testing.T) {
	expected := `From: from
To: to
Subject: subject

Foo build on 2001-02-03 12:34 was a FAILURE.
Done in 1m2s

===================================
ERROR:
-----------------------------------
Argh!
-----------------------------------

--
continuum`
	build := Build{Module: ModuleConfig{Name: "Foo"}, Success: false, Skipped: false, Output: "Argh!"}
	start := time.Date(2001, 2, 3, 12, 34, 0, 0, time.Local)
	duration, _ := time.ParseDuration("1m2s")
	email := EmailConfig{Recipient: "to", Sender: "from"}
	actual := EmailMessage(build, start, duration, email, "subject")
	if expected != actual {
		t.Error("Email message broken")
	}
}
