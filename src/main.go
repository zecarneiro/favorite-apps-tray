package main

import (
	"main/src/lib"

	"github.com/energye/systray"
)

func main() {
	systray.Run(onReady, onQuit)
}

func onReady() {
	lib.Start()
}

func onQuit() {
	// TODO: Not implemented yet
}
