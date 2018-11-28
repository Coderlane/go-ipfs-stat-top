package main

import (
	"fmt"
	"os"

	"github.com/gizak/termui"
	ipfs "github.com/ipfs/go-ipfs-api"
)

func main() {
	shell := ipfs.NewLocalShell()
	if shell == nil || !shell.IsUp() {
		fmt.Fprintf(os.Stderr, "Could not connect to IPFS Daemon.\n")
		return
	}

	if err := termui.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize termui.\nError: %v\n", err)
		return
	}
	defer termui.Close()

	ui := NewUserInterface(shell)
	ui.Run()
}
