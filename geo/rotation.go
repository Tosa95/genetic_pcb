package geo

import "math"

func RotatePoint(x, y, centerX, centerY, angle float64) (float64, float64) {
	// Translate the point relative to the center
	translatedX := x - centerX
	translatedY := y - centerY

	// Convert angle from degrees to radians
	radians := angle * math.Pi / 180

	// Rotate the point
	rotatedX := translatedX*math.Cos(radians) - translatedY*math.Sin(radians)
	rotatedY := translatedX*math.Sin(radians) + translatedY*math.Cos(radians)

	// Translate the point back to its original position relative to the center
	rotatedX += centerX
	rotatedY += centerY

	return rotatedX, rotatedY
}
