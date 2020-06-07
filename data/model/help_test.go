package model

import "testing"

func TestHelpModel(t *testing.T) {
	got := Help{
		"help",
		"displays all available commands",
		[]string{
			"help",
		},
	}
	want := "help"

	if got.Command != want {
		t.Errorf("got %s want %s", got, want)
	}
}
