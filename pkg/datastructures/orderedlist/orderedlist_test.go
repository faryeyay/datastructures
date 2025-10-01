// Package orderedlist - A simple implementation of an ordered list
package orderedlist

import (
	"testing"
)

// TestOrderedListInsertAndOrderString tests that the elements in the list are
// sorted in lexicographical order.
func TestOrderedListInsertAndOrderString(t *testing.T) {
	list := OrderedList[string]{}

	list.Insert("abc")
	list.Insert("ghi")
	list.Insert("jkl")
	list.Insert("def")

	compareList := []string{"abc", "def", "ghi", "jkl"}
	compareListLength := len(compareList)

	// check the count of the lists
	if len(list.ToList()) != compareListLength {
		t.Errorf("the length of list does not match the expected length of %d", compareListLength)
	}
	// ensure sure the two lists match
	if !list.Compare(compareList) {
		t.Errorf("expected the order of the list to be: %s, got %s", compareList, list.Data)
	}
}

// TestOrderedListInsertAndOrderInt tests that the elements in the list are
// sorted in order.
func TestOrderedListInsertAndOrderInt(t *testing.T) {
	list := OrderedList[int]{}

	list.Insert(1)
	list.Insert(0)
	list.Insert(7)
	list.Insert(9)
	list.Insert(5)

	compareList := []int{0, 1, 5, 7, 9}
	compareListLength := len(compareList)

	// check the count of the lists
	if len(list.ToList()) != compareListLength {
		t.Errorf("the length of list does not match the expected length of %d", compareListLength)
	}
	// ensure sure the two lists match
	if !list.Compare(compareList) {
		t.Errorf("expected the order of the list to be: %d, got %d", compareList, list.Data)
	}
}
