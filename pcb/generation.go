package pcb

import (
	"math/rand"
	"time"
)

func generateRandomNode(maxX, maxY float64, randomGenerator *rand.Rand) Node {
	return Node{X: randomGenerator.Float64() * maxX, Y: randomGenerator.Float64() * maxY}
}

func GeneratePcb(nNodes, nEdges int, maxX, maxY float64) *Pcb {
	nodes := make([]Node, nNodes)
	edges := make([]Edge, nEdges)

	s1 := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(s1)

	for i := 0; i < nNodes; i++ {
		nodes[i] = generateRandomNode(maxX, maxY, randomGenerator)
	}

	for i := 0; i < nEdges; i++ {
		startNode := randomGenerator.Intn(nNodes)
		endNode := randomGenerator.Intn(nNodes)

		for ; endNode == startNode; endNode = randomGenerator.Intn(nNodes) {
		}

		edges[i] = Edge{From: startNode, To: endNode}
	}

	g := Genome{
		Nodes: nodes,
		Edges: edges,
	}

	return NewPcb(&g)
}

func GeneratePcbWithNets(netSz, netN int, maxX, maxY float64) *Pcb {
	s1 := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(s1)

	nodeN := netSz * netN
	pcb := &Pcb{
		Genome: &Genome{
			Nodes: make([]Node, 0, nodeN),
			Edges: make([]Edge, 0, (netSz-1)*netN),
			Nets:  make([]Net, netN),
		},
	}

	nodesI := 0

	for netI := 0; netI < netN; netI++ {
		net := &pcb.Genome.Nets[netI]

		for j := 0; j < netSz; j++ {
			node := generateRandomNode(maxX, maxY, randomGenerator)
			pcb.Genome.Nodes = append(pcb.Genome.Nodes, node)
			net.Nodes = append(net.Nodes, nodesI)
			nodesI++
		}

		GenerateNet(pcb, netI, randomGenerator)
	}

	return pcb
}

func ScrumblePcb(original *Pcb, maxX, maxY float64) *Pcb {
	res := original.Genome.copy()

	copy(res.Edges, original.Genome.Edges)

	s1 := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(s1)

	for i := 0; i < len(original.Genome.Nodes); i++ {
		res.Nodes[i] = generateRandomNode(maxX, maxY, randomGenerator)
	}

	return NewPcb(res)
}
