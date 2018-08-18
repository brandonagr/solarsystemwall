package solar

import (
	"image/color"
	"math"

	"github.com/golang/geo/r2"
)

// DrawLine renders a line that moves back and forth
type DrawLine struct {

	// where line starts
	startPosition r2.Point

	// where line neds
	endPosition r2.Point

	// time it takes for line to travel from start to end
	traverseTime float64

	// position of the line along movement vector, 0 to 1
	currentPosition r2.Point

	// normal direction of the line
	lineDirection r2.Point

	// width of the line
	lineWidth float64

	// color of the line
	color color.RGBA

	// z position of ball
	zindex ZIndex
}

var _ Drawable = &DrawLine{}

// NewLine Construct a circle
func NewLine(solarSystem *System) *DrawLine {

	return &DrawLine{
		startPosition:   r2.Point{X: 0, Y: 0},
		endPosition:     r2.Point{X: 130, Y: 0},
		traverseTime:    8.0,
		currentPosition: r2.Point{X: 0, Y: 0},
		lineDirection:   r2.Point{X: 1, Y: 0},
		lineWidth:       12.0,
		color:           color.RGBA{R: 255, G: 255, B: 0, A: 255},
		zindex:          1,
	}
}

// Affects returns bounding circle check
func (line *DrawLine) Affects(position r2.Point, radius float64) bool {
	distance := line.distanceToPoint(position)
	return (distance > line.lineWidth+radius)
}

// ColorAt Returns the color at position blended on top of baseColor
func (line *DrawLine) ColorAt(position r2.Point, baseColor RGBA) (color RGBA) {

	distance := line.distanceToPoint(position)
	if distance > line.lineWidth {
		return baseColor
	}
	distance = distance / line.lineWidth
	color = RGBA{line.color.R, line.color.G, line.color.B, uint8((1.0 - distance) * 255.0)}

	result := color.BlendWith(baseColor)

	//fmt.Println(baseColor, color, result, position, distance, line.currentPosition)

	return result
}

// Computer distance from point to line https://brilliant.org/wiki/dot-product-distance-between-point-and-a-line/
func (line *DrawLine) distanceToPoint(position r2.Point) float64 {
	return math.Abs(position.Sub(line.currentPosition).Dot(line.lineDirection))
}

// ZIndex of the circle
func (line *DrawLine) ZIndex() ZIndex {
	return line.zindex
}

// Animate circle
func (line *DrawLine) Animate(dt float64) bool {
	totalDistance := line.endPosition.Sub(line.startPosition).Norm()
	distance := line.currentPosition.Sub(line.startPosition).Norm()

	distance += (totalDistance / line.traverseTime) * dt
	//fmt.Println(distance, totalDistance, line.traverseTime)
	if distance > totalDistance {
		distance = 0.0
	}

	moveDir := line.endPosition.Sub(line.startPosition).Normalize()
	line.currentPosition = moveDir.Mul(distance).Add(line.startPosition)

	//fmt.Println("Animated to ", line.currentPosition, dt)

	return true
}
