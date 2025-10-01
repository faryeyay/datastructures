// Package bplustree
package bplustree

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
)

// BPlusTree - Represents a B+ Tree. We'll want to track the root node of the B+ Tree.
type BPlusTree[K cmp.Ordered, V any] struct {
	root    *Node[K, V]
	order   int // The order of the tree
	maxKeys int // The maximum number of keys in a node (2*order - 1)
	minKeys int // The minimum number of keys in a node (order - 1)
}

// CreateTree - creates an empty B+ Tree
func CreateTree[K cmp.Ordered, V any](order int) *BPlusTree[K, V] {
	maxKeys := calculateMaxKeys(order)
	minKeys := order - 1

	// create the B+ tree
	return &BPlusTree[K, V]{
		order:   order,
		maxKeys: maxKeys,
		minKeys: minKeys,
	}
}

// CreateTreeWithRootKey - create a new B+ Tree with a root key
func CreateTreeWithRootKey[K cmp.Ordered, V any](order int, rootKey K) *BPlusTree[K, V] {
	// create the root node
	rootNode := NewNode[K, V](rootKey, *new(V), false)
	maxKeys := calculateMaxKeys(order)
	minKeys := order - 1

	// create the B+ Tree
	return &BPlusTree[K, V]{
		root:    rootNode,
		order:   order,
		maxKeys: maxKeys,
		minKeys: minKeys,
	}
}

// calculateMaxKeys - a func that returns an integer for the maximum keys of a tree based on the
// order of the tree
func calculateMaxKeys(order int) int {
	return 2*order - 1
}

// GetRootNode - Returns the root node of the B+ Tree
// *Node[K, V] is the node that we are returning
func (b *BPlusTree[K, V]) GetRootNode() (*Node[K, V], error) {
	if b.root != nil {
		return b.root, nil
	}

	return nil, fmt.Errorf("the b+ tree doesn't have a root")
}

// Insert - inserts a node into the B+ Tree
func (b *BPlusTree[K, V]) Insert(key K, value V) error {
	// If the tree is empty, create a root that is a leaf.
	if b.root == nil {
		b.root = NewNode[K, V](key, value, true)
		return nil
	}

	// Traverse the tree to find the correct leaf node.
	leaf, parentNode, err := b.leafSearch(key, b.root, nil)
	if err != nil {
		return err
	}

	// Add the key and value to the leaf.
	leaf.AddKeyValue(key, value)

	// If the leaf is not full, we are done.
	if leaf.GetNumberOfKeys() <= b.maxKeys {
		return nil
	}

	// --- The leaf node is full, so we need to split it. ---

	// Create a new leaf node.
	newLeaf := NewEmptyNode[K, V]()
	newLeaf.leaf = true

	// Link the leaves.
	newLeaf.next = leaf.next
	leaf.next = newLeaf

	// Get all keys and find the split point.
	allKeys := leaf.GetIterableKeys()
	allValues := leaf.GetValues()
	splitPoint := (len(allKeys) + 1) / 2

	// The key to be copied up to the parent is the first key of the new leaf.
	promotedKey := allKeys[splitPoint]

	// Move the upper half of the keys from the original leaf to the new leaf.
	// Also move the corresponding values.
	for i := splitPoint; i < len(allKeys); i++ {
		newLeaf.AddKeyValue(allKeys[i], allValues[i])
	}

	// Truncate the keys in the original leaf.
	leaf.keys.Data = allKeys[:splitPoint]
	// Truncate the values in the original leaf.
	leaf.values = allValues[:splitPoint]

	// Insert the promoted key into the parent. This may cause further splits.
	return b.insertInParent(parentNode, promotedKey, newLeaf)
}

