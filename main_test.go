package main

import (
	"testing"
)

func TestGetCommands(t *testing.T) {
	commands := getCommands()
	if commands.Commands == nil {
		t.Fail()
	}
}

func BenchmarkGetCommands(b *testing.B) {
	commands := getCommands()
	if commands.Commands == nil {
		b.Error("Unexpected result")
	}
}
