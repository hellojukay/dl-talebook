package main

import (
	"runtime/debug"
)

func PrintVersion() {
	if info, ok := debug.ReadBuildInfo(); ok {
		println(info.Main.Version)
	}
}
