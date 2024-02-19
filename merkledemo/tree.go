package merkledemo

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"hash"
	"strconv"
	"strings"
)

type NodeData struct {
	NodeID      string
	NodeBalance int
}

type Node struct {
	isLeaf   bool
	Data     *NodeData
	Children []*Node
	parent   *Node

	Hash    []byte
	Idx     int //position as sibilng, here 0 .. branchCount-1
	Lvl     int
	IsRoot  bool
	IsWrong bool
}

type Tree struct {
	Hash            hash.Hash
	Root            *Node
	BranchCount     int
	Leaves          []*Node
	BalanceStrategy Strategy
}

type Strategy string

const Sum Strategy = "Sum"
const Max Strategy = "Max"
const Min Strategy = "Min"
const Both Strategy = "Both"
const None Strategy = "None"

func MatchStrategy(label string) Strategy {
	label = strings.ToLower(label)
	switch label {
	case "sum":
		return Sum
	case "max":
		return Max
	case "min":
		return Min
	case "both":
		return Both
	default:
		return None
	}
}

func nodeWithChildrenString(nd *Node, i int) string {
	if nd == nil {
		return ""
	}
	s := "" //fmt.Sprintf("i:%v::%s\n", i, nd)
	for j, c := range nd.Children {

		s += fmt.Sprintf("i:%v::%s\n", j, c)
	}
	for j, c := range nd.Children {
		s += nodeWithChildrenString(c, j)
	}
	return s

}

func (nd *Node) String() string {
	s := fmt.Sprintf("level: %v index: %v isLeaf:%v children: %v\n", nd.Lvl, nd.Idx, nd.isLeaf, len(nd.Children))

	s += fmt.Sprintf(" %v %v\n", nd.Data, hex.EncodeToString(nd.Hash))
	return s
}

func (nd *Node) FillFromChildern(hash hash.Hash, strategy Strategy) {
	if nd.isLeaf {
		return
	}
	nd.Data = &NodeData{}
	nd.Data.NodeBalance = 0
	hash.Reset()
	for _, c := range nd.Children {
		if c != nil {
			nd.Data.NodeBalance = udateUpperBalance(nd.Data.NodeBalance, c.Data.NodeBalance, strategy)
			if strategy == Both {
				hash.Write([]byte(strconv.Itoa(c.Data.NodeBalance)))

			}
			hash.Write(c.Hash)

		}
	}
	if strategy != Both && strategy != None {
		hash.Write([]byte(strconv.Itoa(nd.Data.NodeBalance)))
	}
	nd.Hash = hash.Sum(nil)

}

// Test if two nodes are equal
func (nd *Node) Eq(n *Node) bool {
	if nd == nil && n == nil {
		return true
	}
	if nd == nil || n == nil {
		return false
	}

	if nd.Data.NodeID != n.Data.NodeID || nd.Data.NodeBalance != n.Data.NodeBalance {
		return false
	}
	if !bytes.Equal(nd.Hash, n.Hash) {
		return false
	}

	return true

}

func (tr *Tree) String() string {
	s := tr.Root.String()
	s += nodeWithChildrenString(tr.Root, 0)

	return s
}

// https://graphics.stanford.edu/~seander/bithacks.html#RoundUpPowerOf2
func nextPowerOf2(i int) int { //int32 actually
	if i == 0 {
		return 1
	}
	i--
	i |= i >> 1
	i |= i >> 2
	i |= i >> 4
	i |= i >> 8
	i |= i >> 16
	i++
	return i

}

func NewTree(hf hash.Hash, data []NodeData, branchCount int, strategy Strategy) *Tree {
	if branchCount < 2 {
		branchCount = 2
	}

	leaves := make([]*Node, 0)

	tr := new(Tree)
	tr.BalanceStrategy = strategy
	tr.BranchCount = branchCount //nextPowerOf2(branchCount)
	tr.Hash = hf
	for i := 0; i < len(data); i++ {
		nod := new(Node)
		nod.isLeaf = true
		nod.Data = &NodeData{}
		nod.Data.NodeBalance = data[i].NodeBalance
		nod.Data.NodeID = data[i].NodeID
		tr.Hash.Reset()
		tr.Hash.Write([]byte(data[i].NodeID))
		tr.Hash.Write([]byte(strconv.Itoa(data[i].NodeBalance)))

		nod.Hash = tr.Hash.Sum(nil)
		nod.Lvl = 0

		leaves = append(leaves, nod)

	}
	tr.Leaves = leaves

	tr.makeParents(leaves, 1)
	return tr
}

