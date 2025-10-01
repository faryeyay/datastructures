// Package bplustree - Test-Driven Development test file for BPlusTree
package bplustree

import (
	"reflect"
	"testing"
)

func TestCreateTree(t *testing.T) {
	t.Run("Create empty tree", func(t *testing.T) {
		tree := CreateTree[int, string](3)
		if tree == nil {
			t.Fatal("CreateTree returned nil")
		}
		if tree.root != nil {
			t.Error("Expected new tree to have a nil root")
		}
	})
}

func TestCreateTreeWithRootKey(t *testing.T) {
	t.Run("Create tree with root key", func(t *testing.T) {
		tree := CreateTreeWithRootKey[int, string](3, 10)
		if tree == nil {
			t.Fatal("CreateTreeWithRootKey returned nil tree")
		}
		if tree.root == nil {
			t.Fatal("Expected root to be initialized")
		}
		if tree.root.GetFirstKey() != 10 {
			t.Errorf("Expected root key to be 10, got %d", tree.root.GetFirstKey())
		}
		if tree.root.IsLeaf() { // CreateTreeWithRootKey creates an internal node by default
			t.Error("Expected root created with key to be an internal node, not a leaf")
		}
	})
}

func TestGetRootNode(t *testing.T) {
	t.Run("Get root from initialized tree", func(t *testing.T) {
		tree := CreateTreeWithRootKey[int, string](3, 10)
		root, err := tree.GetRootNode()
		if err != nil {
			t.Fatalf("GetRootNode returned an unexpected error: %v", err)
		}
		if root == nil {
			t.Fatal("GetRootNode returned a nil root")
		}
		if root.GetFirstKey() != 10 {
			t.Errorf("Expected root key to be 10, got %d", root.GetFirstKey())
		}
	})

	t.Run("Get root from empty tree", func(t *testing.T) {
		tree := CreateTree[int, string](3)
		_, err := tree.GetRootNode()
		if err == nil {
			t.Error("Expected an error when getting root from an empty tree, but got nil")
		}
	})
}

// --- Test-Driven Development for Insert ---
// The following tests are expected to fail until the Insert function is implemented.

