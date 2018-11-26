package main

import (
	"fmt"
	"github.com/gizak/termui"
	"os"
)

func main() {
	if err := termui.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize termui.\nError: %v\n", err)
		return
	}

	ui := NewUserInterface()
	ui.Run()

	defer termui.Close()
}
