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

// gridItem represents one primitive and its possible position on a grid.
type GridItem struct {
	Item     tview.Primitive // The item to be positioned. May be nil for an empty item.
	Metadata GridItemMetadata
}

func NewDefaultGridItem(p tview.Primitive) *GridItem {
	m := GridItemMetadata{
		RowSpan: 1,
		ColSpan: 3,
	}

	return &GridItem{
		Item:     p,
		Metadata: m,
	}
}

// p Primitive, row, column, rowSpan, colSpan, minGridHeight, minGridWidth int, focus bool
func (gi *GridItem) ConvertToAddItemInput() (tview.Primitive, int, int, int, int, int, int, bool) {
	return gi.Item, gi.Metadata.Row, gi.Metadata.Column, gi.Metadata.RowSpan, gi.Metadata.ColSpan, gi.Metadata.MinGridHeight, gi.Metadata.MinGridWidth, gi.Metadata.Focus
}

type GridItemMetadata struct {
	Row, Column                 int  // The top-left grid cell where the item is placed.
	ColSpan, RowSpan            int  // The number of rows and columns the item occupies.
	MinGridWidth, MinGridHeight int  // The minimum grid width/height for which this item is visible.
	Focus                       bool // Whether or not this item attracts the layout's focus.
}