func TestInsert(t *testing.T) {
	t.Run("Insert into an empty tree", func(t *testing.T) {
		tree := CreateTree[int, string](3)
		tree.Insert(10, "ten")
		root, err := tree.GetRootNode()
		if err != nil {
			t.Fatalf("Failed to get root after insert: %v", err)
		}
		if !root.IsLeaf() {
			t.Error("Expected root of a single-node tree to be a leaf")
		}
		if root.GetNumberOfKeys() != 1 || root.GetFirstKey() != 10 {
			t.Errorf("Expected root to contain key 10, got keys: %v", root.GetIterableKeys())
		}
		val, _ := root.GetValueAtIndex(0)
		if val != "ten" {
			t.Errorf("Expected root value to be 'ten', got '%s'", val)
		}
	})

	t.Run("Insert multiple keys without split", func(t *testing.T) {
		tree := CreateTree[int, string](3) // maxKeys = 5
		tree.Insert(10, "ten")
		tree.Insert(20, "twenty")
		tree.Insert(5, "five")
		root, _ := tree.GetRootNode()
		keys := root.GetIterableKeys()
		expected := []int{5, 10, 20}
		// This requires reflect.DeepEqual, but for now, we'll check length and order.
		if len(keys) != 3 {
			t.Fatalf("Expected 3 keys, got %d", len(keys))
		}
		if keys[0] != 5 || keys[1] != 10 || keys[2] != 20 {
			t.Errorf("Keys not in correct order. Got %v, expected %v", keys, expected)
		}
		values := root.GetValues()
		expectedValues := []string{"five", "ten", "twenty"}
		if !reflect.DeepEqual(values, expectedValues) {
			t.Errorf("Values not in correct order. Got %v, expected %v", values, expectedValues)
		}
	})

	t.Run("Insert causing a leaf node split", func(t *testing.T) {
		// Using order=2 means maxKeys=3. The 4th key will cause a split.
		tree := CreateTree[int, string](2)
		tree.Insert(10, "ten")
		tree.Insert(20, "twenty")
		tree.Insert(30, "thirty")
		tree.Insert(5, "five") // This should cause a split

		root, _ := tree.GetRootNode()
		if root.IsLeaf() {
			t.Fatal("Root should be an internal node after a split")
		}
		if root.GetNumberOfKeys() != 1 {
			t.Fatalf("Expected root to have 1 key after split, got %d", root.GetNumberOfKeys())
		}
		// With leaf [5, 10, 20, 30], split is at index 2. Promoted key is 20.
		// Left leaf [5, 10], right leaf [20, 30].
		if root.GetFirstKey() != 20 {
			t.Errorf("Expected root key to be 20, got %d", root.GetFirstKey())
		}
		if root.GetNumberOfChildren() != 2 {
			t.Fatalf("Expected root to have 2 children, got %d", root.GetNumberOfChildren())
		}
		child1, _ := root.GetChildAtIndex(0)
		child2, _ := root.GetChildAtIndex(1)
		if !reflect.DeepEqual(child1.GetIterableKeys(), []int{5, 10}) {
			t.Errorf("Expected first child's keys to be [5, 10], got %v", child1.GetIterableKeys())
		}
		if !reflect.DeepEqual(child1.GetValues(), []string{"five", "ten"}) {
			t.Errorf("Expected first child's values to be [five, ten], got %v", child1.GetValues())
		}
		if !reflect.DeepEqual(child2.GetIterableKeys(), []int{20, 30}) {
			t.Errorf("Expected second child's keys to be [20, 30], got %v", child2.GetIterableKeys())
		}
		if !reflect.DeepEqual(child2.GetValues(), []string{"twenty", "thirty"}) {
			t.Errorf("Expected second child's values to be [twenty, thirty], got %v", child2.GetValues())
		}
	})

	t.Run("Insert causing multi-level split (root split)", func(t *testing.T) {
		// order=2 means maxKeys=3. 4th key causes a split.
		tree := CreateTree[int, string](2)

		// These insertions will cause a series of leaf and internal node splits,
		// culminating in a root split that increases the tree's height.
		keysToInsert := []int{10, 20, 30, 40, 5, 15, 25, 35, 28, 50}
		valuesToInsert := []string{"ten", "twenty", "thirty", "forty", "five", "fifteen", "twentyfive", "thirtyfive", "twentyeight", "fifty"}
		for i, key := range keysToInsert {
			err := tree.Insert(key, valuesToInsert[i])
			if err != nil {
				t.Fatalf("Insert(%d, \"val\") failed: %v", key, err)
			}
		}
		// Note: This test doesn't use values, so we're just checking structure.

		// --- Verify final tree structure after root split ---

		// After the final insert (50), the root [15, 25, 30] splits.
		// The key '30' is promoted to a new root.
		root, _ := tree.GetRootNode()
		if root.IsLeaf() {
			t.Fatal("Root should be an internal node")
		}
		if !reflect.DeepEqual(root.GetIterableKeys(), []int{30}) {
			t.Fatalf("Expected new root key to be [30], got %v", root.GetIterableKeys())
		}
		if root.GetNumberOfChildren() != 2 {
			t.Fatalf("Expected new root to have 2 children, got %d", root.GetNumberOfChildren())
		}

		// Check the children of the new root
		child1, _ := root.GetChildAtIndex(0) // This is the old root
		child2, _ := root.GetChildAtIndex(1) // This is the new internal node

		expectedChild1Keys := []int{15, 25}
		if !reflect.DeepEqual(child1.GetIterableKeys(), expectedChild1Keys) {
			t.Errorf("Expected root's first child keys to be %v, got %v", expectedChild1Keys, child1.GetIterableKeys())
		}
		if child1.GetNumberOfChildren() != 3 {
			t.Errorf("Expected root's first child to have 3 children, got %d", child1.GetNumberOfChildren())
		}

		expectedChild2Keys := []int{40}
		if !reflect.DeepEqual(child2.GetIterableKeys(), expectedChild2Keys) {
			t.Errorf("Expected root's second child keys to be %v, got %v", expectedChild2Keys, child2.GetIterableKeys())
		}
		if child2.GetNumberOfChildren() != 2 {
			t.Errorf("Expected root's second child to have 2 children, got %d", child2.GetNumberOfChildren())
		}

		// Verify a leaf node to ensure data is in the right place
		// Search for 28, which should be in a leaf under the first child of the root.
		c1, _ := root.GetChildAtIndex(0)
		c1c3, _ := c1.GetChildAtIndex(2) // Child of [15, 25] for keys >= 25

		expectedLeafKeys := []int{25, 28}
		if !reflect.DeepEqual(c1c3.GetIterableKeys(), expectedLeafKeys) {
			t.Errorf("Expected a specific leaf to have keys %v, got %v", expectedLeafKeys, c1c3.GetIterableKeys())
		}

		// Verify another leaf
		c2, _ := root.GetChildAtIndex(1)
		c2c2, _ := c2.GetChildAtIndex(1) // Child of [40] for keys >= 40
		expectedLeafKeys2 := []int{40, 50}
		if !reflect.DeepEqual(c2c2.GetIterableKeys(), expectedLeafKeys2) {
			t.Errorf("Expected a specific leaf to have keys %v, got %v", expectedLeafKeys2, c2c2.GetIterableKeys())
		}
	})
}

