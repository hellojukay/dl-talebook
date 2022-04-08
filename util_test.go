package main

import (
	"testing"
)

func TestRemoveChars(t *testing.T) {
	filename := `\Wha's Like Us? (Say It in Scots!).epub`
	if IsValidFilename(filename) {
		t.Fail()
	}
	newName := removeChars(filename, IllegalCharacters)
	if !IsValidFilename(newName) {
		t.Fail()
	}
}
