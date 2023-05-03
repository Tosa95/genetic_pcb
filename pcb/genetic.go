package pcb

import (
	"genetic_pcb/genetic"
	"math"
)

type PcbGeneticOperators struct {
	fitnessExp                float64
	mutateProb                float64
	mutateSingleComponentProb float64
	maxX                      float64
	maxY                      float64
	nodeSz                    float64
	edgeSz                    float64
	localMutationMaxDelta     float64
	mutationWeights           MutationWeights
	evaluationParams          EvaluationParams
	mutationChooser           mutationChooser
}

func NewPcbGeneticOperators(
	fitnessExp float64,
	mutateProb float64,
	mutateSinglePointProb float64,
	maxX float64,
	maxY float64,
	nodeSz float64,
	edgeSz float64,
	localMutationMaxDelta float64,
	mutationWeights MutationWeights,
	evaluationParams EvaluationParams,
) *PcbGeneticOperators {

	pgo := PcbGeneticOperators{
		fitnessExp:                fitnessExp,
		mutateProb:                mutateProb,
		mutateSingleComponentProb: mutateSinglePointProb,
		maxX:                      maxX,
		maxY:                      maxY,
		nodeSz:                    nodeSz,
		edgeSz:                    edgeSz,
		localMutationMaxDelta:     localMutationMaxDelta,
		mutationWeights:           mutationWeights,
		evaluationParams:          evaluationParams,
	}

	pgo.mutationChooser = *pgo.buildMutationChooser()

	return &pgo
}

func (pgo *PcbGeneticOperators) Evaluate(i *Pcb, c *genetic.GeneticContext) float64 {
	cost := pgo.EvaluatePcbIntersections(i)
	cost += pgo.EvaluatePcbEdgeLengths(i)
	cost += pgo.EvaluateNonZeroPlaneEdges(i)
	cost += pgo.EvaluateComponentsOutOfBounds(i)
	cost = math.Pow(cost, pgo.fitnessExp)
	fitness := -cost
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

func (pgo *PcbGeneticOperators) Mutate(i *Pcb, c *genetic.GeneticContext) {
	if c.RandomGenerator.Float64() < pgo.mutateProb {
		mutation := pgo.mutationChooser.Pick()
		mutation(i, c)
	}
}

func (pgo *PcbGeneticOperators) Grow(i *Pcb, c *genetic.GeneticContext) {
	i.ComputeGeometry(pgo.nodeSz, pgo.edgeSz)
}
