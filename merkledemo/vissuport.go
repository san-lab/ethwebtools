package merkledemo

import (
	"encoding/hex"
	"fmt"
	"html/template"
	"regexp"
	"strconv"
)

type Net struct {
	Nodes template.JS
	Edges template.JS
}

func (tr *Tree) VisNet() Net {
	return Net{tr.VisNodes(), tr.VisEdges()}
}

func (tr *Tree) visAddWithChildren(p *Node, s *string) {
	if p != nil {
		*s += "{"
		color := "99ffff"
		if p.isLeaf {
			color = "33ff99"
			*s += "isLeaf: true,"
		} else if p.Data != nil {
			color = "fee9b4"
		}
		if p.IsRoot {
			color = "ffA050"
		}

		label := ""
		if p.Data != nil {
			if len(p.Data.NodeID) > 0 {
				label += fmt.Sprintf("%s: ", p.Data.NodeID)
			}
			label += fmt.Sprintf("<b>$%v</b>", p.Data.NodeBalance) + "\\n"
		}
		label += hex.EncodeToString(p.Hash)[:5] + "..." + hex.EncodeToString(p.Hash)[59:]
		*s += fmt.Sprintf("id:\"%s\", label: \"%s\", margin: { top: 15, right: 15, bottom: 15, left: 15 }, color: \"#%s\", heightConstraint:15,",
			p.VisId(), label, color)
		*s += " hash: \"" + hex.EncodeToString(p.Hash) + "\","
		if p.Data != nil {
			*s += " userId: \"" + p.Data.NodeID + "\"," + " balance: " + strconv.Itoa(p.Data.NodeBalance) + ","
		}

		*s += "title:'Node hash: " + hex.EncodeToString(p.Hash) +
			"<br/>Node balance: " + strconv.Itoa(p.Data.NodeBalance) +
			"<br/>hashing method: " + tr.StrategyString() + "',"

		*s += "},\n"
		for _, c := range p.Children {
			tr.visAddWithChildren(c, s)
		}
	}
}

func (tr *Tree) VisEdges() template.JS {
	s := "["
	arrowsFromCildren(tr.Root, &s)
	s += "]"
	return template.JS(s)
}

func arrowsFromCildren(p *Node, s *string) {
	for _, c := range p.Children {
		if c == nil {
			continue
		}
		*s += fmt.Sprintf("{from: \"%s\", to: \"%s\"},", c.VisId(), p.VisId())
		arrowsFromCildren(c, s)
	}
}

func (n *Node) VisId() string {
	s := fmt.Sprintf("L%vC%v", n.Lvl, n.Idx)
	if n.parent != nil {
		s += n.parent.VisId()
	}
	return s
}

var visidmatch = regexp.MustCompile(`C[0-9]+`)

// find node by visid
func (tr *Tree) VisIDToIdx(visid string) int {
	digits := visidmatch.FindAllString(visid, -1)
	base := 1
	idx := 0
	for _, pwrs := range digits {
		digit, _ := strconv.Atoi(pwrs[1:])
		idx += base * digit
		base *= tr.BranchCount

	}
	return idx

}

func (tr *Tree) VisNodes() template.JS {
	s := "["
	tr.visAddWithChildren(tr.Root, &s)

	s += "]"
	return template.JS(s)
}

func (tr *Tree) StrategyString() string {
	switch tr.BalanceStrategy {
	case Sum:
		return "Hash(L.Hash+R.Hash+(L.Balance+R.Balance))"
	case Max:
		return "Hash(L.Hash+R.Hash+Max(L.Balance,R.Balance))"
	case Min:
		return "Min"
	case Both:
		return "Hash(L.Balance+L.Hash+R.Balance+R.Hash)"
	case None:
		return "Hash(L.Hash+R.Hash)"
	default:
		return "Unknown strategy"
	}
}
