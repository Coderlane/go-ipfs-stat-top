package main

import (
	"context"

	tui "github.com/gizak/termui"
)

// RepoPar a paragraph describing the repository.
type RepoPar struct {
	Par *tui.Par
}

// NewRepoPar Create a new paragraph describing the repository's status.
func NewRepoPar() *RepoPar {
	rp := &RepoPar{
		Par: tui.NewPar("TODO"),
	}

	rp.Par.Border = true

	return rp
}

// Resize NOOP, keep the paragraph at a static height.
func (rp *RepoPar) Resize(tui.Resize) {
	// We have a static height
	rp.Par.Height = 8
}

// Refresh Collect the latest repository data
func (rp *RepoPar) Refresh(ctx context.Context) {
	// TODO
}
