// +build windows

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

// NewDisplay create a new WebDisplay
func NewDisplay(solarSystem *System) *WebDisplay {

	width := 136
	height := 91

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

// Dispose cleanup any resources
func (display *WebDisplay) Dispose() {
}

// Render the field to an internal structure, that can be read out by the webserver
func (display *WebDisplay) Render(solarSystem *System) {
	// given the solarSystem, which contains info on planets and drawable items

	image := image.NewRGBA(image.Rect(0, 0, display.width, display.height))

	// loop through every planet
	for planetIndex := 0; planetIndex < PlanetCount; planetIndex++ {
		planet := display.solarSystem.planets[planetIndex]
		//pos := planet.position

		// initialize default color for each led position
		for led := 0; led < planet.ledCount; led++ {
			ledPosition := solarSystem.LedPosition(PlanetIndex(planetIndex), led)

			imageX := int(ledPosition.X * display.scale.X)
			imageY := int(ledPosition.Y * display.scale.Y)

			image.SetRGBA(imageX, imageY, color.RGBA{R: 0, G: 0, B: 0, A: 255})
		}

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

				curColor := RGBA(image.RGBAAt(imageX, imageY))
				curColor = drawable.ColorAt(ledPosition, RGBA(curColor))

				curColor.A = 255 // for rendering to image dont want to blend to nothing
				image.SetRGBA(imageX, imageY, color.RGBA(curColor))
			}
		}
	}

	display.image = image
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
	<body bgcolor="#888888"><img id="gameBoard" src="image/test.png" height="910" width="1360"/></body>
</html>`)

}

// imageHandler Return newly generating image
func (display *WebDisplay) imageHandler(w http.ResponseWriter, r *http.Request) {

	startTime := time.Now()

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-control", "max-age=0, must-revalidate, no-store")

	encoder := &png.Encoder{CompressionLevel: png.NoCompression}
	encoder.Encode(w, display.image)
	log.Print("Generated", r.URL, " in ", time.Since(startTime))
}
