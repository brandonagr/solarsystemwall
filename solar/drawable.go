package solar

import (
	"image/color"
)

// RGBA color used in game
type RGBA color.RGBA

// ZIndex defines the z order of different Drawable
type ZIndex int

// Drawable methods required to draw something
type Drawable interface {

	// Computes the color at the given position with the given baseColor
	// If this drawable does not influence the given XYPosition then it simply returns baseColor
	ColorAt(position XYPosition, baseColor RGBA) RGBA

	// The ZIndex of this Drawable thing
	ZIndex() ZIndex

	// Move this Drawable forward in time by dt, returns if this should be kept alive or deleted
	Animate(dt float64) (keepAlive bool)
}

// BlendWith helper function to blend two colors together
func (foreground RGBA) BlendWith(background RGBA) (color RGBA) {

	fr, fg, fb, fa := uint(foreground.R), uint(foreground.G), uint(foreground.B), uint(foreground.A)
	br, bg, bb, ba := uint(background.R), uint(background.G), uint(background.B), uint(255) // want background to be fully colored

	opacity := fa
	backgroundOpacity := (ba * (255 - fa)) >> 8

	newColor := RGBA{
		uint8((fr*opacity)>>8 + (br*backgroundOpacity)>>8),
		uint8((fg*opacity)>>8 + (bg*backgroundOpacity)>>8),
		uint8((fb*opacity)>>8 + (bb*backgroundOpacity)>>8),
		uint8(opacity),
	}

	return newColor
}
