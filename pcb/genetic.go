package pcb

import (
	"genetic_pcb/genetic"
	"math"

	"github.com/mroth/weightedrand/v2"
)

type PcbGeneticOperators struct {
	fitnessExp            float64
	mutateProb            float64
	mutateSinglePointProb float64
	minDist               float64
	maxX                  float64
	maxY                  float64
	nodeSz                float64
	edgeSz                float64
	localMutationMaxDelta float64
	localMutationProb     float64
	netMutationProb       float64
	edgeLengthPenalty     float64
}

func NewPcbGeneticOperators(
	fitnessExp float64,
	mutateProb float64,
	mutateSinglePointProb float64,
	minDist float64,
	maxX float64,
	maxY float64,
	nodeSz float64,
	edgeSz float64,
	localMutationMaxDelta float64,
	localMutationProb float64,
	netMutationProb float64,
	edgeLengthPenalty float64,
) *PcbGeneticOperators {

	pgo := PcbGeneticOperators{
		fitnessExp:            fitnessExp,
		mutateProb:            mutateProb,
		mutateSinglePointProb: mutateSinglePointProb,
		minDist:               minDist,
		maxX:                  maxX,
		maxY:                  maxY,
		nodeSz:                nodeSz,
		edgeSz:                edgeSz,
		localMutationMaxDelta: localMutationMaxDelta,
		localMutationProb:     localMutationProb,
		netMutationProb:       netMutationProb,
		edgeLengthPenalty:     edgeLengthPenalty,
	}

	return &pgo
}

func (pgo *PcbGeneticOperators) Evaluate(i *Pcb, c *genetic.GeneticContext) float64 {
	fitness := EvaluatePcb(i, pgo.minDist)
	fitness -= (GetTotalPcbLength(i) / (pgo.maxX + pgo.maxY)) * pgo.edgeLengthPenalty
	fitness = math.Pow(fitness, pgo.fitnessExp)
	return fitness
}

func (pgo *PcbGeneticOperators) copyComponentNodesToChild(c *Genome, p *Genome, component int) {
	for i := 0; i < len(p.Components[component].Nodes); i++ {
		nI := p.Components[component].Nodes[i].Node
		c.Nodes[nI] = p.Nodes[nI]
	}
}

func (pgo *PcbGeneticOperators) CrossOver(i1 *Pcb, i2 *Pcb, c *genetic.GeneticContext) *Pcb {
	child := i1.Genome.copy()

	for i := 0; i < len(i1.Genome.Components); i++ {
		v := c.RandomGenerator.Float64()
		if v < 0.5 {
			pgo.copyComponentNodesToChild(child, i1.Genome, i)
			child.Components[i] = *i1.Genome.Components[i].copy()
		} else {
			pgo.copyComponentNodesToChild(child, i2.Genome, i)
			child.Components[i] = *i2.Genome.Components[i].copy()
		}
	}

	currentNet := i1.Genome.Edges[0].Net
	v := c.RandomGenerator.Float64()

	// Suppose both parent edges are sorted by net and omologous nets of both parents have same size
	for i, e1 := range i1.Genome.Edges {
		e2 := i2.Genome.Edges[i]

		if e1.Net != currentNet {
			v = c.RandomGenerator.Float64()
			currentNet = e1.Net
		}

		if v < 0.5 {
			child.Edges[i] = e1
		} else {
			child.Edges[i] = e2
		}

	}

	return NewPcb(child)
}

func (pgo *PcbGeneticOperators) globalMutation(i *Pcb, c *genetic.GeneticContext) {
	for j := range i.Genome.Components {
		if c.RandomGenerator.Float64() < pgo.mutateSinglePointProb {
			component := &i.Genome.Components[j]
			CX, CY := GetComponentRandomPositionInBoundaries(component, pgo.maxX, pgo.maxY, c.RandomGenerator)
			component.Rotation = c.RandomGenerator.Float64() * 360
			MoveComponent(i.Genome.Nodes, component, CX, CY)
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

func (pgo *PcbGeneticOperators) Mutate(i *Pcb, c *genetic.GeneticContext) {
	chooser, _ := weightedrand.NewChooser(
		weightedrand.NewChoice(0, int(pgo.localMutationProb*100)),
		weightedrand.NewChoice(1, int(pgo.netMutationProb*100)),
		weightedrand.NewChoice(2, int((1.0-(pgo.netMutationProb+pgo.localMutationProb))*100)),
	)

	if c.RandomGenerator.Float64() < pgo.mutateProb {
		chosen := chooser.Pick()

		switch chosen {
		case 0:
			// pgo.localMutation(i, c)
			// fmt.Println("LOCAL")
		case 1:
			pgo.netMutation(i, c)
			// fmt.Println("NET")
		case 2:
			pgo.globalMutation(i, c)
			// fmt.Println("GLOBAL")
		}
	}

}

func (pgo *PcbGeneticOperators) Grow(i *Pcb, c *genetic.GeneticContext) {
	i.ComputeGeometry(pgo.nodeSz, pgo.edgeSz)
}