// --- Test-Driven Development for Search ---
// The following tests may fail until Search and Insert are fully implemented.

func TestSearch(t *testing.T) {
	// Setup a tree for searching
	// Use Insert to build a predictable tree structure.
	// With order=2 (maxKeys=3), inserting 4 keys will cause a split
	tree := CreateTree[int, string](2)
	tree.Insert(10, "ten")
	tree.Insert(30, "thirty")
	tree.Insert(5, "five")
	tree.Insert(20, "twenty") // This insert causes a split.
	// The tree should look like this:
	//      [20]
	//     /    \
	// [5, 10]  [20, 30] (with values)

	t.Run("Search in an empty tree", func(t *testing.T) {
		emptyTree := CreateTree[int, string](3)
		_, found, err := emptyTree.Search(10)
		if err == nil {
			t.Error("Expected error when searching in an empty tree, got nil")
		}
		if found {
			t.Error("Expected not to find key in an empty tree")
		}
	})

	t.Run("Search for existing key in leaf", func(t *testing.T) {
		val, found, err := tree.Search(10)
		if err != nil {
			t.Fatalf("Unexpected error during search: %v", err)
		}
		if !found {
			t.Error("Failed to find existing key 10")
		}
		if val != "ten" {
			t.Errorf("Search for key 10 returned wrong value. Got '%s', want 'ten'", val)
		}

		val, found, err = tree.Search(30)
		if err != nil {
			t.Fatalf("Unexpected error during search: %v", err)
		}
		if !found {
			t.Error("Failed to find existing key 30")
		}
		if val != "thirty" {
			t.Errorf("Search for key 30 returned wrong value. Got '%s', want 'thirty'", val)
		}
	})

	t.Run("Search for non-existent key", func(t *testing.T) {
		_, found, err := tree.Search(99)
		if err != nil {
			t.Fatalf("Unexpected error during search for non-existent key: %v", err)
		}
		if found {
			t.Error("Found key 99, which should not exist")
		}
	})
}

// --- Test-Driven Development for Delete ---
// The following tests are expected to fail until a Delete/Remove function is implemented.

