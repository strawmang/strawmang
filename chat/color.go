package chat

import (
	"github.com/lucasb-eyer/go-colorful"
	"log"
)

var colors = []colorful.Color{}

// PopColor returns a new color from our color stack.  If we run
// out of colors and fail to allocate a new batch we will return
// an error
func PopColor() colorful.Color {
	if len(colors) == 0 {
		AllocateColors()
	}
	var color colorful.Color
	color, colors = colors[len(colors)-1], colors[:len(colors)-1]
	return color
}

// AllocateColors returns a fresh batch of colors
func AllocateColors() {
	log.Printf("Allocating new generateColor batch")
	colors = colorful.FastHappyPalette(1000)
}

// TODO: This
func generateColor() string {
	// Too hacky looking?
	return PopColor().Hex()[1:]
}
