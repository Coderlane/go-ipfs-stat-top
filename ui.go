package main

import (
	tui "github.com/gizak/termui"
)

type UserInterface struct {
	Elements    []UserElement
	BWLineChart *BWLineChart
	RepoPar     *RepoPar
	BWTable     *BWTable
}

type UserElement interface {
	Resize(rs tui.Resize)
}

func NewUserInterface() *UserInterface {

	bwLineChart := NewBWLineChart()
	repoPar := NewRepoPar()
	bwTable := NewBWTable()

	ui := &UserInterface{
		BWLineChart: bwLineChart,
		RepoPar:     repoPar,
		BWTable:     bwTable,
		Elements:    []UserElement{bwLineChart, repoPar, bwTable},
	}
	ui.init()
	return ui
}

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

	// First render
	tui.Body.Align()
	tui.Render(tui.Body)
}

func (ui *UserInterface) loop() {
	tui.Render(tui.Body)
}

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

func (ui *UserInterface) quit(tui.Event) {
	tui.StopLoop()
}

func (ui *UserInterface) Run() {
	tui.Loop()
}
