package main

import (
	"context"
	"testing"

	tui "github.com/gizak/termui"
	ipfs "github.com/ipfs/go-ipfs-api"
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

// Test integration with IPFS
func TestBWTableIntegration(t *testing.T) {
	shell := ipfs.NewLocalShell()
	if shell == nil || !shell.IsUp() {
		t.Skipf("Could not connect to IPFS Daemon.")
		return
	}

	rs := tui.Resize{
		Width:  4,
		Height: 80,
	}

	bwt := NewBWTable(shell, rs)

	// See if we successfully refresh
	bwt.Refresh(context.Background())

	if len(bwt.Table.Rows) < 1 {
		t.Errorf("Too few rows: %d", len(bwt.Table.Rows))
	}
}
