package solar

import (
	"container/list"

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

	// verticalLedOffset the led pointing up
	verticalLedOffset int
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

	system.planets[Sun] = Planet{r2.Point{X: 9, Y: 11}, 6, 27, 0}
	system.planets[Mercury] = Planet{r2.Point{X: 7, Y: 25}, 4, 17, 0}
	system.planets[Venus] = Planet{r2.Point{X: 30, Y: 30}, 4, 17, 0}
	system.planets[Earth] = Planet{r2.Point{X: 40, Y: 10}, 4, 17, 0}
	system.planets[Mars] = Planet{r2.Point{X: 60, Y: 19}, 4, 17, 0}
	system.planets[Jupiter] = Planet{r2.Point{X: 78, Y: 30}, 6, 27, 0}
	system.planets[Saturn] = Planet{r2.Point{X: 94, Y: 14}, 6, 27, 0}
	system.planets[Uranus] = Planet{r2.Point{X: 106, Y: 45}, 6, 27, 0}
	system.planets[Neptune] = Planet{r2.Point{X: 126, Y: 25}, 4, 17, 0}

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
func (solarSystem *System) LedPosition(planet PlanetIndex, ledIndex int) r2.Point {
	return solarSystem.planets[planet].position
}
