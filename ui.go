package main

import (
	"context"
	"time"

	tui "github.com/gizak/termui"
	ipfs "github.com/ipfs/go-ipfs-api"
)

// UserInterface The full user-facing interface.
type UserInterface struct {
	Context context.Context
	Shell   *ipfs.Shell
	schan   chan bool

	Elements    []UserElement
	BWLineChart *BWLineChart
	RepoPar     *RepoPar
	BWTable     *BWTable
}

// UserElement a single element in the UserInterface.
type UserElement interface {
	Resize(rs tui.Resize)
	Refresh(ctx context.Context)
}

// NewUserInterface Create a new user interface
func NewUserInterface(shell *ipfs.Shell) *UserInterface {
	// Get the initial size
	rs := tui.Resize{
		Height: tui.TermHeight(),
		Width:  tui.TermWidth(),
	}

	// Setup the widgets
	bwLineChart := NewBWLineChart(shell, rs)
	repoPar := NewRepoPar(shell)
	bwTable := NewBWTable(shell, rs)

	ui := &UserInterface{
		Context: context.Background(),
		Shell:   shell,
		schan:   make(chan bool, 1),

		Elements:    []UserElement{bwLineChart, repoPar, bwTable},
		BWLineChart: bwLineChart,
		RepoPar:     repoPar,
		BWTable:     bwTable,
	}
	ui.init()
	return ui
}

// init Initialize the UI for the first time
func (ui *UserInterface) init() {
	// Setup Keyboard Handlers
	tui.Handle("q", ui.quit)
	tui.Handle("C-q", ui.quit)
	tui.Handle("C-c", ui.quit)

	// Setup Event Handlers
	tui.Handle("<Resize>", ui.resize)

	// Layout the widgets, kinda like this:
	// +------------------+ +----------+
	// |     BW Chart     | | Repo Par |
	// |                  | +----------+
	// |                  |
	// +------------------+
	// +-------------------------------+
	// |            BW Table           |
	// +-------------------------------+
	summary := tui.NewRow(
		tui.NewCol(8, 0, ui.BWLineChart.LineChart),
		tui.NewCol(4, 0, ui.RepoPar.Par))
	table := tui.NewRow(
		tui.NewCol(12, 0, ui.BWTable.Table))
	tui.Body.AddRows(summary, table)

	// Fill out the intial sizes
	rs := tui.Resize{
		Height: tui.TermHeight(),
		Width:  tui.TermWidth(),
	}
	for _, elem := range ui.Elements {
		elem.Resize(rs)
	}
	tui.Body.Align()

	// First render
	ui.refresh()
}

func (ui *UserInterface) refresh() {
	for _, elem := range ui.Elements {
		elem.Refresh(ui.Context)
	}

	tui.Clear()
	tui.Render(tui.Body)
}

// loop Call in a loop to re-render
func (ui *UserInterface) loop() {

	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ticker.C:
			ui.refresh()
		case <-ui.schan:
			return
		}
	}
}

// resize Resize the UI to match the event
func (ui *UserInterface) resize(e tui.Event) {
	rs := e.Payload.(tui.Resize)

	// Resize all of the user elments
	for _, elem := range ui.Elements {
		elem.Resize(rs)
	}

	// Resize the terminal body
	tui.Body.Width = rs.Width
	tui.Body.Align()

	// Re-render
	tui.Clear()
	tui.Render(tui.Body)
}

// quit Quit out of the userinterface
func (ui *UserInterface) quit(tui.Event) {
	ui.schan <- true
	tui.StopLoop()
}

// Run Run the main loop
func (ui *UserInterface) Run() {

	go ui.loop()
	tui.Loop()
}