// insertInParent handles inserting a key and a new child into a parent node,
// and propagates splits upwards if necessary.
func (b *BPlusTree[K, V]) insertInParent(parent *Node[K, V], keyToInsert K, newChild *Node[K, V]) error {
	// If the parent is nil, it means we split the root. We need to create a new root.
	if parent == nil {
		newRoot := NewEmptyNode[K, V]()
		newRoot.leaf = false
		newRoot.AddKey(keyToInsert)
		// The old root and the new child become children of the new root.
		newRoot.AddChild(b.root)
		newRoot.AddChild(newChild)
		b.root = newRoot
		return nil
	}

	// Insert the key and child pointer into the parent.
	parent.AddKey(keyToInsert)
	parent.AddChild(newChild)

	// If the parent is not full, we are done.
	if parent.GetNumberOfKeys() <= b.maxKeys {
		return nil
	}

	// --- The parent node is full, so we need to split it. ---

	// Find the parent of the current parent node to prepare for the next promotion.
	// This is inefficient but necessary without parent pointers.
	grandParent, _ := b.findParent(b.root, parent)

	// Create a new internal node.
	newInternal := NewEmptyNode[K, V]()
	newInternal.leaf = false

	// Get all keys and children from the full parent.
	allKeys := parent.GetIterableKeys()
	allChildren := parent.GetIterableChildren()

	// Find the median key to promote. For internal nodes, this key is moved up, not copied.
	keySplitPoint := len(allKeys) / 2
	keyToPromote := allKeys[keySplitPoint]

	// Move keys to the right of the median to the new internal node.
	for i := keySplitPoint + 1; i < len(allKeys); i++ {
		newInternal.AddKey(allKeys[i])
	}

	// Truncate keys in the original parent node.
	parent.keys.Data = allKeys[:keySplitPoint]

	// The children must also be split.
	childSplitPoint := keySplitPoint + 1
	for i := childSplitPoint; i < len(allChildren); i++ {
		newInternal.AddChild(allChildren[i])
	}
	parent.children = allChildren[:childSplitPoint]

	// Recursively call insertInParent for the next level up.
	return b.insertInParent(grandParent, keyToPromote, newInternal)
}

// Delete - deletes a key from the B+ Tree, rebalancing if necessary.
func (b *BPlusTree[K, V]) Delete(key K) error {
	// Find the leaf node where the key should exist.
	leaf, parent, err := b.leafSearch(key, b.root, nil)
	if err != nil {
		return err // Should only happen on an empty tree.
	}

	// Remove the key from the leaf. If it doesn't exist, it's a no-op.
	if err := leaf.RemoveKey(key); err != nil {
		return errors.New("key not found") // Key not found, successful no-op.
	}

	// If the leaf is the root and is now empty, the tree is empty.
	if leaf == b.root && leaf.GetNumberOfKeys() == 0 {
		b.root = nil
		return nil
	}

	// If the node is still valid (not in underflow), we are done.
	// Note: A more robust implementation would update ancestor keys if the first key of the leaf was deleted.
	if leaf.GetNumberOfKeys() >= b.minKeys {
		return nil
	}

	// The leaf is in underflow. We need to rebalance.
	return b.handleUnderflow(leaf, parent)
}

// handleUnderflow manages rebalancing a node that has too few keys.
func (b *BPlusTree[K, V]) handleUnderflow(node, parent *Node[K, V]) error {
	// If the underflow node is the root, special handling is needed.
	if node == b.root {
		// If the root is an internal node with only one child, that child becomes the new root.
		if !node.IsLeaf() && node.GetNumberOfChildren() == 1 {
			b.root, _ = node.GetChildAtIndex(0)
		}
		return nil
	}

	childIndex, err := b.findChildIndex(parent, node)
	if err != nil {
		return err
	}

	// Try to redistribute with the right sibling.
	if childIndex < parent.GetNumberOfChildren()-1 {
		rightSibling, _ := parent.GetChildAtIndex(childIndex + 1)
		if rightSibling.GetNumberOfKeys() > b.minKeys {
			b.redistributeFromRight(node, rightSibling, parent, childIndex)
			return nil
		}
	}

	// Try to redistribute with the left sibling.
	if childIndex > 0 {
		leftSibling, _ := parent.GetChildAtIndex(childIndex - 1)
		if leftSibling.GetNumberOfKeys() > b.minKeys {
			b.redistributeFromLeft(node, leftSibling, parent, childIndex-1)
			return nil
		}
	}

	// If redistribution is not possible, merge.
	if childIndex < parent.GetNumberOfChildren()-1 {
		// Merge with the right sibling.
		rightSibling, _ := parent.GetChildAtIndex(childIndex + 1)
		return b.merge(node, rightSibling, parent, childIndex)
	}
	// Merge with the left sibling.
	leftSibling, _ := parent.GetChildAtIndex(childIndex - 1)
	return b.merge(leftSibling, node, parent, childIndex-1)
}

