// Package testbench - Test Bench for Data Structures
package testbench

import (
	"log"

	"github.com/faryeyay/insights/primitives/pkg/datastructures/bplustree"
)

func TestBPlus() error {
	rootKey := "root0"
	order := 5
	log.Print("Creating a B+ Tree")
	log.Printf("Setting the root of the tree to %s", rootKey)
	bPlusTree := bplustree.CreateTreeWithRootKey[string, string](order, rootKey)

	_, _, err := bPlusTree.Search("terminator")
	return err
}
