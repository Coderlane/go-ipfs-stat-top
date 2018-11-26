package main

import (
	tui "github.com/gizak/termui"
)

type RepoPar struct {
	Par *tui.Par
}

func NewRepoPar() *RepoPar {
	rp := &RepoPar{
		Par: tui.NewPar("TODO"),
	}

	rp.Par.Height = 8
	rp.Par.Border = true

	return rp
}

func (rp *RepoPar) Resize(tui.Resize) {
	// We have a static height
	rp.Par.Height = 8
}
