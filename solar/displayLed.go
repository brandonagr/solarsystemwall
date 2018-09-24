// +build !windows

package solar

import (
	"log"

	"github.com/jgarff/rpi_ws281x/golang/ws2811"
)

const (
	pin        = 18
	brightness = 255
)

// LedDisplay info needed to render to an image
type LedDisplay struct {
	// total number of Leds for all planets
	totalLedCount int

	// render color for each led
	renderColor []RGBA
}

var testLedDisplay Display = &LedDisplay{}

// NewDisplay return new LedDisplay
func NewDisplay(solarSystem *System) *LedDisplay {
	ledCount := 0
	for planetIndex := 0; planetIndex < PlanetCount; planetIndex++ {
		ledCount += solarSystem.planets[planetIndex].ledCount
	}

	err := ws2811.Init(pin, ledCount, brightness)
	if err != nil {
		log.Fatal(err)
	}

	return &LedDisplay{
		totalLedCount: ledCount,
		renderColor:   make([]RGBA, ledCount),
	}
}

// Dispose cleanup any resources
func (display *LedDisplay) Dispose() {
	ws2811.Clear()
	ws2811.Render()
	ws2811.Fini()
}

// Render the field to an internal structure, that can be read out by the webserver
func (display *LedDisplay) Render(solarSystem *System) {
	firstLedOffset := 0

	// loop through every planet
	for planetIndex := 0; planetIndex < PlanetCount; planetIndex++ {
		planet := solarSystem.planets[planetIndex]

		for led := 0; led < planet.ledCount; led++ {
			ledIndex := led + firstLedOffset
			display.renderColor[ledIndex] = RGBA{R: 0, G: 0, B: 0, A: 255}
		}

		// loop through every drawable object
		for curElement := solarSystem.drawables.Front(); curElement != nil; curElement = curElement.Next() {
			drawable := curElement.Value.(Drawable)

			// bounding circle check to see if this should affect this planet
			if !drawable.Affects(planet.position, planet.radius) {
				continue
			}

			// loop through every led on this planet
			for led := 0; led < planet.ledCount; led++ {
				ledIndex := led + firstLedOffset
				ledPosition := solarSystem.LedPosition(PlanetIndex(planetIndex), led)
				display.renderColor[ledIndex] = drawable.ColorAt(ledPosition, display.renderColor[ledIndex])
			}
		}

		// set each led color
		for led := 0; led < planet.ledCount; led++ {
			ledIndex := led + firstLedOffset
			color := display.renderColor[ledIndex]

			red := (uint32(color.R) * uint32(color.A) ) >> 8
			green := (uint32(color.G) * uint32(color.A) ) >> 8
			blue := (uint32(color.B) * uint32(color.A) ) >> 8

			ws2811.SetLed(ledIndex, uint32(gammaCorrectionLookup[green])<<16 | uint32(gammaCorrectionLookup[red])<<8 | uint32(gammaCorrectionLookup[blue]))
			//ws2811.SetLed(ledIndex, uint32(0x000020))
		}

		firstLedOffset += planet.ledCount
	}

	ws2811.Render()
}

// based on a pull request found at http://forums.adafruit.com/viewtopic.php?f=47&t=26591
// is basically precomputing x = pow(i / 255, 3.0) * 127
var gammaCorrectionLookup = [256]uint8{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 4, 4, 4,
	4, 4, 4, 4, 5, 5, 5, 5, 5, 6, 6, 6, 6, 6, 7, 7,
	7, 7, 7, 8, 8, 8, 8, 9, 9, 9, 9, 10, 10, 10, 10, 11,
	11, 11, 12, 12, 12, 13, 13, 13, 13, 14, 14, 14, 15, 15, 16, 16,
	16, 17, 17, 17, 18, 18, 18, 19, 19, 20, 20, 21, 21, 21, 22, 22,
	23, 23, 24, 24, 24, 25, 25, 26, 26, 27, 27, 28, 28, 29, 29, 30,
	30, 31, 32, 32, 33, 33, 34, 34, 35, 35, 36, 37, 37, 38, 38, 39,
	40, 40, 41, 41, 42, 43, 43, 44, 45, 45, 46, 47, 47, 48, 49, 50,
	50, 51, 52, 52, 53, 54, 55, 55, 56, 57, 58, 58, 59, 60, 61, 62,
	62, 63, 64, 65, 66, 67, 67, 68, 69, 70, 71, 72, 73, 74, 74, 75,
	76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91,
	92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 104, 105, 106, 107, 108,
	109, 110, 111, 113, 114, 115, 116, 117, 118, 120, 121, 122, 123, 125, 126, 127,
}
