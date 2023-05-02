package pcb

import (
	"genetic_pcb/genetic"

	"github.com/mroth/weightedrand/v2"
)

type mutation = func(i *Pcb, c *genetic.GeneticContext)
type mutationChooser = weightedrand.Chooser[mutation, int]

type MutationWeights struct {
	GlobalMutationWeight                  int
	TranslateComponentGroupMutationWeight int
	RegenerateNetMutationWeight           int
	RotateComponentMutationWeight         int
}

func (pgo *PcbGeneticOperators) buildMutationChooser() *mutationChooser {
	chooser, _ := weightedrand.NewChooser(
		weightedrand.NewChoice(pgo.globalMutation, pgo.mutationWeights.GlobalMutationWeight),
		weightedrand.NewChoice(pgo.netMutation, pgo.mutationWeights.RegenerateNetMutationWeight),
		weightedrand.NewChoice(pgo.translateComponentGroup, pgo.mutationWeights.TranslateComponentGroupMutationWeight),
		weightedrand.NewChoice(pgo.rotateComponent, pgo.mutationWeights.RotateComponentMutationWeight),
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
