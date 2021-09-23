package abstracts

import (
	"io/ioutil"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type InteractiveTree struct {
	RawTree *tview.TreeView
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
func NewInteractiveTree() *InteractiveTree {
	// here a create dummy root tree for elb
	// need to rethink now how this gonna work
	// go run  main.go into -R -d 3 -r 'pro.similarweb.com'
	// will display dummy tree, before implementation i need to think more about the tree and resource connectivity,
	// also how do i represent all?
	//tree := dummyFSTree()
	tree := dummyElbRootTree()
	return &InteractiveTree{
		RawTree: tree,
	}
}
