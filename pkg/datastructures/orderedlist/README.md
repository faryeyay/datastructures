# `orderedlist` Package

The `orderedlist` package provides a simple, generic implementation of an ordered list in Go. This list automatically maintains its elements in ascending order as they are added, supporting any type that implements the `cmp.Ordered` interface (e.g., `int`, `string`, `float64`).

## Features

  * **Generic:** Works with any ordered type.
  * **Automatic Ordering:** Elements are kept in ascending order upon insertion.
  * **Basic Operations:** Includes methods for inserting, retrieving, and comparing lists.

## Usage

Here's how to create and use an `OrderedList`:

```go
package main

import (
	"fmt"
	"orderedlist" // Assuming orderedlist is in your GOPATH or module
)

func main() {
	// Create an ordered list of integers
	intList := orderedlist.OrderedList[int]{}
	intList.Insert(5)
	intList.Insert(2)
	intList.Insert(8)
	intList.Insert(1)

	fmt.Println("Integer List:", intList.String()) // Output: Integer List: [1 2 5 8]

	// Create an ordered list of strings
	stringList := orderedlist.OrderedList[string]{}
	stringList.Insert("banana")
	stringList.Insert("apple")
	stringList.Insert("cherry")

	fmt.Println("String List:", stringList.String()) // Output: String List: [apple banana cherry]

	// Get an element
	fmt.Println("Element at index 1 (intList):", intList.Get(1)) // Output: Element at index 1 (intList): 2

	// Get the length
	fmt.Println("Length of intList:", intList.Len()) // Output: Length of intList: 4

	// Convert to a slice
	slice := intList.ToList()
	fmt.Println("intList as slice:", slice) // Output: intList as slice: [1 2 5 8]

	// Compare lists
	listToCompare := []int{1, 2, 5, 8}
	fmt.Println("intList matches [1 2 5 8]:", intList.Compare(listToCompare)) // Output: intList matches [1 2 5 8]: true

	listToCompareFalse := []int{1, 2, 8, 5}
	fmt.Println("intList matches [1 2 8 5]:", intList.Compare(listToCompareFalse)) // Output: intList matches [1 2 8 5]: false
}

```

## API

### `type OrderedList[T cmp.Ordered]`

Represents an ordered list. `T` must be a type that implements `cmp.Ordered`.

#### `func (o *OrderedList[T]) Insert(element T)`

Inserts an `element` into the list while maintaining its sorted order.

  * **Time Complexity:** O(n), where n is the number of elements in the list, due to potential shifting of elements.
  * **Space Complexity:** O(1) amortized for the insertion itself, but O(n) for the overall list storage.

#### `func (o *OrderedList[T]) ToList() []T`

Returns the elements of the ordered list as a Go slice.

#### `func (o *OrderedList[T]) Get(idx int) T`

Returns the element at the specified `idx`.

#### `func (o *OrderedList[T]) Len() int`

Returns the number of elements in the ordered list.

#### `func (o *OrderedList[T]) String() string`

Returns a string representation of the ordered list.

#### `func (o *OrderedList[T]) Compare(compare []T) bool`

Compares the `OrderedList` with a provided slice `compare`. It returns `true` if both the length and all elements at corresponding positions are identical, otherwise `false`.

## Time and Space Complexity

  * **Space Complexity (Overall):** O(n), where n is the number of elements stored in the list.
  * **Insert Operation:**
      * **Time Complexity:** O(n) in the worst case, as existing elements might need to be shifted to accommodate the new element.
      * **Space Complexity:** O(1) amortized for the insertion itself, not counting the overall list storage.

-----