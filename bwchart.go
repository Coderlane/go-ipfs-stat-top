package main

import (
	"context"
	"math"

	tui "github.com/gizak/termui"
	ipfs "github.com/ipfs/go-ipfs-api"
)

// Bandwith History, keep a history of floats
type bwHistory struct {
	history    []float64
	historyLen int
}

// Append a new float, popping the oldest if necessary
func (bwh *bwHistory) Append(a float64) {
	hlen := len(bwh.history)
	if hlen > bwh.historyLen {
		// Pop the last history entry
		bwh.history = bwh.history[2:]
	}

	bwh.history = append(bwh.history, a)
}

// newHistory create a new History object
func newHistory(mlen int) *bwHistory {
	return &bwHistory{
		history:    make([]float64, 0),
		historyLen: mlen,
	}
}

// UpdateLen update the maximum number of history entries
func (bwh *bwHistory) UpdateLen(mlen int) {
	if mlen == bwh.historyLen {
		// NOOP
		return
	}

	hlen := len(bwh.history)
	if mlen < hlen {
		// Shrink
		bwh.history = bwh.history[hlen-mlen:]
	}
	bwh.historyLen = mlen
}

// GetHistory get the history tracked by this object
func (bwh *bwHistory) GetHistory() []float64 {
	return bwh.history
}

// BWLineChart a LineChart for tracking total bandwith usage over time.
type BWLineChart struct {
	shell *ipfs.Shell

	LineChart *tui.LineChart

	txHistory *bwHistory
	rxHistory *bwHistory
}

var txLabel = "TX"
var rxLabel = "RX"

// NewBWLineChart Create a new Bandwith Line Chart to track total bandwith
// usage over time.
func NewBWLineChart(shell *ipfs.Shell, rs tui.Resize) *BWLineChart {
	bwc := &BWLineChart{
		shell:     shell,
		LineChart: tui.NewLineChart(),

		txHistory: newHistory(32),
		rxHistory: newHistory(32),
	}

	bwc.LineChart.BorderLabel = "Bandwith In/Out"
	bwc.LineChart.Mode = "dot"
	bwc.LineChart.YFloor = 0
	bwc.LineChart.Data[txLabel] = []float64{0}
	bwc.LineChart.Data[rxLabel] = []float64{0}
	bwc.LineChart.AxesColor = tui.ColorWhite
	bwc.LineChart.LineColor[txLabel] = tui.ColorGreen
	bwc.LineChart.LineColor[rxLabel] = tui.ColorBlue

	// Setup the initial size
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

	bwc.rxHistory.UpdateLen(rs.Width * 2)
	bwc.txHistory.UpdateLen(rs.Width * 2)

	// Figure out the maximum sizes
	bwc.LineChart.Data[txLabel] = bwc.txHistory.GetHistory()
	bwc.LineChart.Data[rxLabel] = bwc.rxHistory.GetHistory()
	bwc.LineChart.DataLabels = make([]string, len(bwc.LineChart.Data[txLabel]))
}

// Refresh Collect the latest series of data for the Line Chart
func (bwc *BWLineChart) Refresh(ctx context.Context) {
	// Get the latests stats, default to 0 on error
	tx := float64(0)
	rx := float64(0)

	stats, _ := bwc.shell.StatsBW(ctx)
	if stats != nil {
		tx = stats.RateOut / 1024
		rx = stats.RateIn / 1024
	}

	bwc.refresh(tx, rx)
}

func (bwc *BWLineChart) refresh(tx, rx float64) {
	bwc.txHistory.Append(tx)
	bwc.rxHistory.Append(rx)

	// Figure out the maximum sizes
	bwc.LineChart.Data[txLabel] = bwc.txHistory.GetHistory()
	bwc.LineChart.Data[rxLabel] = bwc.rxHistory.GetHistory()
	bwc.LineChart.DataLabels = make([]string, len(bwc.LineChart.Data[txLabel]))
}
