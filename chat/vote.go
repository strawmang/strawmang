package chat

type voteOption uint16

func (v voteOption) Has(opt voteOption) bool {
	if (v & opt) != 0 {
		return true
	}
	return false
}

const (
	VOTE_OPTION_A voteOption = 1 << iota
	VOTE_OPTION_B
	// Add more as needed
)
