package main

import (
	"runtime"
	"testing"
)

func TestRemoveChars(t *testing.T) {
	if runtime.GOOS == "windows" {
		filename := `\Wha's Like Us? (Say It in Scots!).epub`
		if IsValidFilename(filename) {
			t.Fail()
		}
		newName := removeChars(filename, IllegalCharacters)
		if !IsValidFilename(newName) {
			t.Fail()
		}
	}

}
