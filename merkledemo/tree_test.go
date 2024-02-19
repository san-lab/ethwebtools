package merkledemo

import (
	"crypto/sha256"
	"fmt"
	"testing"
)

func TestTtee(t *testing.T) {
	leaves := append(AliceToHarry, AliceToHarry...)
	tree := NewTree(sha256.New(), leaves, 2, Sum)
	fmt.Println(tree)
	for _, leaf := range tree.Leaves {
		fmt.Println(leaf.VisId())
	}

}

func TestVisID(t *testing.T) {
	leaves := append(AliceToHarry, AliceToHarry...)
	tree := NewTree(sha256.New(), leaves, 5, Sum)
	for n := range tree.Leaves {
		visid := tree.Leaves[n].VisId()
		fmt.Println(n, visid)
		n2 := tree.VisIDToIdx(visid)
		fmt.Println((n == n2))
	}

}
