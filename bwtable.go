package main

import (
	tui "github.com/gizak/termui"
)

type BWTable struct {
	Table *tui.Table
}

var bwtHeader = []string{"Protocol", "Total In", "Total Out", "Rate In", "Rate Out"}

func NewBWTable() *BWTable {
	bwt := &BWTable{
		Table: tui.NewTable(),
	}

	bwt.Table.FgColor = tui.ColorWhite
	bwt.Table.BgColor = tui.ColorDefault

	// Add our default column headers
	bwt.Table.Rows = [][]string{
		bwtHeader,
	}

	return bwt
}

func (bwt *BWTable) Resize(rs tui.Resize) {
	// Let the table take up 2/3 of the terminal
	bwt.Table.Height = (rs.Height / 3) * 2
	if bwt.Table.Height < 8 {
		// At a minimum, 8 rows
		bwt.Table.Height = 8
	}
}
