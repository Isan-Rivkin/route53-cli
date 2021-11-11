package v1

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type NodeInput struct {
	Txt        string
	Reference  interface{}
	Color      tcell.Color
	Selectable bool
	Expanded   bool
}

func NewNode(i *NodeInput) *tview.TreeNode {
	return tview.NewTreeNode(i.Txt).
		SetReference(i.Reference).
		SetSelectable(i.Selectable).
		SetColor(i.Color).
		SetExpanded(i.Expanded)
}

func NewRootNodeWithChildren(root *NodeInput, children []*NodeInput) *tview.TreeNode {
	rootNode := NewNode(root)

	for _, child := range children {
		n := NewNode(child)
		rootNode.AddChild(n)
	}

	return rootNode
}
