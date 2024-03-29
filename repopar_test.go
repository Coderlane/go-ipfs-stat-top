package main

import (
	"context"
	"testing"

	ipfs "github.com/ipfs/go-ipfs-api"

	"github.com/stretchr/testify/assert"
)

// Test setting up a repository paragraph
func TestRepoPar(t *testing.T) {

	rp := NewRepoPar(nil)

	// Try failing to refresh
	rp.refresh(nil)

	expected := "Error, unable to connect to IPFS Daemon"
	assert.Equal(t, rp.Par.Text, expected, "")

	// Refresh with some test data
	stats := &repoStats{
		NumObjects: 16,
		RepoSize:   102400,
		StorageMax: 10000000000,
		Version:    "test",
		RepoPath:   "test",
	}
	rp.refresh(stats)

	expected = "Path: test\nVersion: test\nSize: 102 kB\n" +
		"Max Size: 10 GB\nTotal Objects: 16\n"
	assert.Equal(t, rp.Par.Text, expected, "")
}

// Test integration with IPFS
func TestRepoParIntegration(t *testing.T) {
	shell := ipfs.NewLocalShell()
	if shell == nil || !shell.IsUp() {
		t.Skipf("Could not connect to IPFS Daemon.")
		return
	}

	rp := NewRepoPar(shell)

	// See if we successfully refresh
	rp.Refresh(context.Background())

	assert.Regexp(t, ".*Path:.*", rp.Par.Text)
}
