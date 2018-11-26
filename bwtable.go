package main

import (
	tui "github.com/gizak/termui"
	"math"
)

// BWTable A table to track bandwith usage by protocol
type BWTable struct {
	Table *tui.Table
}

var bwtHeader = []string{"Protocol", "Total In", "Total Out", "Rate In", "Rate Out"}

// NewBWTable Create a new table to track bandwith usage by protocol
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

// Resize Resize the table to match the new terminal size
func (bwt *BWTable) Resize(rs tui.Resize) {
	// Let the table take up 2/3 of the terminal
	bwt.Table.Height = int(math.Ceil((float64(rs.Height) / 3) * 2))
	if bwt.Table.Height < 8 {
		// At a minimum, 8 rows
		bwt.Table.Height = 8
	}
}
