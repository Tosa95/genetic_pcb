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

func GeneratePcbFull(componentTemplates []Component, componentN, netN int, maxX, maxY float64, randomGenerator *rand.Rand) *Pcb {
	nodeI := 0

	pcb := &Pcb{
		Genome: &Genome{
			Nodes:      make([]Node, 0),
			Components: make([]Component, componentN),
			Nets:       make([]Net, netN),
		},
	}

	// Generate components and nodes
	for i := 0; i < componentN; i++ {
		c := componentTemplates[randomGenerator.Intn(len(componentTemplates))].copy()

		for j := range c.Nodes {
			cn := &c.Nodes[j]
			n := Node{X: cn.DX, Y: cn.DY, Component: i}
			cn.Node = nodeI
			pcb.Genome.Nodes = append(pcb.Genome.Nodes, n)
			nodeI++

		}

		c.CX, c.CY = GetComponentRandomPositionInBoundaries(c, maxX, maxY, randomGenerator)
		PlaceComponentNodes(pcb.Genome.Nodes, c)

		pcb.Genome.Components[i] = *c
	}

	// Randomly assign each node to a net

	for i := range pcb.Genome.Nodes {
		net := randomGenerator.Intn(netN)
		pcb.Genome.Nets[net].Nodes = append(pcb.Genome.Nets[net].Nodes, i)
	}

	for i := range pcb.Genome.Nets {
		GenerateNet(pcb, i, randomGenerator)
	}

	return pcb
}

func ScrumblePcb(original *Pcb, maxX, maxY float64) *Pcb {
	res := original.Genome.copy()

	copy(res.Edges, original.Genome.Edges)

	s1 := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(s1)

	for i := 0; i < len(res.Components); i++ {
		c := &res.Components[i]
		c.CX, c.CY = GetComponentRandomPositionInBoundaries(c, maxX, maxY, randomGenerator)
		PlaceComponentNodes(res.Nodes, c)
	}

	return NewPcb(res)
}
