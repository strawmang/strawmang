package chat

import (
	"testing"
)

func TestMaskHas(t *testing.T) {
	mask := voteOption(VOTE_OPTION_A | VOTE_OPTION_B)
	if !mask.Has(VOTE_OPTION_A) {
		t.Fail()
	}

	if !mask.Has(VOTE_OPTION_B) {
		t.Fail()
	}

	if mask.Has(1 << 3) {
		t.Fail()
	}
}
