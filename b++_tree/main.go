package main

import (
	"fmt"
)

const max_keys_size = 4

type BPlusTreeNode struct {
	child_bplus_tree_nodes [max_keys_size]*BPlusTreeNode
	child_bplus_tree_keys  [max_keys_size]int

	is_leaf_node    bool
	next_bplus_node *BPlusTreeNode

	total_keys int

	parent *BPlusTreeNode
}

type BPlusTree struct {
	root *BPlusTreeNode
}

func NewBPlusTree() *BPlusTree {
	return &BPlusTree{root: &BPlusTreeNode{is_leaf_node: true}}
}

func (bplus_tree *BPlusTree) insert(key int) {
	if bplus_tree.root == nil {
		bplus_tree.root = &BPlusTreeNode{is_leaf_node: true}
	}

	if bplus_tree.root.total_keys == max_keys_size {
		old_root := bplus_tree.root
		bplus_tree.root = &BPlusTreeNode{}
		bplus_tree.root.child_bplus_tree_nodes[0] = old_root

		bplus_tree.splitChild(bplus_tree.root, 0)
	}

	bplus_tree.insertNotNull(bplus_tree.root, key)
}

func (bplus_tree *BPlusTree) splitChild(parent *BPlusTreeNode, index int) {
	new_child := &BPlusTreeNode{}
	old_child := parent.child_bplus_tree_nodes[index]

	new_child.is_leaf_node = old_child.is_leaf_node

	new_child.total_keys = max_keys_size / 2

	for ii := 0; ii < (max_keys_size / 2); ii++ {
		new_child.child_bplus_tree_keys[ii] =
			old_child.child_bplus_tree_keys[ii+(max_keys_size/2)]
	}

	if !new_child.is_leaf_node {
		for ii := 0; ii < (max_keys_size / 2); ii++ {
			new_child.child_bplus_tree_nodes[ii] =
				old_child.child_bplus_tree_nodes[ii+(max_keys_size/2)]
		}
	}

	old_child.total_keys = (max_keys_size) / 2

	for ii := parent.total_keys; ii > index; ii-- {
		parent.child_bplus_tree_nodes[ii+1] = parent.child_bplus_tree_nodes[ii]
	}

	parent.child_bplus_tree_nodes[index+1] = new_child

	for ii := parent.total_keys - 1; ii >= index; ii-- {
		parent.child_bplus_tree_keys[ii+1] = parent.child_bplus_tree_keys[ii]
	}
	parent.child_bplus_tree_keys[index] = old_child.child_bplus_tree_keys[(max_keys_size-1)/2]

	parent.total_keys++
}

func (bplus_tree *BPlusTree) insertNotNull(node *BPlusTreeNode, key int) {
	index := node.total_keys - 1

	if node.is_leaf_node {
		for index >= 0 && key < node.child_bplus_tree_keys[index] {
			node.child_bplus_tree_keys[index+1] = node.child_bplus_tree_keys[index]
			index--
		}

		node.child_bplus_tree_keys[index+1] = key
		node.total_keys++
	} else {
		for index >= 0 && key < node.child_bplus_tree_keys[index] {
			index--
		}
		index++

		if node.child_bplus_tree_nodes[index].total_keys == max_keys_size {
			bplus_tree.splitChild(node, index)
			if key > node.child_bplus_tree_keys[index] {
				index++
			}
		}

		bplus_tree.insertNotNull(node.child_bplus_tree_nodes[index], key)
	}
}

func (bplus_tree *BPlusTree) search(key int) bool {
	return bplus_tree.root.search(key)
}

func (bplus_treenode *BPlusTreeNode) search(key int) bool {
	ii := 0

	for ii < bplus_treenode.total_keys && key > bplus_treenode.child_bplus_tree_keys[ii] {
		ii++
	}

	if ii < bplus_treenode.total_keys && key == bplus_treenode.child_bplus_tree_keys[ii] {
		return true
	}

	if bplus_treenode.is_leaf_node {
		return false
	}

	return bplus_treenode.child_bplus_tree_nodes[ii].search(ii)
}

func (bplus_tree *BPlusTree) Print() {
	bplus_tree.root.print(0)
}

func (bplus_treenode *BPlusTreeNode) print(level int) {
	if bplus_treenode == nil {
		return
	}

	fmt.Printf("Level %d: ", level)
	for i := 0; i < bplus_treenode.total_keys; i++ {
		fmt.Printf("%d ", bplus_treenode.child_bplus_tree_keys[i])
	}
	fmt.Printf("\n")

	if !bplus_treenode.is_leaf_node {
		for i := 0; i <= bplus_treenode.total_keys; i++ {
			bplus_treenode.child_bplus_tree_nodes[i].print(level + 1)
		}
	}
}

func main() {
	tree := NewBPlusTree()

	keys := []int{3, 7, 2, 8, 5, 1, 4, 9, 6}
	for _, key := range keys {
		tree.insert(key)
	}
	tree.Print()

}