// redistributeFromRight borrows a key from the right sibling.
func (b *BPlusTree[K, V]) redistributeFromRight(node, rightSibling, parent *Node[K, V], parentKeyIndex int) {
	keyToMove := rightSibling.GetFirstKey()

	if node.IsLeaf() {
		valueToMove, _ := rightSibling.GetValueAtIndex(0)
		node.AddKeyValue(keyToMove, valueToMove)
		rightSibling.RemoveKey(keyToMove)
		parent.keys.Data[parentKeyIndex] = rightSibling.GetFirstKey()
	} else {
		separatorKey := parent.GetKeys().Get(parentKeyIndex)
		node.AddKey(separatorKey)
		childToMove, _ := rightSibling.GetChildAtIndex(0)
		rightSibling.RemoveChildAtIndex(0)
		node.AddChild(childToMove)
		rightSibling.RemoveKey(keyToMove)
		parent.keys.Data[parentKeyIndex] = keyToMove
	}
}

// redistributeFromLeft borrows a key from the left sibling.
func (b *BPlusTree[K, V]) redistributeFromLeft(node, leftSibling, parent *Node[K, V], parentKeyIndex int) {
	if node.IsLeaf() {
		lastKeyIndex := leftSibling.GetNumberOfKeys() - 1
		keyToMove := leftSibling.GetIterableKeys()[lastKeyIndex]
		valueToMove, _ := leftSibling.GetValueAtIndex(lastKeyIndex)
		node.AddKeyValue(keyToMove, valueToMove)
		leftSibling.RemoveKey(keyToMove)
		parent.keys.Data[parentKeyIndex] = node.GetFirstKey()
	} else {
		separatorKey := parent.GetKeys().Get(parentKeyIndex)
		node.AddKey(separatorKey)
		childToMove, _ := leftSibling.GetChildAtIndex(leftSibling.GetNumberOfChildren() - 1)
		leftSibling.RemoveChildAtIndex(leftSibling.GetNumberOfChildren() - 1)
		node.children = slices.Insert(node.children, 0, childToMove)
		lastKeyIndex := leftSibling.GetNumberOfKeys() - 1
		keyToMoveUp := leftSibling.GetIterableKeys()[lastKeyIndex]
		leftSibling.RemoveKey(keyToMoveUp)
		parent.keys.Data[parentKeyIndex] = keyToMoveUp
	}
}

// merge combines an underflow node with an adjacent sibling.
func (b *BPlusTree[K, V]) merge(leftNode, rightNode, parent *Node[K, V], parentKeyIndex int) error {
	separatorKey := parent.GetKeys().Get(parentKeyIndex)

	if leftNode.IsLeaf() {
		// Move all keys and values from right to left.
		for i, k := range rightNode.GetIterableKeys() {
			v, _ := rightNode.GetValueAtIndex(i)
			leftNode.AddKeyValue(k, v)
		}
		leftNode.next = rightNode.next
	} else {
		leftNode.AddKey(separatorKey)
		for _, k := range rightNode.GetIterableKeys() {
			leftNode.AddKey(k)
		}
		for _, c := range rightNode.GetChildren() {
			leftNode.AddChild(c)
		}
	}

	// Remove the key and pointer from the parent, which might cause a recursive underflow.
	parent.RemoveKey(separatorKey)
	parent.RemoveChildAtIndex(parentKeyIndex + 1)

	if parent.GetNumberOfKeys() < b.minKeys {
		grandParent, _ := b.findParent(b.root, parent)
		return b.handleUnderflow(parent, grandParent)
	}

	return nil
}

