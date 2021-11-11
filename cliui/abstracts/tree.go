package abstracts

import (
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func changeAsync(n *tview.TreeNode) {
	go func() {
		time.Sleep(time.Duration(3) * time.Second)
		a := tview.NewTreeNode("A").
			SetReference("A").
			SetSelectable(true).
			SetColor(tcell.ColorYellow).
			SetExpanded(true)
		b := tview.NewTreeNode("B").
			SetReference("B").
			SetSelectable(true).
			SetColor(tcell.ColorYellow).
			SetExpanded(true)

		n.AddChild(a)
		n.AddChild(b)
	}()
}

func dummyElbRootTree() *tview.TreeView {
	rootDir := "dualstack.some-pp-1011.us-east-3.elb.amazonaws.com."
	root := tview.NewTreeNode(rootDir).
		SetColor(tcell.ColorRed).
		SetExpanded(true)

	listeners := tview.NewTreeNode("Listeners").
		SetReference("listenersRef").
		SetSelectable(true).
		SetColor(tcell.ColorGreen).
		SetExpanded(true)

	certs := tview.NewTreeNode("Certificates").
		SetReference("certsRef").
		SetSelectable(true).
		SetExpanded(true)

	targetGroup := tview.NewTreeNode("Target Groups").
		SetReference("tgRef").
		SetSelectable(true).
		SetColor(tcell.ColorGreen).
		SetExpanded(true)

	ec2Instances := tview.NewTreeNode("EC2 Instances").
		SetReference("ec2Ref").
		SetSelectable(true).
		SetColor(tcell.ColorGreen).
		SetExpanded(true)

	targetGroup.AddChild(ec2Instances)

	listeners.AddChild(certs)
	listeners.AddChild(targetGroup)

	root.AddChild(listeners)

	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return // Selecting the root node does nothing.
		}
		children := node.GetChildren()
		if len(children) == 0 {
			changeAsync(node)
			// Load and show files in this directory.
			//path := reference.(string)

		} else {
			// Collapse if visible, expand if collapsed.
			node.SetExpanded(!node.IsExpanded())
		}
	})
	// A helper function which adds the files and directories of the given path
	// to the given target node.
	// r53Record:
	// -> click (elb)
	// -> E: elb
	// --> E: network info...
	// --> E: certs
	// --> E: TG
	// --> E: Listeners
	// ---> E: Ec2 instances
	// ----> E: network info...

	// add := func(target *tview.TreeNode, path string) {
	// 	files, err := ioutil.ReadDir(path)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	for _, file := range files {
	// 		node := tview.NewTreeNode(file.Name()).
	// 			SetReference(filepath.Join(path, file.Name())).
	// 			SetSelectable(file.IsDir())
	// 		if file.IsDir() {
	// 			node.SetColor(tcell.ColorGreen)
	// 		}
	// 		target.AddChild(node)
	// 	}
	// }

	// tree.SetChangedFunc(func(node *tview.TreeNode) {
	// 	fmt.Println("hovered a node!!!! ")
	// 	reference := node.GetReference()
	// 	if reference == nil {
	// 		fmt.Println("Hover: root selected @@")
	// 		return
	// 	}
	// 	val := reference.(string)
	// 	fmt.Println("Hover: ", val)
	// })
	return tree
}

func dummyFSTree() *tview.TreeView {
	rootDir := "."
	root := tview.NewTreeNode(rootDir).
		SetColor(tcell.ColorRed)
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	// A helper function which adds the files and directories of the given path
	// to the given target node.
	add := func(target *tview.TreeNode, path string) {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			node := tview.NewTreeNode(file.Name()).
				SetReference(filepath.Join(path, file.Name())).
				SetSelectable(file.IsDir())
			if file.IsDir() {
				node.SetColor(tcell.ColorGreen)
			}
			target.AddChild(node)
		}
	}

	// Add the current directory to the root node.
	add(root, rootDir)

	// If a directory was selected, open it.
	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return // Selecting the root node does nothing.
		}
		children := node.GetChildren()
		if len(children) == 0 {
			// Load and show files in this directory.
			path := reference.(string)
			add(node, path)
		} else {
			// Collapse if visible, expand if collapsed.
			node.SetExpanded(!node.IsExpanded())
		}
	})
	return tree
}

type TreeHandler func(node *tview.TreeNode)

type InteractiveTree struct {
	RawTree *tview.TreeView
}

func NewInteractiveTree(root *tview.TreeNode, onSelect TreeHandler) *InteractiveTree {

	//tree := dummyFSTree()
	//tree := dummyElbRootTree()
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	// when a user clicks on a node
	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		onSelect(node)
	})

	return &InteractiveTree{
		RawTree: tree,
	}
}

// Initial:
// - tree: root node
// - tree: display reachable types
// - tree: display expanded names
// - info: display root info

// OnSelectedNode:
// - tree: get expandable resources
