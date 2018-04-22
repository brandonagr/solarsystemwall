package solar

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"net/http"
	"time"
)

// A display that can render the field
type Display interface {
	Render(*System)
}

// XYPositionColor a color at a given position
type XYPositionColor struct {
	position XYPosition
	color    RGBA
}

// Web display
type WebDisplay struct {
	colorSamples []XYPositionColor
	width        int
	height       int

	// scale to fit all the planets
	scale XYPosition

	// offset applied to location of each planet
	offset XYPosition
}

var testWebDisplay Display = &WebDisplay{}

// Create a new WebDisplay
func NewWebDisplay(solarSystem *System, width, height int) *WebDisplay {

	min := XYPosition{10000, 10000}
	max := XYPosition{0, 0}

	for _, planet := range solarSystem.planets {
		// use 2 * radius just to give some extra spacing around edges
		planetMin := XYPosition{planet.position.X - 2*planet.radius, planet.position.Y - 2*planet.radius}
		planetMax := XYPosition{planet.position.X + 2*planet.radius, planet.position.Y + 2*planet.radius}

		if planetMin.X < min.X {
			min.X = planetMin.X
		}
		if planetMin.Y < min.Y {
			min.Y = planetMin.Y
		}
		if planetMax.X > max.X {
			max.X = planetMax.X
		}
		if planetMax.Y > max.Y {
			max.Y = planetMax.Y
		}
	}

	scaleX := (max.X - min.X) / float64(width)
	scaleY := (max.Y - min.Y) / float64(height)

	display := &WebDisplay{
		colorSamples: make([]XYPositionColor, solarSystem.LedCount()),
		width:        width,
		height:       height,
		scale:        XYPosition{scaleX, scaleY},
		offset:       XYPosition{-min.X, -min.Y},
	}

	go display.LaunchWebServer()
	return display
}

// Render the field to an internal structure, that can be read out by the webserver
func (display *WebDisplay) Render(solarSystem *System) {

}

// LaunchWebServer Launches the webserver
func (display *WebDisplay) LaunchWebServer() {

	http.HandleFunc("/", htmlPageHandler)
	http.HandleFunc("/image/", func(w http.ResponseWriter, r *http.Request) { display.imageHandler(w, r) })

	log.Print("Server listening on 8080")
	http.ListenAndServe(":8080", nil)
}

// Serve static html page
func htmlPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
<html>
	<head><script type="text/javascript"><!--
		function reloadpic()
        {
			document.images["gameBoard"].src = "image/test.png";
			setTimeout(reloadpic, 100);
        }
        setTimeout(reloadpic, 100)
	--></script></head>
	<body><img id="gameBoard" src="image/test.png" height="42" width="1024"/></body>
</html>`)

}

// imageHandler Return newly generating image
func (display *WebDisplay) imageHandler(w http.ResponseWriter, r *http.Request) {

	startTime := time.Now()

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-control", "max-age=0, must-revalidate, no-store")

	image := image.NewRGBA(image.Rect(0, 0, display.width, display.height))

	// for dataIndex := 0; dataIndex < width; dataIndex++ {

	// 	for y := 0; y < height; y++ {
	// 		displayedColor := color.RGBA(data[dataIndex])
	// 		displayedColor.A = 255
	// 		image.Set(dataIndex*spacing, y, displayedColor)
	// 	}
	// }

	png.Encode(w, image)
	log.Print("Generated", r.URL, " in", time.Since(startTime))
}

// based on a pull request found at http://forums.adafruit.com/viewtopic.php?f=47&t=26591
// is basically precomputing x = pow(i / 255, 3.0) * 127
var gammaCorrectionLookup [256]uint8 = [256]uint8{
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