func TestDelete(t *testing.T) {
	t.Run("Delete non-existent key", func(t *testing.T) {
		tree := CreateTree[int, string](2)
		tree.Insert(10, "ten")
		tree.Insert(20, "twenty")
		err := tree.Delete(99)
		if err == nil {
			t.Errorf("Deleting a non-existent key should produce an error, got no error")
		}
		expectedKeys := []int{10, 20}
		if !reflect.DeepEqual(tree.root.GetIterableKeys(), expectedKeys) {
			t.Errorf("Tree should be unchanged. Expected %v, got %v", expectedKeys, tree.root.GetIterableKeys())
		}
	})

	t.Run("Delete from leaf without underflow", func(t *testing.T) {
		// order=2 means minKeys=1, maxKeys=3
		tree := CreateTree[int, string](2)
		tree.Insert(10, "ten")
		tree.Insert(20, "twenty") // Leaf has [10, 20], count=2 > minKeys
		tree.Delete(10)
		root, _ := tree.GetRootNode()
		expectedKeys := []int{20}
		if !reflect.DeepEqual(root.GetIterableKeys(), expectedKeys) {
			t.Errorf("Deletion from leaf failed. Expected %v, got %v", expectedKeys, root.GetIterableKeys())
		}
	})

	t.Run("Delete causing re-distribution (borrow from right)", func(t *testing.T) {
		// order=2 means minKeys=1, maxKeys=3
		tree := CreateTree[int, string](2)
		tree.Insert(10, "ten")
		tree.Insert(20, "twenty")
		tree.Insert(30, "thirty")
		tree.Insert(40, "forty") // Causes split -> Root [30], Left [10, 20], Right [30, 40]
		tree.Delete(10)          // Left leaf is now [20], which is at minKeys.

		// Now delete 20, causing underflow. Left leaf will borrow from right.
		tree.Delete(20)

		root, _ := tree.GetRootNode()
		// After redistribution, parent key is updated.
		// Right sibling [30, 40] gives 30 to left leaf. Parent key becomes 40.
		// Left leaf becomes [30], right leaf becomes [40].
		expectedRootKeys := []int{40}
		if !reflect.DeepEqual(root.GetIterableKeys(), expectedRootKeys) {
			t.Fatalf("Expected root keys to be %v, got %v", expectedRootKeys, root.GetIterableKeys())
		}

		left, _ := root.GetChildAtIndex(0)
		right, _ := root.GetChildAtIndex(1)
		expectedLeftKeys := []int{30}
		expectedRightKeys := []int{40}
		if !reflect.DeepEqual(left.GetIterableKeys(), expectedLeftKeys) {
			t.Errorf("Expected left leaf to be %v, got %v", expectedLeftKeys, left.GetIterableKeys())
		}
		if !reflect.DeepEqual(right.GetIterableKeys(), expectedRightKeys) {
			t.Errorf("Expected right leaf to be %v, got %v", expectedRightKeys, right.GetIterableKeys())
		}
	})

	t.Run("Delete causing merge and tree height decrease", func(t *testing.T) {
		// order=2 means minKeys=1, maxKeys=3
		tree := CreateTree[int, string](2)
		tree.Insert(10, "ten")
		tree.Insert(20, "twenty")
		tree.Insert(30, "thirty")
		tree.Insert(5, "five") // Causes split -> Root [20], Left [5, 10], Right [20, 30]

		// Delete to get both children to minKeys
		tree.Delete(5)  // Left is now [10]
		tree.Delete(30) // Right is now [20]

		// Now, deleting 10 will cause left to underflow and merge with right
		tree.Delete(10)

		// The merge causes the parent (root) to underflow and be removed.
		// The merged leaf becomes the new root, decreasing tree height.
		root, _ := tree.GetRootNode()
		if !root.IsLeaf() {
			t.Errorf("Expected new root to be a leaf node, but it's not")
		}
		expectedKeys := []int{20}
		if !reflect.DeepEqual(root.GetIterableKeys(), expectedKeys) {
			t.Errorf("Expected merged root to have keys %v, got %v", expectedKeys, root.GetIterableKeys())
		}
	})

	t.Run("Delete last element in tree", func(t *testing.T) {
		tree := CreateTree[int, string](2)
		tree.Insert(10, "ten")
		tree.Delete(10)
		root, err := tree.GetRootNode()
		if err == nil {
			t.Errorf("Expected error getting root from empty tree, but got nil")
		}
		if root != nil {
			t.Errorf("Expected root to be nil after deleting last element, but got %v", root)
		}
	})
}
