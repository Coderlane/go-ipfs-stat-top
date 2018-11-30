package main

import (
	"context"
	"fmt"
	"math"
	"sort"

	"github.com/dustin/go-humanize"
	tui "github.com/gizak/termui"
	ipfs "github.com/ipfs/go-ipfs-api"
	p2pmetrics "github.com/libp2p/go-libp2p-metrics"
)

// BWTable A table to track bandwith usage by protocol
type BWTable struct {
	shell *ipfs.Shell

	Table *tui.Table

	protocols map[string]interface{}
}

// A protocol name combined with its statistics
type protoStat struct {
	Protocol string
	Stats    p2pmetrics.Stats
}

// A list of protocols and their statistics
type protoStatSlice []protoStat

func (p protoStatSlice) Len() int { return len(p) }

// Order based on total input and output
func (p protoStatSlice) Less(i, j int) bool {
	istat := p[i].Stats
	jstat := p[j].Stats

	cmp := (istat.TotalIn + istat.TotalOut) - (jstat.TotalIn + jstat.TotalOut)
	if cmp == 0 {
		return p[i].Protocol < p[j].Protocol
	}

	return cmp > 0
}

// Implement the sortable interface
func (p protoStatSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// A single stream in the lists of streams
type streamData struct {
	Protocol string
}

// A single peer in the list of peers
type peerData struct {
	Addr    string
	Peer    string
	Latency string
	Streams []streamData
}

// Data returned from "swarm/peers"
type peersData struct {
	Peers []peerData
}

var bwtHeader = []string{"Protocol",
	"Total In", "Total Out", "Rate In", "Rate Out"}

// NewBWTable Create a new table to track bandwith usage by protocol
func NewBWTable(shell *ipfs.Shell, rs tui.Resize) *BWTable {
	bwt := &BWTable{
		shell:     shell,
		Table:     tui.NewTable(),
		protocols: make(map[string]interface{}),
	}

	bwt.Table.FgColor = tui.ColorWhite
	bwt.Table.BgColor = tui.ColorDefault

	// Add our default column headers
	bwt.Table.Rows = [][]string{
		bwtHeader,
	}

	bwt.Resize(rs)

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

// Refresh Collect the latest per-protocol bandwith data
func (bwt *BWTable) Refresh(ctx context.Context) {
	var (
		peers peersData
		stats protoStatSlice
	)

	req := bwt.shell.Request("swarm/peers")
	req = req.Option("streams", true)

	err := req.Exec(ctx, &peers)
	if err != nil {
		bwt.refresh(stats)
		return
	}

	for _, peer := range peers.Peers {
		for _, stream := range peer.Streams {
			_, ok := bwt.protocols[stream.Protocol]
			if !ok {
				bwt.protocols[stream.Protocol] = struct{}{}
			}
		}
	}

	for proto := range bwt.protocols {
		var stat protoStat

		stat.Protocol = proto

		req := bwt.shell.Request("stats/bw")
		req.Option("proto", proto)
		_ = req.Exec(ctx, &stat.Stats)

		stats = append(stats, stat)
	}

	bwt.refresh(stats)
}

// Print a padded representation of the bytes
func paddedHumanBytes(bytes uint64) string {
	return fmt.Sprintf("%8s", humanize.Bytes(bytes))
}

// Update the table with the latest statistics
func (bwt *BWTable) refresh(stats protoStatSlice) {
	rows := [][]string{bwtHeader}

	sort.Sort(stats)
	for _, stat := range stats {
		row := []string{stat.Protocol,
			paddedHumanBytes(uint64(stat.Stats.TotalIn)),
			paddedHumanBytes(uint64(stat.Stats.TotalOut)),
			paddedHumanBytes(uint64(stat.Stats.RateIn)),
			paddedHumanBytes(uint64(stat.Stats.RateOut))}
		rows = append(rows, row)
	}

	bwt.Table.Rows = rows
}
