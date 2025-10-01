package bplustree

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewNode(t *testing.T) {
	t.Run("Create leaf node", func(t *testing.T) {
		node := NewNode[int, string](10, "ten", true)
		if node == nil {
			t.Fatal("NewNode returned nil")
		}
		if !node.IsLeaf() {
			t.Error("Expected node to be a leaf")
		}
		if node.GetFirstKey() != 10 {
			t.Errorf("Expected first key to be 10, got %d", node.GetFirstKey())
		}
		if len(node.GetChildren()) != 0 {
			t.Error("Expected new leaf node to have no children")
		}
	})

	t.Run("Create internal node", func(t *testing.T) {
		node := NewNode[string, int]("test", 0, false)
		if node == nil {
			t.Fatal("NewNode returned nil")
		}
		if node.IsLeaf() {
			t.Error("Expected node to not be a leaf")
		}
		if node.GetFirstKey() != "test" {
			t.Errorf("Expected first key to be 'test', got %s", node.GetFirstKey())
		}
		if node.GetChildren() != nil {
			t.Error("Expected new internal node to have nil children slice")
		}
	})
}

func TestNodeKeys(t *testing.T) {
	t.Run("GetFirstKey", func(t *testing.T) {
		node := NewNode[int, string](20, "twenty", true)
		if key := node.GetFirstKey(); key != 20 {
			t.Errorf("Expected first key to be 20, got %d", key)
		}
	})

	t.Run("GetFirstKey on empty node", func(t *testing.T) {
		emptyNode := &Node[int, string]{}
		if key := emptyNode.GetFirstKey(); key != 0 {
			t.Errorf("Expected first key of empty node to be zero value, got %d", key)
		}
	})

	t.Run("AddKey and GetKeys", func(t *testing.T) {
		node := NewNode[int, string](20, "twenty", true)
		node.AddKey(10)
		node.AddKey(30)
		keys := node.GetKeys().GetData()
		expected := []int{10, 20, 30}
		if !reflect.DeepEqual(keys, expected) {
			t.Errorf("Expected keys to be %v, got %v", expected, keys)
		}
	})

	t.Run("GetIterableKeys", func(t *testing.T) {
		node := NewNode[int, string](20, "twenty", true)
		node.AddKey(10)
		node.AddKey(30)
		keys := node.GetIterableKeys()
		expected := []int{10, 20, 30}
		if !reflect.DeepEqual(keys, expected) {
			t.Errorf("Expected iterable keys to be %v, got %v", expected, keys)
		}
	})

	t.Run("RemoveKey", func(t *testing.T) {
		node := NewNode[int, string](20, "twenty", true)
		node.AddKey(10)
		node.AddKey(30) // Keys are now [10, 20, 30]

		// Positive test: remove an existing key from the middle
		err := node.RemoveKey(20)
		if err != nil {
			t.Fatalf("RemoveKey(20) returned an unexpected error: %v", err)
		}
		expectedAfterRemove1 := []int{10, 30}
		if !reflect.DeepEqual(node.GetIterableKeys(), expectedAfterRemove1) {
			t.Errorf("Expected keys to be %v after removing 20, got %v", expectedAfterRemove1, node.GetIterableKeys())
		}

		// Negative test: remove a non-existent key
		err = node.RemoveKey(99)
		if err == nil {
			t.Fatal("Expected an error when removing a non-existent key, but got nil")
		}
		// Ensure keys are unchanged after failed removal
		if !reflect.DeepEqual(node.GetIterableKeys(), expectedAfterRemove1) {
			t.Errorf("Keys should not have changed after attempting to remove a non-existent key. Expected %v, got %v", expectedAfterRemove1, node.GetIterableKeys())
		}

		// Positive test: remove the first key
		err = node.RemoveKey(10)
		if err != nil {
			t.Fatalf("RemoveKey(10) returned an unexpected error: %v", err)
		}
		expectedAfterRemove2 := []int{30}
		if !reflect.DeepEqual(node.GetIterableKeys(), expectedAfterRemove2) {
			t.Errorf("Expected keys to be %v after removing 10, got %v", expectedAfterRemove2, node.GetIterableKeys())
		}

		// Positive test: remove the last remaining key
		err = node.RemoveKey(30)
		if err != nil {
			t.Fatalf("RemoveKey(30) returned an unexpected error: %v", err)
		}
		if len(node.GetIterableKeys()) != 0 {
			t.Errorf("Expected keys to be empty after removing the last key, got %v", node.GetIterableKeys())
		}

		// Negative test: remove from a node with no keys left
		err = node.RemoveKey(100)
		if err == nil {
			t.Fatal("Expected an error when removing from a node with no keys left, but got nil")
		}
	})
}