// findChildIndex finds the index of a child within a parent's children slice.
func (b *BPlusTree[K, V]) findChildIndex(parent, child *Node[K, V]) (int, error) {
	if parent == nil || parent.IsLeaf() {
		return -1, fmt.Errorf("cannot find child index in a nil or leaf parent")
	}
	for i, c := range parent.GetChildren() {
		if c == child {
			return i, nil
		}
	}
	return -1, fmt.Errorf("child not found in parent")
}

// findParent is a helper function to find the parent of a given node.
// This is inefficient and is a consequence of not having parent pointers.
// It starts searching from startNode.
func (b *BPlusTree[K, V]) findParent(startNode, childNode *Node[K, V]) (*Node[K, V], error) {
	if startNode == nil || startNode.IsLeaf() {
		return nil, nil
	}

	for _, child := range startNode.GetChildren() {
		if child == childNode {
			return startNode, nil
		}
	}

	// Figure out which child's subtree to search in.
	key := childNode.GetFirstKey()
	keys := startNode.GetIterableKeys()
	children := startNode.GetChildren()

	// Find the correct child to descend into.
	for i := 0; i < len(keys); i++ {
		if key < keys[i] {
			if i < len(children) {
				return b.findParent(children[i], childNode)
			}
			return nil, fmt.Errorf("child index out of bounds during parent search")
		}
	}

	// If key is greater than all keys, descend into the rightmost child.
	if len(keys) < len(children) {
		return b.findParent(children[len(keys)], childNode)
	}

	return nil, nil
}

// Search - search for a node in the tree based on the value
// bool represents whether or not the node has been found
// *Node[K, V] is the node if it is found. Otherwise, it is nil
// V is the value if found
func (b *BPlusTree[K, V]) Search(searchKey K) (V, bool, error) {
	var zeroV V
	rootNode, err := b.GetRootNode()
	if err != nil {
		return zeroV, false, err
	}

	// we'll want to perform a search. We'll need to start
	// from the root and go from.
	leaf, _, err := b.leafSearch(searchKey, rootNode, nil)
	if err != nil {
		return zeroV, false, err
	}

	// look for the key within the leafs.
	for i, k := range leaf.GetIterableKeys() {
		// if the key is found, get the value and return it
		if k == searchKey {
			if value, err := leaf.GetValueAtIndex(i); err == nil {
				return value, true, nil
			}
		}
	}

	return zeroV, false, nil
}

// leafSearch returns the leaf node in question
// *Node[K,V] - leaf node
// *Node[K,V] - parent node
// error - error if there is any
func (b *BPlusTree[K, V]) leafSearch(searchKey K, currentNode *Node[K, V], parentNode *Node[K, V]) (*Node[K, V], *Node[K, V], error) {
	// if the current node is a leaf, return the node and its parent
	if currentNode.IsLeaf() {
		return currentNode, parentNode, nil
	}

	// get all of the children of this node
	children := currentNode.GetIterableChildren()

	// retrieve a list of keys that can be considered left sided
	// intervals. these are are for the B+ tree so that we can search
	// the child nodes to the left of these keys
	leftSidedIntervals := currentNode.GetIterableKeys()

	if len(children) != (len(leftSidedIntervals) + 1) {
		return nil, nil, fmt.Errorf("the alignment of keys to children is incorrect. there's an issue with the table.")
	}

	m := len(leftSidedIntervals)

	// loop through the keys, if you find a key that the searchKey
	// is less than or equal to. If we find that key, we know that we should look
	// at the children of the current node
	for i := 0; i < m; i++ {
		if searchKey <= leftSidedIntervals[i] {
			searchNode, err := currentNode.GetChildAtIndex(i)

			if err != nil {
				return nil, nil, fmt.Errorf("error during the search for the node: %w", err)
			}

			return b.leafSearch(searchKey, searchNode, currentNode)
		}
	}

	// If searchKey is greater than all keys, traverse the rightmost child.
	searchNode, err := currentNode.GetChildAtIndex(m)
	if err != nil {
		return nil, nil, fmt.Errorf("error during search for rightmost node: %w", err)
	}
	return b.leafSearch(searchKey, searchNode, currentNode)
}
