package geo

import (
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy"
)

func isAtLeastOnePointContained(p1 *geom.Polygon, p2 *geom.Polygon) bool {
	for _, c := range p1.Coords() {
		for _, pt := range c {
			if xy.IsPointInRing(geom.XY, pt, p2.FlatCoords()) {
				return true
			}
		}
	}
	return false
}

func IsContained(p1 *geom.Polygon, p2 *geom.Polygon) bool {
	for _, c := range p1.Coords() {
		for _, pt := range c {
			if !(xy.IsPointInRing(geom.XY, pt, p2.FlatCoords())) {
				return false
			}
		}
	}
	return true
}

func LineToPolyDist(lstart geom.Coord, lend geom.Coord, poly *geom.Polygon) float64 {
	min := -1.0

	for _, c := range poly.Coords() {
		prev := c[0]
		for _, curr := range c[1:] {
			distance := xy.DistanceFromLineToLine(lstart, lend, prev, curr)

			if distance == 0 {
				return 0
			}

			if min < 0 || distance < min {
				min = distance
			}

			prev = curr
		}
	}

	return min
}

func PointToPolyDist(pt geom.Coord, poly *geom.Polygon) float64 {
	min := -1.0

	if xy.IsPointInRing(geom.XY, pt, poly.Clone().FlatCoords()) {
		return 0
	}

	for _, c := range poly.Coords() {
		prev := c[0]
		for _, curr := range c[1:] {
			distance := xy.DistanceFromPointToLine(pt, prev, curr)

			if distance == 0 {
				return 0
			}

			if min < 0 || distance < min {
				min = distance
			}

			prev = curr
		}
	}

	return min
}

func PolyDistance(p1 *geom.Polygon, p2 *geom.Polygon) float64 {
	min := -0.1

	if isAtLeastOnePointContained(p1, p2) || isAtLeastOnePointContained(p2, p1) {
		return 0
	}

	for _, c := range p1.Coords() {
		prev := c[0]
		for _, curr := range c[1:] {
			distance := LineToPolyDist(prev, curr, p2)

			if distance == 0 {
				return 0
			}

			if min < 0 || distance < min {
				min = distance
			}

			prev = curr
		}
	}

	return min
}
