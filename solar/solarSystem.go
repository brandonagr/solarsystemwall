package solar

import (
	"container/list"
	"math"

	"github.com/golang/geo/r2"
)

// PlanetIndex used to reference array of all planets
type PlanetIndex int

// Planet indexes
const (
	Sun         PlanetIndex = 0
	Mercury     PlanetIndex = 1
	Venus       PlanetIndex = 2
	Earth       PlanetIndex = 3
	Mars        PlanetIndex = 4
	Jupiter     PlanetIndex = 5
	Saturn      PlanetIndex = 6
	Uranus      PlanetIndex = 7
	Neptune     PlanetIndex = 8
	PlanetCount int         = 9
)

// String return name of a given PlanetIndex
func (planet PlanetIndex) String() string {
	if planet < Sun || planet > Neptune {
		return "Unknown"
	}
	names := [...]string{"Sun", "Mercury", "Venus", "Earth", "Mars", "Jupiter", "Saturn", "Uranus", "Neptune"}
	return names[planet]
}

// Planet information about a planet
type Planet struct {

	// position of planet on the wall
	position r2.Point

	// size of the planet in mm
	radius float64

	// number of leds
	ledCount int

	// angleOffset radians offset so that axis aligns with physical leds
	angleOffset float64

	// angleDirection -1 if led strip was installed counter clockwise
	angleDirection float64
}

// System all objects that exist in the system, the state of the world
type System struct {

	// The planets
	planets [PlanetCount]Planet

	// All of the drawable items, stored in increasing ZIndex order
	drawables *list.List
}

// DefaultSystem create a solar system with all the data for the planets initialized
func DefaultSystem() *System {
	system := &System{}

	system.planets[Sun] = Planet{r2.Point{X: 9, Y: 11}, 6, 27, 1.39, -1.0}
	system.planets[Mercury] = Planet{r2.Point{X: 7, Y: 25}, 4, 17, 0.78, 1.0}
	system.planets[Venus] = Planet{r2.Point{X: 30, Y: 30}, 4, 17, -0.26, 1.0}
	system.planets[Earth] = Planet{r2.Point{X: 40, Y: 10}, 4, 17, 2.79, 1.0}
	system.planets[Mars] = Planet{r2.Point{X: 60, Y: 19}, 4, 17, 1.04, 1.0}
	system.planets[Jupiter] = Planet{r2.Point{X: 78, Y: 30}, 6, 27, -0.26, -1.0}
	system.planets[Saturn] = Planet{r2.Point{X: 94, Y: 14}, 6, 27, -0.52, -1.0}
	system.planets[Uranus] = Planet{r2.Point{X: 106, Y: 45}, 6, 27, 2.35, -1.0}
	system.planets[Neptune] = Planet{r2.Point{X: 126, Y: 25}, 4, 27, 1.57, 1.0}

	system.drawables = list.New()

	system.drawables.PushFront(NewRotatingLine(Sun, system))
	system.drawables.PushFront(NewRotatingLine(Mercury, system))
	system.drawables.PushFront(NewRotatingLine(Venus, system))
	system.drawables.PushFront(NewRotatingLine(Earth, system))
	system.drawables.PushFront(NewRotatingLine(Mars, system))
	system.drawables.PushFront(NewRotatingLine(Jupiter, system))
	system.drawables.PushFront(NewRotatingLine(Saturn, system))
	system.drawables.PushFront(NewRotatingLine(Uranus, system))
	system.drawables.PushFront(NewRotatingLine(Neptune, system))

	// system.drawables.PushFront(&DrawLine{
	// 	startPosition:   r2.Point{X: 0, Y: 0},
	// 	endPosition:     r2.Point{X: 130, Y: 50},
	// 	traverseTime:    15.0,
	// 	currentPosition: r2.Point{X: 0, Y: 0},
	// 	lineDirection:   r2.Point{X: .707, Y: .707},
	// 	lineWidth:       24.0,
	// 	color:           color.RGBA{R: 255, G: 255, B: 0, A: 128},
	// 	zindex:          2,
	// })

	// system.drawables.PushFront(&DrawLine{
	// 	startPosition:   r2.Point{X: 0, Y: 0},
	// 	endPosition:     r2.Point{X: 130, Y: 0},
	// 	traverseTime:    10.0,
	// 	currentPosition: r2.Point{X: 0, Y: 0},
	// 	lineDirection:   r2.Point{X: 1, Y: 0},
	// 	lineWidth:       18.0,
	// 	color:           color.RGBA{R: 255, G: 0, B: 0, A: 128},
	// 	zindex:          1,
	// })

	// system.drawables.PushFront(&DrawLine{
	// 	startPosition:   r2.Point{X: 0, Y: 0},
	// 	endPosition:     r2.Point{X: 0, Y: 50},
	// 	traverseTime:    12.0,
	// 	currentPosition: r2.Point{X: 0, Y: 0},
	// 	lineDirection:   r2.Point{X: 0, Y: 1},
	// 	lineWidth:       12.0,
	// 	color:           color.RGBA{R: 0, G: 255, B: 0, A: 128},
	// 	zindex:          1,
	// })

	return system
}

// LedCount return total number of LEDs that exist in system
func (solarSystem *System) LedCount() int {
	count := 0
	for _, planet := range solarSystem.planets {
		count += planet.ledCount
	}
	return count
}

// LedPosition return XYPosition of a given Led on the planet
func (solarSystem *System) LedPosition(planetI PlanetIndex, ledIndex int) r2.Point {

	planet := solarSystem.planets[planetI]
	radiansPerLed := (2.0 * math.Pi) / float64(planet.ledCount)

	ledOffset := r2.Point{
		X: math.Cos(float64(ledIndex)*radiansPerLed*planet.angleDirection+planet.angleOffset) * planet.radius * 0.5,
		Y: math.Sin(float64(ledIndex)*radiansPerLed*planet.angleDirection+planet.angleOffset) * planet.radius * 0.5}

	return ledOffset.Add(planet.position)
}

// Animate moves all drawables forward in time
func (solarSystem *System) Animate(dt float64) {

	for curElement := solarSystem.drawables.Front(); curElement != nil; {

		drawable := curElement.Value.(Drawable)

		if !drawable.Animate(dt) {
			nextElement := curElement.Next()
			solarSystem.drawables.Remove(curElement)
			curElement = nextElement
		} else {
			curElement = curElement.Next()
		}
	}
}
