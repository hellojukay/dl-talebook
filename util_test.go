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
		newName := tosafeFileName(filename)
		if !IsValidFilename(newName) {
			t.Fail()
		}
	}

	filename := "80/20法则"
	newName := tosafeFileName(filename)
	if !IsValidFilename(newName) {
		t.Fail()
	}

}
