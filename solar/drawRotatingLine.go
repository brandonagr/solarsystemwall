package solar

import (
	"image/color"
	"math"

	"github.com/golang/geo/r2"
)

// DrawRotatingLine renders a line that moves back and forth
type DrawRotatingLine struct {

	// where line starts
	startPosition r2.Point

	// end of the line, updated after call to Animate
	endPosition r2.Point

	// where line neds
	length float64

	// time it takes for line to travel from start to end
	traverseTime float64

	// angle of the line
	currentAngle float64

	// width of the line
	lineWidth float64

	// color of the line
	color color.RGBA

	// z position of line
	zindex ZIndex
}

var _ Drawable = &DrawLine{}

// NewRotatingLine Construct a circle
func NewRotatingLine(planet PlanetIndex, solarSystem *System) *DrawRotatingLine {

	return &DrawRotatingLine{
		startPosition: solarSystem.planets[planet].position,
		length:        7.0,
		traverseTime:  12.0,
		currentAngle:  0.0,
		lineWidth:     2.0,
		color:         color.RGBA{R: 255, G: 255, B: 255, A: 255},
		zindex:        1,
	}
}

// Affects returns bounding circle check
func (line *DrawRotatingLine) Affects(position r2.Point, radius float64) bool {
	distance := line.startPosition.Sub(position).Norm()

	//fmt.Println(distance, position, radius)
	return (distance > line.length+radius)
}

// ColorAt Returns the color at position blended on top of baseColor
func (line *DrawRotatingLine) ColorAt(position r2.Point, baseColor RGBA) (color RGBA) {

	distance := line.distanceToPoint(position)

	//fmt.Println(distance, line.endPosition, position)

	if distance > line.lineWidth {
		return baseColor
	}
	distance = distance / line.lineWidth
	color = RGBA{line.color.R, line.color.G, line.color.B, uint8((1.0 - distance) * 255.0)}

	result := color.BlendWith(baseColor)

	//fmt.Println(baseColor, color, result, position, distance)

	return result
}

// Computer distance from point to line segment https://stackoverflow.com/questions/849211/shortest-distance-between-a-point-and-a-line-segment
func (line *DrawRotatingLine) distanceToPoint(position r2.Point) float64 {

	length := line.startPosition.Sub(line.endPosition).Norm()
	length = length * length
	if length == 0 {
		return line.startPosition.Sub(position).Norm()
	}

	startToEndVector := line.endPosition.Sub(line.startPosition)
	t := (position.Sub(line.startPosition).Dot(startToEndVector)) / length
	if t > 1 {
		t = 1
	}
	if t < 0 {
		t = 0
	}

	projectedPoint := line.startPosition.Add(startToEndVector.Mul(t))

	return position.Sub(projectedPoint).Norm()
}

// ZIndex of the circle
func (line *DrawRotatingLine) ZIndex() ZIndex {
	return line.zindex
}

// Animate circle
func (line *DrawRotatingLine) Animate(dt float64) bool {
	line.currentAngle += 2 * math.Pi * dt / line.traverseTime
	if line.currentAngle > 2*math.Pi {
		line.currentAngle -= 2 * math.Pi
	}

	endOffset := r2.Point{X: math.Cos(line.currentAngle) * line.length, Y: math.Sin(line.currentAngle) * line.length}
	line.endPosition = line.startPosition.Add(endOffset)

	//fmt.Println("Animated to ", line.currentAngle, line.endPosition, dt)

	return true
}
