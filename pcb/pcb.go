package pcb

import (
	"fmt"
	"genetic_pcb/draw"
	"genetic_pcb/geo"
	"image"
	"image/color"
	"math"
	"math/rand"
	"sort"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/twpayne/go-geom"
)

type Node struct {
	X float64
	Y float64
}

type Edge struct {
	From int
	To   int
	Net  int
}

type Net struct {
	Nodes []int
}

type ComponentNode struct {
	Node int
	DX   float64
	DY   float64
}

type Component struct {
	Nodes []ComponentNode
	X1    float64
	Y1    float64
	X2    float64
	Y2    float64
	CX    float64
	CY    float64
}

type Genome struct {
	Nodes      []Node
	Edges      []Edge
	Nets       []Net
	Components []Component
}

type Geometry struct {
	Nodes      []*geom.Polygon
	Edges      []*geom.Polygon
	Components []*geom.Polygon
}

type Pcb struct {
	Genome   *Genome
	Geometry *Geometry
}

func (p *Pcb) String() string {
	return fmt.Sprintf("%v", p.Genome.Nodes)
}

func (g *Genome) copy() *Genome {
	nodes := make([]Node, len(g.Nodes))
	copy(nodes, g.Nodes)

	edges := make([]Edge, len(g.Edges))
	copy(edges, g.Edges)

	nets := make([]Net, len(g.Nets))
	copy(nets, g.Nets)

	components := make([]Component, len(g.Components))
	copy(components, g.Components)

	return &Genome{
		Nodes:      nodes,
		Edges:      edges,
		Nets:       nets,
		Components: components,
	}
}

func (n *Net) copy() *Net {
	nodes := make([]int, len(n.Nodes))
	copy(nodes, n.Nodes)

	return &Net{
		Nodes: nodes,
	}
}

func nodeToPoly(x, y, nodeSz float64) *geom.Polygon {
	halfSz := nodeSz / 2
	return geom.NewPolygonFlat(geom.XY, []float64{x - halfSz, y - halfSz, x + halfSz, y - halfSz, x + halfSz, y + halfSz, x - halfSz, y + halfSz}, []int{8})
}

func componentToPoly(c *Component) *geom.Polygon {
	finalX1 := c.CX + c.X1
	finalY1 := c.CY + c.Y1
	finalX2 := c.CX + c.X2
	finalY2 := c.CY + c.Y2

	return geom.NewPolygonFlat(geom.XY, []float64{
		finalX1, finalY1,
		finalX2, finalY1,
		finalX2, finalY2,
		finalX1, finalY2,
	}, []int{8})
}

func edgeToPoly(from, to int, nodes []Node, edgeSz float64) *geom.Polygon {
	// https://www.quora.com/How-do-I-find-a-vector-orthogonal-to-another-vector#:~:text=One%20way%20to%20find%20a,%C2%B7%20u%20%3D%200%20for%20u.

	x1, y1 := nodes[from].X, nodes[from].Y
	x2, y2 := nodes[to].X, nodes[to].Y
	xs, ys := x2-x1, y2-y1
	xp := 1.0
	yp := -(xs / ys) * xp
	pNorm := math.Sqrt(xp*xp + yp*yp)

	xp = (xp / pNorm) * (edgeSz / 2)
	yp = (yp / pNorm) * (edgeSz / 2)

	return geom.NewPolygonFlat(geom.XY, []float64{x1 - xp, y1 - yp, x2 - xp, y2 - yp, x2 + xp, y2 + yp, x1 + xp, y1 + yp, x1 - xp, y1 - yp}, []int{10})
}

func NewPcb(genome *Genome) *Pcb {
	res := Pcb{
		Genome: genome,
	}

	return &res
}

func (pcb *Pcb) ComputeGeometry(nodeSz, edgeSz float64) {
	nodes := make([]*geom.Polygon, len(pcb.Genome.Nodes))
	edges := make([]*geom.Polygon, len(pcb.Genome.Edges))
	components := make([]*geom.Polygon, len(pcb.Genome.Components))

	for i, coords := range pcb.Genome.Nodes {
		x, y := coords.X, coords.Y
		nodePoly := nodeToPoly(x, y, nodeSz)
		nodes[i] = nodePoly
	}

	for i, edge := range pcb.Genome.Edges {
		edgePoly := edgeToPoly(edge.From, edge.To, pcb.Genome.Nodes, edgeSz)
		edges[i] = edgePoly
	}

	for i, component := range pcb.Genome.Components {
		componentPoly := componentToPoly(&component)
		components[i] = componentPoly
	}

	geometry := Geometry{
		Nodes:      nodes,
		Edges:      edges,
		Components: components,
	}

	pcb.Geometry = &geometry
}

func (g *Genome) AreAdjacent(edgeIndex1, edgeIndex2 int) bool {
	f1, t1 := g.Edges[edgeIndex1].From, g.Edges[edgeIndex1].To
	f2, t2 := g.Edges[edgeIndex2].From, g.Edges[edgeIndex2].To

	return f1 == f2 || f1 == t2 || t1 == f2 || t1 == t2
}

