// utilities_test.go
package utilities

import (
	"reflect"
	"testing"
)

func TestCloneSlice(t *testing.T) {
	// Test case 1: Empty slice of integers
	t.Run("EmptyIntSlice", func(t *testing.T) {
		input := []int{}
		cloned := CloneSlice(input)

		if cloned == nil {
			t.Errorf("Cloned slice is nil, expected empty slice")
		}
		if len(cloned) != 0 {
			t.Errorf("Cloned slice length is %d, expected 0", len(cloned))
		}
		// Ensure it's a new slice, not just the same underlying array (though for empty, it's often the same address)
		if !reflect.DeepEqual(input, cloned) {
			t.Errorf("Cloned slice %v is not deep equal to input %v", cloned, input)
		}
	})

	// Test case 2: Nil slice of strings
	t.Run("NilStringSlice", func(t *testing.T) {
		var input []string = nil
		cloned := CloneSlice(input)

		if cloned == nil {
			t.Errorf("Cloned slice is nil, expected empty slice (make([]T, 0) returns non-nil empty slice for nil input)")
		}
		if len(cloned) != 0 {
			t.Errorf("Cloned slice length is %d, expected 0", len(cloned))
		}
		// For nil input, make([]T, 0) is returned, which is not strictly DeepEqual to nil,
		// but it's the correct and safe behavior for `make` with len 0.
	})

	// Test case 3: Slice of integers
	t.Run("IntSlice", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		cloned := CloneSlice(input)

		if cloned == nil {
			t.Fatalf("Cloned slice is nil, expected non-nil")
		}
		if len(cloned) != len(input) {
			t.Errorf("Cloned slice length is %d, expected %d", len(cloned), len(input))
		}
		if cap(cloned) != len(input) { // make(len) sets cap to len
			t.Errorf("Cloned slice capacity is %d, expected %d", cap(cloned), len(input))
		}
		if !reflect.DeepEqual(input, cloned) {
			t.Errorf("Cloned slice %v is not deep equal to input %v", cloned, input)
		}
		// Verify independence: modify cloned slice and check original
		cloned[0] = 99
		if input[0] == 99 {
			t.Errorf("Original slice was modified, expected independence. Original: %v", input)
		}
	})

	// Test case 4: Slice of strings
	t.Run("StringSlice", func(t *testing.T) {
		input := []string{"apple", "banana", "cherry"}
		cloned := CloneSlice(input)

		if cloned == nil {
			t.Fatalf("Cloned slice is nil, expected non-nil")
		}
		if len(cloned) != len(input) {
			t.Errorf("Cloned slice length is %d, expected %d", len(cloned), len(input))
		}
		if !reflect.DeepEqual(input, cloned) {
			t.Errorf("Cloned slice %v is not deep equal to input %v", cloned, input)
		}
		// Verify independence
		cloned[1] = "grape"
		if input[1] == "grape" {
			t.Errorf("Original slice was modified, expected independence. Original: %v", input)
		}
	})

	// Test case 5: Slice of custom structs (value type)
	t.Run("StructSliceValue", func(t *testing.T) {
		type MyStruct struct {
			ID   int
			Name string
		}
		input := []MyStruct{{ID: 1, Name: "A"}, {ID: 2, Name: "B"}}
		cloned := CloneSlice(input)

		if cloned == nil {
			t.Fatalf("Cloned slice is nil, expected non-nil")
		}
		if len(cloned) != len(input) {
			t.Errorf("Cloned slice length is %d, expected %d", len(cloned), len(input))
		}
		if !reflect.DeepEqual(input, cloned) {
			t.Errorf("Cloned slice %v is not deep equal to input %v", cloned, input)
		}
		// Verify independence (modifying a field in the cloned struct)
		cloned[0].Name = "X"
		if input[0].Name == "X" {
			t.Errorf("Original struct in slice was modified, expected independence. Original: %v", input)
		}
	})

	// Test case 6: Slice of pointers to custom structs (shallow copy behavior)
	t.Run("StructPointerSlice", func(t *testing.T) {
		type MyStruct struct {
			ID   int
			Name string
		}
		s1 := &MyStruct{ID: 1, Name: "A"}
		s2 := &MyStruct{ID: 2, Name: "B"}
		input := []*MyStruct{s1, s2}
		cloned := CloneSlice(input)

		if cloned == nil {
			t.Fatalf("Cloned slice is nil, expected non-nil")
		}
		if len(cloned) != len(input) {
			t.Errorf("Cloned slice length is %d, expected %d", len(cloned), len(input))
		}
		if !reflect.DeepEqual(input, cloned) {
			t.Errorf("Cloned slice %v is not deep equal to input %v", cloned, input)
		}

		// Important: Verify shallow copy. The slice itself is new, but the pointers within it are the same.
		// Modifying the *content* of a struct pointed to by an element in the cloned slice
		// should affect the original slice's element as well.
		cloned[0].Name = "X"
		if input[0].Name != "X" {
			t.Errorf("Original struct content was NOT modified, expected shallow copy behavior. Original: %v", input)
		}
		// However, replacing a pointer in the cloned slice should NOT affect the original.
		cloned[1] = &MyStruct{ID: 3, Name: "C"}
		if input[1] == cloned[1] {
			t.Errorf("Replacing pointer in cloned slice affected original, expected independence of slice elements. Original: %v", input)
		}
		if input[1].Name != "B" { // Ensure original's second element is still the original s2
			t.Errorf("Original slice's second element changed unexpectedly: %v", input[1])
		}
	})
}