func (tr *Tree) makeParents(level []*Node, lvl int) {
	p := len(level)
	if p == 1 {
		tr.Root = level[0]
		tr.Root.IsRoot = true
		return
	}

	nextlvl := make([]*Node, 0)
	for i := 0; i < p; {
		parent := new(Node)
		parent.Lvl = lvl
		parent.Children = make([]*Node, tr.BranchCount)
		tr.Hash.Reset()
		for j := 0; i < p && j < tr.BranchCount; j++ {
			parent.Children[j] = level[i]
			level[i].Idx = j
			level[i].parent = parent
			i++

		}
		parent.Idx = (i - 1) / tr.BranchCount
		parent.Lvl = lvl
		parent.FillFromChildern(tr.Hash, tr.BalanceStrategy)
		nextlvl = append(nextlvl, parent)
	}

	tr.makeParents(nextlvl, lvl+1)

}

func udateUpperBalance(upperBalance, leafBalance int, strategy Strategy) int {
	switch strategy {
	case Sum, Both:
		return upperBalance + leafBalance
	case Max:
		if leafBalance > upperBalance {
			return leafBalance
		}
		return upperBalance
	case Min:
		if leafBalance < upperBalance {
			return leafBalance
		}
		return upperBalance
	case None:
		return 0
	default:
		fmt.Println("Unknown strategy", strategy)
		return upperBalance
	}

}

func (tr *Tree) cloneProofNode(org *Node) *Node {
	cp := new(Node)
	if org.Data != nil {
		cp.Data = new(NodeData)
		cp.Data.NodeBalance = org.Data.NodeBalance
	}
	cp.Hash = org.Hash
	cp.Lvl = org.Lvl
	cp.Idx = org.Idx
	return cp
}

// TODO rework this
func (tr *Tree) GetProof(leafID string) (*Tree, error) {
	leafIdx := -1
	for i, l := range tr.Leaves {
		if l.Data.NodeID == leafID {
			leafIdx = i
			break
		}
	}
	if leafIdx == -1 {
		return nil, fmt.Errorf("Leaf not found")
	}

	proof := new(Tree)
	provednode := tr.Leaves[leafIdx]

	shadownode := tr.cloneProofNode(provednode)

	shadownode.isLeaf = true

	root := tr.TracePathToRoot(provednode, shadownode)

	proof.Hash = tr.Hash
	proof.Root = root

	return proof, nil
}

func (tr *Tree) TracePathToRoot(org *Node, cp *Node) *Node {

	if org.parent == nil {
		cp.IsRoot = true

		return cp
	}
	cppar := tr.cloneProofNode(org.parent)

	for i, c := range org.parent.Children {
		if c == nil {
			continue
		}

		if i == cp.Idx {

			cppar.Children = append(cppar.Children, cp)
			cp.parent = cppar
		} else {
			cc := tr.cloneProofNode(c)
			cppar.Children = append(cppar.Children, cc)
			cc.isLeaf = true
			cc.parent = cppar

		}

	}

	return tr.TracePathToRoot(org.parent, cppar)

}

func VerifyProofConsistency(proof *Tree) bool {
	return IsSubtreeConsistent(proof.Root, proof.Hash)

}

func IsSubtreeConsistent(nd *Node, hash hash.Hash) bool {
	if nd.isLeaf {

		hash.Reset()
		hash.Write([]byte(nd.Data.NodeID))
		hash.Write([]byte(strconv.Itoa(nd.Data.NodeBalance)))
		ht := hash.Sum(nil)
		loc := bytes.Equal(ht, nd.Hash)
		if !loc {
			fmt.Println("Wrong leaf!")
			return loc

		}
		return true
	}
	hash.Reset()
	for _, c := range nd.Children {
		if c != nil {
			hash.Write(c.Hash)
		}
	}
	if !bytes.Equal(hash.Sum(nil), nd.Hash) {
		fmt.Println("Wrong node", nd)
		return false
	}
	for _, c := range nd.Children {
		if !IsSubtreeConsistent(c, hash) {
			return false
		}
	}
	return true
}

func (tr *Tree) NiceRoot() string {
	return hex.EncodeToString(tr.Root.Hash)
}
