// Package orderedlist - A simple implementation of an ordered list
package orderedlist

import (
	"cmp"
	"errors"
	"fmt"
	"slices"
)

// This package should allow elements to be listed
// in order. This should be a generic ordered list
// that supports a minimum of ints and strings

// Space - O(n)
// Insert - O(n)

// OrderedList is a list that is ordered. When an element
// is added to the list, we should organize the elements so
// that they are in ascending order depending on the type
type OrderedList[T cmp.Ordered] struct {
	Data []T
}

// Insert adds an element to the ordered list. It uses binary search to find the
// correct position in O(log n) time. The insertion itself is O(n) due to
// shifting elements in the slice.
func (o *OrderedList[T]) Insert(element T) (int, bool) {
	// Find the position to insert the element to maintain order.
	// slices.BinarySearch is O(log n).
	idx, found := slices.BinarySearch(o.Data, element)

	// If the element already exists, do nothing to prevent duplicates.
	if found {
		return idx, false
	}

	// Insert the element at the correct index to maintain the sorted order.
	// slices.Insert is O(n).
	o.Data = slices.Insert(o.Data, idx, element)
	return idx, true
}

// Remove removes an element from the ordered list
func (o *OrderedList[T]) Remove(element T) error {
	idx, found := slices.BinarySearch(o.Data, element)
	if !found {
		return errors.New("the element doesn't exist in the list")
	}

	// delete the element using slices.Delete
	o.Data = slices.Delete(o.Data, idx, idx+1)

	return nil
}

// ToList return the data from the ordered list as a slice
func (o *OrderedList[T]) ToList() []T {
	return o.Data
}

// Get returns the element at the index indicated.
func (o *OrderedList[T]) Get(idx int) T {
	return o.Data[idx]
}

// GetData returns the all of the elements in the ordered list.
func (o *OrderedList[T]) GetData() []T {
	return o.Data
}

// Len returns the length of the ordered list
func (o *OrderedList[T]) Len() int {
	return len(o.Data)
}

// String returns a string representation of the ordered list
func (o *OrderedList[T]) String() string {
	return fmt.Sprintf("%v", o.Data)
}

// Compare compares a list with the OrderedList. If each
// element matches and is in the same position as in the
// argument, then return true. Otherwise, return false
func (o *OrderedList[T]) Compare(compare []T) bool {
	// set the default value of the comparison to true
	// this allows us to match cases such as when the lists
	// are empty

	// quick check, if the lengths don't match, return false
	if len(o.Data) != len(compare) {
		return false
	}
	// loop through each element of o.Data
	// if a value doesn't match, return false
	// since each element of both lists must
	// match to be the same. Their positions
	// must match as well.
	for i, v := range o.Data {
		comparedValue := compare[i]

		//
		if v != comparedValue {
			return false
		}
	}

	return true
}
