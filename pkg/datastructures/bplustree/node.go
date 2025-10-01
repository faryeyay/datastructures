// Package bplustree
// Farye Nwede <farye@aeekay.com>
package bplustree

import (
	"cmp"
	"errors"
	"fmt"
	"slices"

	"github.com/faryeyay/insights/primitives/pkg/datastructures/orderedlist"
)

// Node - Represents a node in a B+ Tree
// The keys should be a comparable value of K. We enforced this using
// cmp.Ordered.
// The data (otherwise known as the value) should be V. Upon further though
// we don't need this to be a comparable.
type Node[K cmp.Ordered, V any] struct {
	keys     *orderedlist.OrderedList[K] // Represents the keys of the Node
	children []*Node[K, V]
	values   []V // values associated with keys, for leaf nodes only
	leaf     bool
	next     *Node[K, V]
	// leaf - whether or not the node is a leaf node. If the node is a leaf,
	// then it should have a value and no children. If not, then the
	// children can be populated but value must be nil

	// next - next points to the next leaf node. This allows us to traverse
	// leaf nodes
}

// NewNode - Creates a new node with the key
func NewNode[K cmp.Ordered, V any](key K, value V, leaf bool) *Node[K, V] {
	// Create an empty map of children
	// TODO: Maybe avoid creating the map if the
	// node is a leaf
	keys := &orderedlist.OrderedList[K]{}
	keys.Insert(key)

	node := &Node[K, V]{
		keys: keys,
		leaf: leaf,
	}
	if leaf {
		node.values = []V{value}
	}
	return node
}

// NewEmptyNode - Creates a new empty node with the key
func NewEmptyNode[K cmp.Ordered, V any]() *Node[K, V] {
	return &Node[K, V]{
		keys:     &orderedlist.OrderedList[K]{},
		leaf:     false,
		children: []*Node[K, V]{},
	}
}

// GetKey - Returns the key of the B+ Tree's node
func (n *Node[K, V]) GetFirstKey() K {
	var k K

	// return an empty K since there are no keys
	if n.keys == nil {
		return k
	}

	// return the key at the index 0
	// this is the first element of the ordered list
	if len(n.keys.GetData()) == 0 {
		return k
	}

	k = n.keys.Get(0)
	return k
}

// String - Returns a string representation of the B+ Tree node
func (n *Node[K, V]) String() string {
	return fmt.Sprintf("%v: %v", n.keys, n.children)
}

// GetKeys - Returns the keys of a B+ Tree node.
func (n *Node[K, V]) GetKeys() *orderedlist.OrderedList[K] {
	return n.keys
}

// GetIterableKeys - Returns the keys of a B+ Tree that are iterable
// this allows for loop controls for the data structure
func (n *Node[K, V]) GetIterableKeys() []K {
	return n.keys.GetData()
}

// AddKey adds a key to an internal node.
func (n *Node[K, V]) AddKey(k K) {
	n.keys.Insert(k)
}

// AddKeyValue adds a key-value pair to a leaf node.
// If the key already exists, it updates the value.
func (n *Node[K, V]) AddKeyValue(k K, v V) {
	idx, inserted := n.keys.Insert(k)
	if inserted {
		// New key, insert value at the same index.
		n.values = slices.Insert(n.values, idx, v)
	} else {
		// Key already exists, update the value.
		n.values[idx] = v
	}
}

// RemoveKey - Removes a key from the B+ Tree node.
func (n *Node[K, V]) RemoveKey(k K) error {
	// Before removing the key, find its index to remove the corresponding value.
	idx, found := slices.BinarySearch(n.GetIterableKeys(), k)

	if !found {
		return errors.New("key not found in node")
	}

	if n.leaf && len(n.values) <= 1 {
		// clear the slice
		n.values = []V{}
	}

	if n.leaf && len(n.values) > 1 {
		n.values = slices.Delete(n.values, idx, idx+1)
	}

	return n.keys.Remove(k)
}

// GetChildren - Returns the children of a B+ Tree node.
func (n *Node[K, V]) GetChildren() []*Node[K, V] {
	// handle the case where children is not set
	if n.children == nil {
		return nil
	}

	return n.children
}

// GetValueAtIndex returns the value at a specific index. For leaf nodes only.
func (n *Node[K, V]) GetValueAtIndex(idx int) (V, error) {
	var zeroV V
	if !n.leaf {
		return zeroV, errors.New("cannot get value from an internal node")
	}
	if idx < 0 || idx >= len(n.values) {
		return zeroV, errors.New("index out of bounds for values")
	}
	return n.values[idx], nil
}

// GetValues returns all values from a leaf node.
func (n *Node[K, V]) GetValues() []V {
	if !n.leaf {
		return nil
	}
	return n.values
}

// GetIterableChildren - Returns the children of a B+ Tree node. This
// allows for the iteration of children.
func (n *Node[K, V]) GetIterableChildren() []*Node[K, V] {
	// handle the case where children is not set
	if n.children == nil {
		return nil
	}

	return n.children
}

// AddChild adds a child to the node's children, maintaining sorted order by the child's first key.
// It uses a binary search to find the insertion point in O(log n) time. The insertion
// into the slice is an O(n) operation.
func (n *Node[K, V]) AddChild(node *Node[K, V]) error {
	// Leaf nodes cannot have children.
	if n.leaf {
		return errors.New("leaf nodes can't have children")
	}

	// If children slice is nil, initialize it.
	if n.children == nil {
		n.children = []*Node[K, V]{}
	}

	// Find the correct insertion index using binary search (O(log n)).
	// The search is based on the first key of the child nodes.
	idx, _ := slices.BinarySearchFunc(n.children, node.GetFirstKey(), func(a *Node[K, V], b K) int {
		return cmp.Compare(a.GetFirstKey(), b)
	})

	// Insert the new child at the correct position to maintain order (O(n)).
	n.children = slices.Insert(n.children, idx, node)

	return nil
}

// RemoveChildAtIndex - removes a child from the childrens slice at a given index
func (n *Node[K, V]) RemoveChildAtIndex(idx int) {
	n.children = slices.Delete(n.children, idx, idx+1)
}

// GetChildAtIndex - returns the child at the specific index
func (n *Node[K, V]) GetChildAtIndex(idx int) (*Node[K, V], error) {
	if n.children == nil {
		return nil, fmt.Errorf("there are no children for this node")
	}

	// check for indexes that are out bounds. if the index is out of bounds
	// return an error
	if idx < 0 || len(n.children) <= idx {
		return nil, fmt.Errorf("the index out of range")
	}

	return n.children[idx], nil
}

func (n *Node[K, V]) GetNumberOfKeys() int {
	// handle the case when keys is not set
	if n.keys == nil {
		return 0
	}

	// else return the length of the slice of keys
	return len(n.keys.GetData())
}

// GetNextNode - Returns the next node. If the node is not a leaf,
// return nil and an error. Otherwise, return the next node.
func (n *Node[K, V]) GetNextNode() (*Node[K, V], error) {
	if n.leaf {
		return n.next, nil
	}

	return nil, fmt.Errorf("the node %v is not a leaf node", n)
}

// IsLeaf - Returns whether the node is a leaf or not
func (n *Node[K, V]) IsLeaf() bool {
	return n.leaf
}

func (n *Node[K, V]) GetNumberOfChildren() int {
	// handle the case when children is not set
	if n.children == nil {
		return 0
	}

	// else return the length of the slice of children
	return len(n.children)
}
