package merkledemo

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strconv"

	"github.com/san-lab/commongo/gohttpservice"
	"github.com/san-lab/ethwebtool/templates"
)

var Forest = map[string]*MerkleData{}

func CallMerkleDemo(r *http.Request, rdata *templates.RenderData) {
	sescook, _ := r.Cookie(gohttpservice.SessionIdName)
	fmt.Println(sescook.Value)
	rdata.TemplateName = "merkledemo"
	brcs := r.FormValue("branchcount")
	brc, _ := strconv.Atoi(brcs)

	strategys := r.FormValue("strategy") //'s is for string, not for plural
	if len(strategys) == 0 {
		strategys = "Sum"
	}
	var leafdata []NodeData
	mdat, ok := Forest[sescook.Value]
	if !ok {

		mdat = &MerkleData{}
		mdat.Id = sescook.Value
		mdat.PString = "Merkle Demo"

		mdat.Branchcount = brc
		mdat.Strategy = Sum
		Forest[sescook.Value] = mdat
		leafdata = AliceToHarry
	} else {
		leafdata = NodesToData(mdat.Tree.Leaves)
	}

	mdat.Strategy = MatchStrategy(strategys)
	//see if leaf edit is requested
	action := r.FormValue("action")
	if action == "Change" {
		visid := r.FormValue("leafid")
		if len(visid) > 0 {
			userid := r.FormValue("newUserId")
			nbalance := r.FormValue("newBalance")
			balance, err := strconv.Atoi(nbalance)
			if err != nil {
				fmt.Println(err)
			}
			idx := mdat.Tree.VisIDToIdx(visid)
			leafdata[idx].NodeID = userid
			leafdata[idx].NodeBalance = balance
		}

	}
	mdat.Tree = NewTree(sha256.New(), leafdata, mdat.Branchcount, mdat.Strategy)
	mdat.DisplayTree = mdat.Tree
	//fmt.Println(mdat)
	rdata.BodyData = mdat
}

type MerkleData struct {
	Id          string
	PString     string
	Tree        *Tree
	Proofs      map[int]*Tree
	DisplayTree *Tree
	Branchcount int
	Strategy    Strategy
}
