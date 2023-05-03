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

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/twpayne/go-geom"

	"github.com/golang/freetype"
	"golang.org/x/image/font/gofont/goregular"
)

type Node struct {
	X         float64
	Y         float64
	Component int
}

type Edge struct {
	From  int
	To    int
	Net   int
	Plane int
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
	Nodes    []ComponentNode
	X1       float64
	Y1       float64
	X2       float64
	Y2       float64
	CX       float64
	CY       float64
	Rotation float64
}

func (c *Component) copy() *Component {
	nodes := make([]ComponentNode, len(c.Nodes))
	copy(nodes, c.Nodes)

	return &Component{
		Nodes:    nodes,
		X1:       c.X1,
		Y1:       c.Y1,
		X2:       c.X2,
		Y2:       c.Y2,
		CX:       c.CX,
		CY:       c.CY,
		Rotation: c.Rotation,
	}
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

	for i := range components {
		components[i] = *g.Components[i].copy()
	}

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
	return geom.NewPolygonFlat(geom.XY, []float64{x - halfSz, y - halfSz, x + halfSz, y - halfSz, x + halfSz, y + halfSz, x - halfSz, y + halfSz, x - halfSz, y - halfSz}, []int{10})
}

func componentToPoly(c *Component) *geom.Polygon {
	finalX1 := c.CX + c.X1
	finalY1 := c.CY + c.Y1
	finalX2 := c.CX + c.X2
	finalY2 := c.CY + c.Y2

	p1X, p1Y := geo.RotatePoint(finalX1, finalY1, c.CX, c.CY, c.Rotation)
	p2X, p2Y := geo.RotatePoint(finalX2, finalY1, c.CX, c.CY, c.Rotation)
	p3X, p3Y := geo.RotatePoint(finalX2, finalY2, c.CX, c.CY, c.Rotation)
	p4X, p4Y := geo.RotatePoint(finalX1, finalY2, c.CX, c.CY, c.Rotation)

	return geom.NewPolygonFlat(geom.XY, []float64{
		p1X, p1Y,
		p2X, p2Y,
		p3X, p3Y,
		p4X, p4Y,
		p1X, p1Y,
	}, []int{10})
}

func edgeToPoly(from, to int, nodes []Node, edgeSz float64) *geom.Polygon {
	// https://www.quora.com/How-do-I-find-a-vector-orthogonal-to-another-vector#:~:text=One%20way%20to%20find%20a,%C2%B7%20u%20%3D%200%20for%20u.

	x1, y1 := nodes[from].X, nodes[from].Y
	x2, y2 := nodes[to].X, nodes[to].Y
	xs, ys := x2-x1, y2-y1
	xp := 1.0
	yp := 0.0
	if ys != 0 {
		yp = -(xs / ys) * xp
	} else {
		xp = 0.0
		yp = 1.0
	}

	pNorm := math.Sqrt(xp*xp + yp*yp)

	if pNorm == 0 {
		return geom.NewPolygonFlat(geom.XY, []float64{x1, y1, x2, y1, x2, y2, x1, y2, x1, y1}, []int{10})
	}
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

func reduceLuminosity(originalColor color.Color) color.Color {
	// Get the RGBA values of the original color
	originalR, originalG, originalB, originalA := originalColor.RGBA()

	// Convert the uint32 RGBA values to uint8 values
	r := uint8(originalR / 0x101)
	g := uint8(originalG / 0x101)
	b := uint8(originalB / 0x101)
	a := uint8(originalA / 0x101)

	// Calculate the new RGB values by reducing the luminosity by 50%
	newR := uint8(math.Max(float64(r)/2, 0))
	newG := uint8(math.Max(float64(g)/2, 0))
	newB := uint8(math.Max(float64(b)/2, 0))

	// Create the new color with the updated RGB values
	return color.RGBA{newR, newG, newB, a}
}

func DrawPcbToImage(pcb *Pcb, imgPath string, imgW, imgH int, sx, sy float64, netColors []color.Color) {

	img := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
	gc := draw2dimg.NewGraphicContext(img)

	gc.SetLineWidth(0)

	fontData := goregular.TTF
	font, err := freetype.ParseFont(fontData)
	if err != nil {
		panic(err)
	}

	gc.SetFont(font)
	gc.SetFontSize(50)

	for i, edge := range pcb.Geometry.Edges {
		cl := netColors[pcb.Genome.Edges[i].Net%len(netColors)]

		if pcb.Genome.Edges[i].Plane == 0 {
			gc.SetFillColor(cl)
		} else {
			gc.SetFillColor(reduceLuminosity(cl))
			fmt.Println("HEHEHE")
		}

		draw.DrawPoly(gc, edge, sx, sy, true)
	}

	for i, node := range pcb.Geometry.Nodes {
		gc.SetFillColor(color.RGBA{255, 0, 0, 255})
		draw.DrawPoly(gc, node, sx, sy, true)
		gc.SetFillColor(color.RGBA{200, 200, 200, 255})
		// gc.FontCache = yourFontCache
		gc.FontCache = draw2d.NewFolderFontCache("C:\\Users\\david\\st\\prj\\genetic_pcb\\resource\\font")
		// gc.FontCache = draw2d.GetGlobalFontCache()
		gc.SetFontSize(10)
		gc.FillStringAt(fmt.Sprintf("%v", i), pcb.Genome.Nodes[i].X, pcb.Genome.Nodes[i].Y)
	}

	gc.SetFillColor(color.RGBA{255, 255, 255, 0})
	gc.SetLineWidth(1)
	gc.SetStrokeColor(color.RGBA{255, 255, 255, 255})

	for _, component := range pcb.Geometry.Components {
		draw.DrawPoly(gc, component, sx, sy, false)
	}

	draw2dimg.SaveToPngFile(imgPath, img)

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

func PlaceComponentNodes(nodes []Node, c *Component) {
	for _, n := range c.Nodes {
		nodes[n.Node].X = c.CX + n.DX
		nodes[n.Node].Y = c.CY + n.DY

		nodes[n.Node].X, nodes[n.Node].Y = geo.RotatePoint(nodes[n.Node].X, nodes[n.Node].Y, c.CX, c.CY, c.Rotation)
	}
}

func GetComponentRandomPositionInBoundaries(c *Component, maxX, maxY float64, randomGenerator *rand.Rand) (float64, float64) {
	maxXDisplacement := (maxX - (c.X2 - c.CX) - (c.CX - c.X1))
	maxYDisplacement := (maxY - (c.Y2 - c.CY) - (c.CY - c.Y1))

	return randomGenerator.Float64()*maxXDisplacement - c.X1, randomGenerator.Float64()*maxYDisplacement - c.Y1
}
