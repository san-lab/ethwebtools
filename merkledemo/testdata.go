package merkledemo

var AliceToDerek = []NodeData{{"Alice", 100}, {"Bob", 200}, {"Cecil", 300}, {"Derek", 400}}
var AliceToHarry = append(AliceToDerek, NodeData{"Eric", 250}, NodeData{"Felix", 370},
	NodeData{"George", 525}, NodeData{"Harry", 500})

func NodesToData(nodes []*Node) []NodeData {
	nds := make([]NodeData, len(nodes))
	for i, nd := range nodes {
		nds[i].NodeBalance = nd.Data.NodeBalance
		nds[i].NodeID = nd.Data.NodeID
	}
	return nds
}
