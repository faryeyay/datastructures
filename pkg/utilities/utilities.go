// Package utilities for commonly used methods that need a place
// Phone home E.T. Phone hom
package utilities

// CloneSlice returns a copy of a slice so that I don't have to worry
// about why I'm manipulating something else. 
func CloneSlice[T any](inputSlice []T) []T {
	// 1. Get the lenght of the input slice. We need this so that
	// we can create another slice with the same length
	sliceLength := len(inputSlice)

	// 2. Make any empty slice of type T with the same length. Remember
	// sliceLength from above, we're going to use that here to create
	// the slice of the correct length 
	clonedElements := make([]T, sliceLength)

	// 3. Use the 'copy' built-in function to copy elements from the
	//    original slice's inputSlice into the new slice.
	copy(clonedElements, inputSlice)

	return clonedElements
}