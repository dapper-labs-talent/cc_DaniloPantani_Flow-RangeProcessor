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

// GetActiveRange returns the current active range.
func (r *RangeResponseProcessor) GetActiveRange() (minHeight uint64, maxHeight uint64) {
	minHeight = uint64(r.nextHeight)
	maxHeight = minHeight + r.activeRange - 1
	return
}

// ProcessRange process each block range response it receives and update the processor
func (r *RangeResponseProcessor) ProcessRange(startHeight uint64, blocks []Block) {
	minHeight, maxHeight := r.GetActiveRange()
	var wg sync.WaitGroup
	for i, block := range blocks {
		height := startHeight + uint64(i)

		// skip invalid block range
		if height < minHeight || height > maxHeight {
			continue
		}
		// block already exist
		if _, ok := r.blocks.Load(height); ok {
			continue
		}

		wg.Add(1)
		block := block
		go func() {
			defer wg.Done()

			counter := uint64(0)
			c, ok := r.blockCounter.Load(height)
			if ok {
				counter, _ = c.(uint64)
			}

			counter++
			r.blockCounter.Store(height, counter)

			// if reach the minimal responses, add the new block
			if counter >= r.minResponses {
				r.blocks.Store(height, block)
				if int64(height) >= r.nextHeight-1 {
					r.nextHeight = int64(height) + 1
				}
			}
		}()
	}
	wg.Wait()
}

// Equals returns true if two blocks are equal
func (b Block) Equals(cmp Block) bool {
	return b == cmp
}
