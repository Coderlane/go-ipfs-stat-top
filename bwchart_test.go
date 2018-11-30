package main

import (
	"context"
	"testing"

	tui "github.com/gizak/termui"
	ipfs "github.com/ipfs/go-ipfs-api"

	"github.com/stretchr/testify/assert"
)

// Test setting up and working with a history object
func TestBWHistory(t *testing.T) {
	bwh := newHistory(4)

	// Add a little data
	for i := 0; i < 2; i++ {
		bwh.Append(float64(i))
	}

	a := bwh.GetHistory()
	assert.Equal(t, []float64{0, 1}, a, "")

	// Add more data
	for i := 4; i < 8; i++ {
		bwh.Append(float64(i))
	}

	b := bwh.GetHistory()
	assert.Equal(t, []float64{4, 5, 6, 7}, b, "")

	// Shrink our history
	bwh.UpdateLen(2)

	c := bwh.GetHistory()
	assert.Equal(t, []float64{6, 7}, c, "")
}

// Test setting up a line chart
func TestBWLineChart(t *testing.T) {
	rs := tui.Resize{
		Width:  4,
		Height: 80,
	}

	bwc := NewBWLineChart(nil, rs)

	// Add data
	for i := 0; i < 16; i++ {
		bwc.refresh(float64(i), float64(i+1))
	}

	expected := map[string][]float64(map[string][]float64{
		"TX": []float64{8, 9, 10, 11, 12, 13, 14, 15},
		"RX": []float64{9, 10, 11, 12, 13, 14, 15, 16}})
	assert.Equal(t, bwc.LineChart.Data, expected, "")
}

// Test integration with IPFS
func TestBWLineChartIntegration(t *testing.T) {
	shell := ipfs.NewLocalShell()
	if shell == nil || !shell.IsUp() {
		t.Skipf("Could not connect to IPFS Daemon.")
		return
	}

	rs := tui.Resize{
		Width:  4,
		Height: 80,
	}

	bwc := NewBWLineChart(shell, rs)

	// See if we successfully refresh
	bwc.Refresh(context.Background())

	assert.Len(t, bwc.LineChart.Data["RX"], 1)
	assert.Len(t, bwc.LineChart.Data["TX"], 1)
}
