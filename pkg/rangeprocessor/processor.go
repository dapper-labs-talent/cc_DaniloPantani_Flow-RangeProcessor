package rangeprocessor

import "sync"

const (
	// minRange represents the size of the active range
	minRange = 1
	// minResponses represents the responses containing that height
	minResponses = 1
)

type (
	// Block represents the chain Block
	Block string

	// RangeResponseProcessor represents the chain range processor
	RangeResponseProcessor struct {
		activeRange  uint64
		minResponses uint64
		nextHeight   int64
		blockCounter sync.Map
		blocks       sync.Map
	}
)

// New instantiate a new RangeResponseProcessor passing the active height range, the
// block height for the processor should start, and the minimal response per block height
func New(activeRange, blockResponses uint64, startHeight int64) RangeResponseProcessor {
	if activeRange < minRange {
		activeRange = minRange
	}
	if blockResponses < minResponses {
		blockResponses = minResponses
	}
	return RangeResponseProcessor{
		activeRange:  activeRange,
		nextHeight:   startHeight,
		minResponses: blockResponses,
	}
}

// Equals returns true if two blocks are equal
func (b Block) Equals(cmp Block) bool {
	return b == cmp
}
