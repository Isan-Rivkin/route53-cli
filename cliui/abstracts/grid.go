package abstracts

import (
	"github.com/rivo/tview"
)

type BaseGrid interface {
	RawGrid() *tview.Grid
}

type RootGrid struct {
	Grid *tview.Grid
}

func NewRootGrid() BaseGrid {
	return &RootGrid{
		Grid: tview.NewGrid(),
	}
}

func (g *RootGrid) RawGrid() *tview.Grid {
	return g.Grid
}

func (g *RootGrid) Render() *tview.Grid {

	return g.Grid
}
