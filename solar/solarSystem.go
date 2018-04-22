package solar

import (
	"container/list"
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
	position XYPosition

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

// LedCount return total number of LEDs that exist in system
func (solarSystem *System) LedCount() int {
	count := 0
	for _, planet := range solarSystem.planets {
		count += planet.ledCount
	}
	return count
}
