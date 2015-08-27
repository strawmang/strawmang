package chat

import (
	"testing"
)

func TestColors(t *testing.T) {
	for i := 0; i < 1200; i++ {
		t.Log("Generating color", i)
		color := PopColor()
		t.Logf("Generated color: %v", color.Hex())
	}
}
