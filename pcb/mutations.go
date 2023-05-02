package pcb

import (
	"genetic_pcb/genetic"

	"github.com/mroth/weightedrand/v2"

	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
)

type mutation = func(i *Pcb, c *genetic.GeneticContext)
type mutationChooser = weightedrand.Chooser[mutation, int]

type MutationWeights struct {
	GlobalMutationWeight                  int
	TranslateComponentGroupMutationWeight int
	RegenerateNetMutationWeight           int
	RotateComponentMutationWeight         int
	RerouteEdgeMutationWeight             int
}

func (pgo *PcbGeneticOperators) buildMutationChooser() *mutationChooser {
	chooser, _ := weightedrand.NewChooser(
		weightedrand.NewChoice(pgo.globalMutation, pgo.mutationWeights.GlobalMutationWeight),
		weightedrand.NewChoice(pgo.netMutation, pgo.mutationWeights.RegenerateNetMutationWeight),
		weightedrand.NewChoice(pgo.translateComponentGroup, pgo.mutationWeights.TranslateComponentGroupMutationWeight),
		weightedrand.NewChoice(pgo.rotateComponent, pgo.mutationWeights.RotateComponentMutationWeight),
		weightedrand.NewChoice(pgo.rerouteEdge, pgo.mutationWeights.RerouteEdgeMutationWeight),
	)

	return chooser
}

func (pgo *PcbGeneticOperators) globalMutation(i *Pcb, c *genetic.GeneticContext) {
	for j := range i.Genome.Components {
		if c.RandomGenerator.Float64() < pgo.mutateSingleComponentProb {
			component := &i.Genome.Components[j]
			component.CX, component.CY = GetComponentRandomPositionInBoundaries(component, pgo.maxX, pgo.maxY, c.RandomGenerator)
			component.Rotation = c.RandomGenerator.Float64() * 360
			PlaceComponentNodes(i.Genome.Nodes, component)
		}
	}
}

func (pgo *PcbGeneticOperators) translateComponentGroup(i *Pcb, c *genetic.GeneticContext) {
	DX, DY := (c.RandomGenerator.Float64()-0.5)*pgo.maxX*0.1, (c.RandomGenerator.Float64()-0.5)*pgo.maxY*0.1

	for j := range i.Genome.Components {
		if c.RandomGenerator.Float64() < pgo.mutateSingleComponentProb {
			component := &i.Genome.Components[j]
			newX, newY := component.CX+DX, component.CY+DY

			if newX >= 0 && newX < pgo.maxX && newY >= 0 && newY < pgo.maxY {
				component.CX, component.CY = component.CX+DX, component.CY+DY
				PlaceComponentNodes(i.Genome.Nodes, component)
			}

		}

	}
}

func (pgo *PcbGeneticOperators) rotateComponent(i *Pcb, c *genetic.GeneticContext) {
	for j := range i.Genome.Components {
		if c.RandomGenerator.Float64() < pgo.mutateSingleComponentProb {
			component := &i.Genome.Components[j]
			component.Rotation = c.RandomGenerator.Float64() * 360
			PlaceComponentNodes(i.Genome.Nodes, component)
		}
	}
}

func clip(x, min, max float64) float64 {
	if x < min {
		return min
	}

	if x > max {
		return max
	}

	return x
}

func (pgo *PcbGeneticOperators) localMutation(i *Pcb, c *genetic.GeneticContext) {
	node := c.RandomGenerator.Intn(len(i.Genome.Nodes))

	dx, dy := (c.RandomGenerator.Float64()*pgo.localMutationMaxDelta)-pgo.localMutationMaxDelta/2, (c.RandomGenerator.Float64()*pgo.localMutationMaxDelta)-pgo.localMutationMaxDelta/2
	x, y := i.Genome.Nodes[node].X, i.Genome.Nodes[node].Y

	// fmt.Printf("%v, %v, %v\n", node, dx, dy)

	i.Genome.Nodes[node] = Node{X: clip(x+dx, 0, pgo.maxX), Y: clip(y+dy, 0, pgo.maxY)}

}

func (pgo *PcbGeneticOperators) netMutation(i *Pcb, c *genetic.GeneticContext) {
	netI := c.RandomGenerator.Intn(len(i.Genome.Nets))
	GenerateNet(i, netI, c.RandomGenerator)
}

func (pgo *PcbGeneticOperators) rerouteEdge(i *Pcb, c *genetic.GeneticContext) {
	edgeIndex := c.RandomGenerator.Intn(len(i.Genome.Edges))
	edge := i.Genome.Edges[edgeIndex]
	netIndex := edge.Net
	net := &i.Genome.Nets[netIndex]

	g := simple.NewUndirectedGraph()

	idToNodes := make(map[int64]int, len(i.Genome.Nodes))
	nodesToId := make(map[int]int64, len(i.Genome.Nodes))

	for _, node := range net.Nodes {
		n := g.NewNode()
		idToNodes[n.ID()] = node
		nodesToId[node] = n.ID()
		g.AddNode(n)
	}

	for i, edge := range i.Genome.Edges {
		if i != edgeIndex && edge.Net == netIndex {
			g.SetEdge(g.NewEdge(g.Node(nodesToId[edge.From]), g.Node(nodesToId[edge.To])))
		}
	}

	ccs := topo.ConnectedComponents(g)

	if len(ccs) > 2 {
		panic("Removing an edge resulted in more than 2 connected components, which is impossible for proper nets")
	}

	// fmt.Printf("Removing edge %d, %d from net %d\n", edge.From, edge.To, netIndex)
	// // Stampa delle componenti fortemente connesse trovate
	// for i, scc := range ccs {
	// 	if i > 1 {

	// 	}
	// 	fmt.Printf("Componente fortemente connessa %d: ", i)
	// 	for _, n := range scc {
	// 		fmt.Printf("%v ", idToNodes[n.ID()])
	// 	}
	// 	fmt.Println()
	// }

	n1 := idToNodes[ccs[0][c.RandomGenerator.Intn(len(ccs[0]))].ID()]
	n2 := idToNodes[ccs[1][c.RandomGenerator.Intn(len(ccs[1]))].ID()]

	// fmt.Printf("Chosen nodes %d %d\n", n1, n2)

	i.Genome.Edges = append(
		i.Genome.Edges[:edgeIndex],
		i.Genome.Edges[edgeIndex+1:]...,
	)

	i.Genome.Edges = append(i.Genome.Edges, Edge{From: n1, To: n2, Net: netIndex})

	SortPcbEdges(i)
}
