package solar

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"time"

	"github.com/golang/geo/r2"
)

// Display that can render the field
type Display interface {
	Render(*System)
}

// XYPositionColor a color at a given position
type XYPositionColor struct {
	position r2.Point
	color    color.RGBA
}

// WebDisplay info needed to render to an image
type WebDisplay struct {
	colorSamples []r2.Point
	width        int
	height       int

	solarSystem *System

	// scale to fit all the planets
	scale r2.Point

	// offset applied to location of each planet
	offset r2.Point

	// image that is rendered to
	image *image.RGBA
}

var testWebDisplay Display = &WebDisplay{}

// NewWebDisplay create a new WebDisplay
func NewWebDisplay(solarSystem *System, width, height int) *WebDisplay {

	// min := XYPosition{10000, 10000}
	// max := XYPosition{0, 0}

	// for _, planet := range solarSystem.planets {
	// 	// use 2 * radius just to give some extra spacing around edges
	// 	planetMin := XYPosition{planet.position.X - 2*planet.radius, planet.position.Y - 2*planet.radius}
	// 	planetMax := XYPosition{planet.position.X + 2*planet.radius, planet.position.Y + 2*planet.radius}

	// 	if planetMin.X < min.X {
	// 		min.X = planetMin.X
	// 	}
	// 	if planetMin.Y < min.Y {
	// 		min.Y = planetMin.Y
	// 	}
	// 	if planetMax.X > max.X {
	// 		max.X = planetMax.X
	// 	}
	// 	if planetMax.Y > max.Y {
	// 		max.Y = planetMax.Y
	// 	}
	// }

	// scaleX := (max.X - min.X) / float64(width)
	// scaleY := (max.Y - min.Y) / float64(height)

	display := &WebDisplay{
		colorSamples: make([]r2.Point, solarSystem.LedCount()),
		width:        width * 5,
		height:       height * 5,
		scale:        r2.Point{X: 5, Y: 5},
		offset:       r2.Point{X: 0, Y: 0},
		solarSystem:  solarSystem,
		image:        image.NewRGBA(image.Rect(0, 0, width*5, height*5)),
	}

	go display.LaunchWebServer()
	return display
}

// Render the field to an internal structure, that can be read out by the webserver
func (display *WebDisplay) Render(solarSystem *System) {
	// given the solarSystem, which contains info on planets and drawable items

	image := image.NewRGBA(image.Rect(0, 0, display.width, display.height))

	// loop through every planet
	for planetIndex := 0; planetIndex < PlanetCount; planetIndex++ {
		planet := display.solarSystem.planets[planetIndex]
		//pos := planet.position

		// loop through every drawable object
		for curElement := solarSystem.drawables.Front(); curElement != nil; curElement = curElement.Next() {
			drawable := curElement.Value.(Drawable)

			// for now want default color to be set for all leds
			// bounding circle check to see if this should affect this planet
			//if !drawable.Affects(pos, planet.radius) {
			//	continue
			//}

			// loop through every led on this planet
			for led := 0; led < planet.ledCount; led++ {
				ledPosition := solarSystem.LedPosition(PlanetIndex(planetIndex), led)

				imageX := int(ledPosition.X * display.scale.X)
				imageY := int(ledPosition.Y * display.scale.Y)

				image.SetRGBA(imageX, imageY, color.RGBA{R: 0, G: 0, B: 255, A: 255})

				curColor := RGBA(image.RGBAAt(imageX, imageY))
				curColor = drawable.ColorAt(ledPosition, RGBA(curColor))

				image.SetRGBA(imageX, imageY, color.RGBA(curColor))
			}
		}
	}

	display.image = image
	// for curElement := field.drawables.Front(); curElement != nil; curElement = curElement.Next() {

	// 	drawable := curElement.Value.(Drawable)

	// 	color = drawable.ColorAt(position, color)
	// }
}

// LaunchWebServer Launches the webserver
func (display *WebDisplay) LaunchWebServer() {

	http.HandleFunc("/", htmlPageHandler)
	http.HandleFunc("/image/", func(w http.ResponseWriter, r *http.Request) { display.imageHandler(w, r) })

	log.Print("Server listening on 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Serve static html page
func htmlPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
<html>
	<head><script type="text/javascript"><!--
		function reloadpic()
        {
			document.images["gameBoard"].src = "image/test.png";
			setTimeout(reloadpic, 400);
        }
        setTimeout(reloadpic, 400)
	--></script></head>
	<body><img id="gameBoard" src="image/test.png" height="910" width="1360"/></body>
</html>`)

}

// imageHandler Return newly generating image
func (display *WebDisplay) imageHandler(w http.ResponseWriter, r *http.Request) {

	startTime := time.Now()

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-control", "max-age=0, must-revalidate, no-store")

	// image := image.NewRGBA(image.Rect(0, 0, display.width, display.height))

	// for planetIndex := 0; planetIndex < PlanetCount; planetIndex++ {
	// 	planet := display.solarSystem.planets[planetIndex]
	// 	pos := planet.position

	// 	pos.X = pos.X * display.scale.X
	// 	pos.Y = pos.Y * display.scale.Y

	// 	radiansPerLed := (2.0 * math.Pi) / float64(planet.ledCount)
	// 	for led := 0; led < planet.ledCount; led++ {
	// 		ledX := pos.X + math.Cos(float64(led)*radiansPerLed)*planet.radius*2
	// 		ledY := pos.Y + math.Sin(float64(led)*radiansPerLed)*planet.radius*2

	// 		image.Set(int(ledX), int(ledY), color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), 255})
	// 	}
	// }

	// for dataIndex := 0; dataIndex < display.width; dataIndex++ {

	// 	for y := 0; y < height; y++ {
	// 		displayedColor := color.RGBA(data[dataIndex])
	// 		displayedColor.A = 255
	// 		image.Set(dataIndex*spacing, y, displayedColor)
	// 	}
	// }

	encoder := &png.Encoder{CompressionLevel: png.NoCompression}
	encoder.Encode(w, display.image)
	log.Print("Generated", r.URL, " in ", time.Since(startTime))
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
