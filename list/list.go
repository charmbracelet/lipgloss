package list

import "github.com/charmbracelet/lipgloss/tree"

// New returns a new list.
func New(items ...any) *tree.TreeNode {
	t := tree.New("").
		ItemAddFunc(listItemAddFn).
		Enumerator(Bullet)
	for _, item := range items {
		t = t.Item(item)
	}
	return t
}

// the main difference from this and the tree's default itemAddFunc, is that
// this one has some syntactic sugar to make declaring sublists feel easier.
// e.g.: New("foo", "bar", New("zaz")) would eval to:
// - foo
// - bar
//   - zaz
func listItemAddFn(nodes []tree.Node, item any) []tree.Node {
	switch item := item.(type) {
	case *tree.TreeNode:
		newItem, rm := parentize(nodes, item)
		if rm >= 0 {
			nodes = remove(nodes, rm)
		}
		nodes = append(nodes, newItem)
		return nodes
	case tree.Node:
		nodes = append(nodes, item)
	case string:
		s := tree.StringNode(item)
		nodes = append(nodes, &s)
	}
	return nodes
}

// walks backwards in the existing nodes until it finds a string node, then
// remove it from the list and set it as the parent of the current node.
func parentize(nodes []tree.Node, item tree.Node) (tree.Node, int) {
	for j := len(nodes) - 1; j >= 0; j-- {
		parent := nodes[j]
		switch parent := parent.(type) {
		case tree.StringNode:
			return treenize(parent, item.Children()), j
		case *tree.StringNode:
			return treenize(parent, item.Children()), j
		}
	}
	return item, -1
}

// creates a new TreeNode with the given name and children.
func treenize(parent tree.Node, children []tree.Node) *tree.TreeNode {
	data := make([]any, 0, len(children))
	for _, d := range children {
		data = append(data, d)
	}
	return tree.New(parent.Name(), data...)
}

func remove(data []tree.Node, i int) []tree.Node {
	return append(data[:i], data[i+1:]...)
}