func TestNodeChildren(t *testing.T) {
	parent := NewNode[int, string](100, "one-hundred", false)
	child1 := NewNode[int, string](50, "fifty", true)
	child2 := NewNode[int, string](200, "two-hundred", true)
	child3 := NewNode[int, string](150, "one-hundred-fifty", true)

	t.Run("Add and Get Children", func(t *testing.T) {
		parent.AddChild(child1)
		parent.AddChild(child2)
		parent.AddChild(child3) // Insert between child1 and child2

		children := parent.GetIterableChildren()
		if len(children) != 3 {
			t.Fatalf("Expected 3 children, got %d", len(children))
		}
		if children[0] != child1 || children[1] != child3 || children[2] != child2 {
			t.Errorf("Children not in expected order")
		}
	})

	t.Run("Add child out of bounds", func(t *testing.T) {
		node := NewNode[int, string](10, "ten", false)
		child := NewNode[int, string](5, "five", true)
		node.AddChild(child) // Index is out of bounds
		if len(node.GetChildren()) != 1 || node.GetChildren()[0] != child {
			t.Error("Expected child to be appended when index is out of bounds")
		}
	})

	t.Run("GetChildAtIndex", func(t *testing.T) {
		child, err := parent.GetChildAtIndex(1)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if child != child3 {
			t.Error("Got incorrect child at index")
		}
	})

	t.Run("GetChildAtIndex negative index", func(t *testing.T) {
		_, err := parent.GetChildAtIndex(-1)
		if err == nil {
			t.Error("Expected error for negative index, got nil")
		}
	})

	t.Run("GetChildAtIndex out of bounds", func(t *testing.T) {
		_, err := parent.GetChildAtIndex(len(parent.GetChildren()))
		if err == nil {
			t.Error("Expected error for out of bounds index, got nil")
		}
	})

	t.Run("GetChildAtIndex on node with no children", func(t *testing.T) {
		node := NewNode[int, string](1, "one", false)
		_, err := node.GetChildAtIndex(0)
		if err == nil {
			t.Error("Expected error when getting child from node with no children, got nil")
		}
	})
}

func TestNodeNavigation(t *testing.T) {
	leaf1 := NewNode[int, string](10, "ten", true)
	leaf2 := NewNode[int, string](20, "twenty", true)
	internalNode := NewNode[int, string](30, "thirty", false)
	leaf1.next = leaf2

	t.Run("GetNextNode from leaf", func(t *testing.T) {
		next, err := leaf1.GetNextNode()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if next != leaf2 {
			t.Error("Next node was not the expected node")
		}
	})

	t.Run("GetNextNode from leaf with no next", func(t *testing.T) {
		next, err := leaf2.GetNextNode()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if next != nil {
			t.Error("Expected nil for next node")
		}
	})

	t.Run("GetNextNode from internal node", func(t *testing.T) {
		_, err := internalNode.GetNextNode()
		if err == nil {
			t.Error("Expected error when getting next node from internal node, got nil")
		}
	})
}

func TestIsLeaf(t *testing.T) {
	leaf := NewNode[int, string](10, "ten", true)
	internal := NewNode[int, string](20, "twenty", false)

	if !leaf.IsLeaf() {
		t.Error("Expected node to be a leaf")
	}

	if internal.IsLeaf() {
		t.Error("Expected node to not be a leaf")
	}
}

func TestString(t *testing.T) {
	node := NewNode[int, string](10, "ten", false)
	expected := fmt.Sprintf("%v: %v", node.keys, node.children)
	if node.String() != expected {
		t.Errorf("Expected string representation to be '%s', got '%s'", expected, node.String())
	}
}
