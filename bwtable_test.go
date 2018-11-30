package main

import (
	"testing"

	tui "github.com/gizak/termui"
	p2pmetrics "github.com/libp2p/go-libp2p-metrics"

	"github.com/stretchr/testify/assert"
)

// Test setting up a line chart
func TestBWTable(t *testing.T) {
	rs := tui.Resize{
		Width:  4,
		Height: 80,
	}

	bwt := NewBWTable(nil, rs)

	var stats protoStatSlice

	// Refresh with no data
	bwt.refresh(stats)

	expected := [][]string([][]string{[]string{
		"Protocol", "Total In", "Total Out", "Rate In", "Rate Out"}})
	assert.Equal(t, bwt.Table.Rows, expected, "")

	stats = append(stats, protoStat{
		Protocol: "Test",
		Stats: p2pmetrics.Stats{
			TotalIn:  1024,
			TotalOut: 204800,
		},
	})

	// Refresh with some data
	bwt.refresh(stats)

	expected = append(expected,
		[]string{"Test", "  1.0 kB", "  205 kB", "     0 B", "     0 B"})
	assert.Equal(t, bwt.Table.Rows, expected, "")
}
