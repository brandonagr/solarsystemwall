package solar

// DrawCircle renders an expanding circle
type DrawCircle struct {

	// center position of the ring
	position XYPosition

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
func NewCircle(system *SolarSystem) *DrawCircle {

	return &DrawCircle{
		position: XYPosition{0, 0},
		zindex:   100,
	}
}

// ColorAt Returns the color at position blended on top of baseColor
func (circle *DrawCircle) ColorAt(position XYPosition, baseColor RGBA) (color RGBA) {

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