func (g *Genome) IsNodeOnEdge(nodeIndex, edgeIndex int) bool {
	f, t := g.Edges[edgeIndex].From, g.Edges[edgeIndex].To

	return f == nodeIndex || t == nodeIndex
}

func DrawPcbToImage(pcb *Pcb, imgPath string, imgW, imgH int, sx, sy float64, netColors []color.Color) {

	img := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
	gc := draw2dimg.NewGraphicContext(img)

	gc.SetLineWidth(0)

	for i, edge := range pcb.Geometry.Edges {
		color := netColors[pcb.Genome.Edges[i].Net%len(netColors)]
		gc.SetFillColor(color)
		draw.DrawPoly(gc, edge, sx, sy, true)
	}

	gc.SetFillColor(color.RGBA{255, 0, 0, 255})

	for _, node := range pcb.Geometry.Nodes {
		draw.DrawPoly(gc, node, sx, sy, true)
	}

	gc.SetFillColor(color.RGBA{255, 255, 255, 0})
	gc.SetLineWidth(1)
	gc.SetStrokeColor(color.RGBA{255, 255, 255, 255})

	for _, component := range pcb.Geometry.Components {
		draw.DrawPoly(gc, component, sx, sy, false)
	}

	draw2dimg.SaveToPngFile(imgPath, img)

}

func EvaluatePcb(pcb *Pcb, minDist float64) float64 {

	violatedConstraints := 0.0

	for i1, n1 := range pcb.Geometry.Nodes {
		for i2, n2 := range pcb.Geometry.Nodes {
			if i1 != i2 {
				if geo.PolyDistance(n1, n2) < minDist {
					// fmt.Printf("nn %v %v\n", i1, i2)
					violatedConstraints += 0.5
				}
			}
		}
	}

	for i1, e1 := range pcb.Geometry.Edges {
		for i2, e2 := range pcb.Geometry.Edges {
			if i1 != i2 {
				if !(pcb.Genome.AreAdjacent(i1, i2)) && geo.PolyDistance(e1, e2) < minDist {
					// fmt.Printf("ee %v %v %v\n", i1, i2, geo.PolyDistance(e1, e2))
					violatedConstraints += 0.5
				}
			}
		}
	}

	for i1, n := range pcb.Geometry.Nodes {
		for i2, e := range pcb.Geometry.Edges {
			if !(pcb.Genome.IsNodeOnEdge(i1, i2)) && geo.PolyDistance(n, e) < minDist {
				// fmt.Printf("ne %v %v\n", i1, i2)
				violatedConstraints += 1.0
			}
		}
	}

	return -violatedConstraints

}

func GetTotalPcbLength(pcb *Pcb) float64 {
	res := 0.0

	for _, edge := range pcb.Genome.Edges {
		x1, y1 := pcb.Genome.Nodes[edge.From].X, pcb.Genome.Nodes[edge.From].Y
		x2, y2 := pcb.Genome.Nodes[edge.To].X, pcb.Genome.Nodes[edge.To].Y

		res += math.Sqrt((x1-x2)*(x1-x2) + (y1-y2)*(y1-y2))
	}

	return res
}

func RemoveNetEdges(pcb *Pcb, net int) {
	i := 0 // output index
	for _, edge := range pcb.Genome.Edges {
		if edge.Net != net {
			// copy and increment index
			pcb.Genome.Edges[i] = edge
			i++
		}
	}

	pcb.Genome.Edges = pcb.Genome.Edges[:i]
}

func GenerateNet(pcb *Pcb, net int, randomGenerator *rand.Rand) {
	RemoveNetEdges(pcb, net)

	netInstance := pcb.Genome.Nets[net]
	nNodes := len(netInstance.Nodes)

	shuffledNodes := make([]int, nNodes)

	copy(shuffledNodes, netInstance.Nodes)

	rand.Shuffle(nNodes, func(i, j int) { shuffledNodes[i], shuffledNodes[j] = shuffledNodes[j], shuffledNodes[i] })

	for i := 1; i < nNodes; i++ {
		connectToIndex := randomGenerator.Intn(i)
		to := shuffledNodes[connectToIndex]
		from := shuffledNodes[i]

		edge := Edge{
			From: from,
			To:   to,
			Net:  net,
		}

		pcb.Genome.Edges = append(pcb.Genome.Edges, edge)
	}

	SortPcbEdges(pcb)
}

func SortPcbEdges(pcb *Pcb) {
	sort.Slice(pcb.Genome.Edges, func(i, j int) bool {
		e1 := pcb.Genome.Edges[i]
		e2 := pcb.Genome.Edges[j]

		if e1.Net == e2.Net && e1.From == e2.From {
			return e1.To < e2.To
		}

		if e1.Net == e2.Net {
			return e1.From < e2.From
		}

		return e1.Net < e2.Net
	})
}
