package main

import (
	tui "github.com/gizak/termui"
	"math"
)

// BWLineChart a LineChart for tracking total bandwith usage over time.
type BWLineChart struct {
	LineChart *tui.LineChart
}

var txLabel = "TX"
var rxLabel = "RX"

// NewBWLineChart Create a new Bandwith Line Chart to track total bandwith
// usage over time.
func NewBWLineChart() *BWLineChart {
	bwc := &BWLineChart{
		LineChart: tui.NewLineChart(),
	}

	bwc.LineChart.BorderLabel = "Bandwith In/Out"
	bwc.LineChart.Mode = "dot"
	bwc.LineChart.YFloor = 0
	bwc.LineChart.Data[txLabel] = []float64{3, 2, 1, 0}
	bwc.LineChart.Data[rxLabel] = []float64{0, 1, 2, 3}
	bwc.LineChart.AxesColor = tui.ColorWhite
	bwc.LineChart.LineColor[txLabel] = tui.ColorGreen
	bwc.LineChart.LineColor[rxLabel] = tui.ColorBlue

	// Setup the initial size
	rs := tui.Resize{
		Height: tui.TermHeight(),
		Width:  tui.TermWidth(),
	}
	bwc.Resize(rs)

	return bwc
}

// Resize Resize the line chart to match the new terminal size
func (bwc *BWLineChart) Resize(rs tui.Resize) {
	// Let the table take up 1/3 of the terminal
	bwc.LineChart.Height = int(math.Floor(float64(rs.Height) / 3))
	if bwc.LineChart.Height < 8 {
		// At a minimum, 8 rows
		bwc.LineChart.Height = 8
	}
}
