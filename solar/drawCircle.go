package solar

import (
	"github.com/golang/geo/r2"
)

// DrawCircle renders an expanding circle
type DrawCircle struct {

	// center position of the ring
	position r2.Point

	// speed of radius change in distance per second
	velocity float64

	// radius of the ring
	radius float64

	// the length of the tail of the ball
	maxRadius float64

	// if the ball should be hidden this frame or not
	hideBall bool

	// rate that influence falls off
	falloff float64

	//
	color RGBA

	// z position of ball
	zindex ZIndex
}

var _ Drawable = &DrawCircle{}

// NewCircle Construct a circle
func NewCircle(solarSystem *System) *DrawCircle {

	return &DrawCircle{
		position: r2.Point{X: 0, Y: 0},
		zindex:   100,
	}
}

// Affects returns bounding circle check
func (circle *DrawCircle) Affects(position r2.Point, radius float64) bool {
	distance := circle.position.Sub(position)
	return distance.Norm() < radius+circle.maxRadius
}

// ColorAt Returns the color at position blended on top of baseColor
func (circle *DrawCircle) ColorAt(position r2.Point, baseColor RGBA) (color RGBA) {

	return color
}

// ZIndex of the circle
func (circle *DrawCircle) ZIndex() ZIndex {
	return circle.zindex
}

// Animate circle
func (circle *DrawCircle) Animate(dt float64) bool {
	circle.radius += circle.velocity * dt

	return (circle.radius < circle.maxRadius)
}
