package main

import (
	"context"
	"fmt"

	"github.com/dustin/go-humanize"
	tui "github.com/gizak/termui"
	ipfs "github.com/ipfs/go-ipfs-api"
)

// RepoPar a paragraph describing the repository.
type RepoPar struct {
	shell *ipfs.Shell

	Par *tui.Par
}

// The stats from /api/v0/stats/repo
type repoStats struct {
	NumObjects uint64
	RepoSize   int64
	RepoPath   string
	Version    string
	StorageMax uint64
}

// NewRepoPar Create a new paragraph describing the repository's status.
func NewRepoPar(shell *ipfs.Shell) *RepoPar {
	rp := &RepoPar{
		shell: shell,
		Par:   tui.NewPar(""),
	}

	rp.Par.Border = true
	rp.Par.Height = 8

	return rp
}

// Resize NOOP, keep the paragraph at a static height.
func (rp *RepoPar) Resize(tui.Resize) {
	// We have a static height
}

// Refresh Collect the latest repository data
func (rp *RepoPar) Refresh(ctx context.Context) {
	var stats repoStats

	err := rp.shell.Request("stats/repo").Exec(ctx, &stats)

	if err != nil {
		rp.refresh(nil)
	} else {
		rp.refresh(&stats)
	}
}

func (rp *RepoPar) refresh(stats *repoStats) {

	if stats == nil {
		rp.Par.Text = "Error, unable to connect to IPFS Daemon"
	} else {
		rp.Par.Text = fmt.Sprintf(
			"Path: %v\nVersion: %v\nSize: %s\nMax Size: %s\nTotal Objects: %v\n",
			stats.RepoPath, stats.Version, humanize.Bytes(uint64(stats.RepoSize)),
			humanize.Bytes(stats.StorageMax), stats.NumObjects)
	}
}
