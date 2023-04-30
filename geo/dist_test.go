package geo_test

import (
	"genetic_pcb/geo"
	"testing"

	"github.com/twpayne/go-geom"
)

func TestPolyDistance1(t *testing.T) {
	p1 := geom.NewPolygonFlat(geom.XY, []float64{0, 0, 10, 0, 10, 1, 0, 1, 0, 0}, []int{8})
	p2 := geom.NewPolygonFlat(geom.XY, []float64{5, 5, 5, 2, 6, 2, 6, 5, 5, 5}, []int{8})

	dist := geo.PolyDistance(p1, p2)

	if dist != 1 {
		t.Errorf("Expected distance = 1, got %v", dist)
	}
}

func TestPolyDistance2(t *testing.T) {
	p1 := geom.NewPolygonFlat(geom.XY, []float64{0, 0, 10, 0, 10, 1, 0, 1, 0, 0}, []int{8})
	p2 := geom.NewPolygonFlat(geom.XY, []float64{5, 5, 5, -2, 6, -2, 6, 5, 5, 5}, []int{8})

	dist := geo.PolyDistance(p1, p2)

	if dist != 0 {
		t.Errorf("Expected distance = 0, got %v", dist)
	}
}

func TestPolyDistance3(t *testing.T) {
	p1 := geom.NewPolygonFlat(geom.XY, []float64{0, 0, 10, 0, 10, 10, 0, 10, 0, 0}, []int{8})
	p2 := geom.NewPolygonFlat(geom.XY, []float64{1, 1, 2, 1, 2, 2, 1, 2, 1, 1}, []int{8})

	dist := geo.PolyDistance(p1, p2)

	if dist != 0 {
		t.Errorf("Expected distance = 0, got %v", dist)
	}
}

func TestPolyDistance4(t *testing.T) {
	p1 := geom.NewPolygonFlat(geom.XY, []float64{0.0, 0.0, 1.0, 0.0, 1.0, 1.0, 0.0, 1.0, 0.0, 0.0}, []int{8})
	p2 := geom.NewPolygonFlat(geom.XY, []float64{2.1, 0.5, 2.6, 0.0, 3.1, 0.5, 2.6, 1.0, 2.1, 0.5}, []int{8})

	dist := geo.PolyDistance(p1, p2)

	if dist != 1.1 {
		t.Errorf("Expected distance = 1.1, got %v", dist)
	}
}
