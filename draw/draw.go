package draw

import (
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/twpayne/go-geom"
)

func DrawPoly(gc *draw2dimg.GraphicContext, p *geom.Polygon, sx float64, sy float64, fill bool) {

	ring := p.Coords()[0]

	firstPoint := ring[0]
	originalMatrix := gc.GetMatrixTransform()

	gc.SetMatrixTransform(draw2d.NewScaleMatrix(sx, sy))

	gc.BeginPath()
	// fmt.Printf("-%v %v\n", firstPoint.X(), firstPoint.Y())
	gc.MoveTo(firstPoint.X(), firstPoint.Y())

	for _, pt := range ring[1:] {
		// fmt.Printf("%v %v\n", pt.X(), pt.Y())
		gc.LineTo(pt.X(), pt.Y())
	}

	gc.Close()
	if fill {
		gc.Fill()
	} else {
		gc.Stroke()
	}

	gc.SetMatrixTransform(originalMatrix)

}
