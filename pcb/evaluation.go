package pcb

import (
	"genetic_pcb/geo"

	"github.com/twpayne/go-geom"
)

type EvaluationParams struct {
	SamePlaneIntersectionCost      float64
	DifferentPlaneIntersectionCost float64
	EdgeLengthCost                 float64
	OutOfBoundsCost                float64
	NonZeroPlaneEdgeCost           float64
	MinDist                        float64
}

func (pgo *PcbGeneticOperators) EvaluatePcbIntersections(pcb *Pcb) float64 {

	cost := 0.0

	for i1, n1 := range pcb.Geometry.Nodes {
		for i2, n2 := range pcb.Geometry.Nodes {
			if i1 != i2 && pcb.Genome.Nodes[i1].Component != pcb.Genome.Nodes[i2].Component {
				if geo.PolyDistance(n1, n2) < pgo.evaluationParams.MinDist {
					// fmt.Printf("nn %v %v\n", i1, i2)
					cost += pgo.evaluationParams.SamePlaneIntersectionCost / 2.0
				}
			}
		}
	}

	for i1, e1 := range pcb.Geometry.Edges {
		for i2, e2 := range pcb.Geometry.Edges {
			if i1 != i2 {
				if !(pcb.Genome.AreAdjacent(i1, i2)) && geo.PolyDistance(e1, e2) < pgo.evaluationParams.MinDist {
					if pcb.Genome.Edges[i1].Plane == pcb.Genome.Edges[i2].Plane {
						cost += pgo.evaluationParams.SamePlaneIntersectionCost / 2.0
					} else {
						cost += pgo.evaluationParams.DifferentPlaneIntersectionCost / 2.0
					}

				}
			}
		}
	}

	for i1, n := range pcb.Geometry.Nodes {
		for i2, e := range pcb.Geometry.Edges {
			if !(pcb.Genome.IsNodeOnEdge(i1, i2)) && geo.PolyDistance(n, e) < pgo.evaluationParams.MinDist {
				// fmt.Printf("ne %v %v %v %v\n", i1, i2, geo.PolyDistance(n, e), pcb.Genome.Edges[i2])
				cost += pgo.evaluationParams.SamePlaneIntersectionCost
			}
		}
	}

	for i1, c1 := range pcb.Geometry.Components {
		for i2, c2 := range pcb.Geometry.Components {
			if i1 != i2 && geo.PolyDistance(c1, c2) < pgo.evaluationParams.MinDist {
				// fmt.Printf("ne %v %v %v %v\n", i1, i2, geo.PolyDistance(n, e), pcb.Genome.Edges[i2])
				cost += pgo.evaluationParams.SamePlaneIntersectionCost
			}
		}
	}

	return cost

}

func (pgo *PcbGeneticOperators) EvaluatePcbEdgeLengths(pcb *Pcb) float64 {
	return (GetTotalPcbLength(pcb) / (pgo.maxX + pgo.maxY)) * pgo.evaluationParams.EdgeLengthCost
}

func (pgo *PcbGeneticOperators) getNonZeroPlaneEdgesCount(i *Pcb) int {
	cost := 0

	for _, e := range i.Genome.Edges {
		if e.Plane != 0 {
			cost += 1
		}
	}

	return cost
}

func (pgo *PcbGeneticOperators) EvaluateNonZeroPlaneEdges(pcb *Pcb) float64 {
	return float64(pgo.getNonZeroPlaneEdgesCount(pcb)) * pgo.evaluationParams.NonZeroPlaneEdgeCost
}

func (pgo *PcbGeneticOperators) EvaluateComponentsOutOfBounds(pcb *Pcb) float64 {

	cost := 0.0

	pcbSpace := geom.NewPolygonFlat(geom.XY, []float64{
		0, 0,
		pgo.maxX, 0,
		pgo.maxX, pgo.maxY,
		0, pgo.maxY,
		0, 0,
	}, []int{10})

	for _, c := range pcb.Geometry.Components {
		if !geo.IsContained(c, pcbSpace) {
			cost += pgo.evaluationParams.OutOfBoundsCost
		}
	}

	return cost
}
