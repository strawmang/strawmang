package chat

import (
	"github.com/lucasb-eyer/go-colorful"
	"log"
	"math/rand"
	"time"
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
	rand.Seed(time.Now().UTC().UnixNano())
	log.Printf("Allocating new color batch")
	for i := 0; i < 100; i++ {
		colors = append(colors, colorful.FastHappyColor())
	}
}

// TODO: This
func generateColor() string {
	// Too hacky looking?
	return PopColor().Hex()[1:]
}
