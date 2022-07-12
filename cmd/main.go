package main

import (
	"fmt"

	"github.com/dapper-labs-talent/cc_DaniloPantani_Flow-RangeProcessor/pkg/rangeprocessor"
)

func main() {
	processor := rangeprocessor.New(3, 3, 0)

	minHeight, maxHeight := processor.GetActiveRange()
	fmt.Printf("minHeight: %d | maxHeight: %d\n", minHeight, maxHeight)

	processor.ProcessRange(0, []rangeprocessor.Block{"block_0", "block_1", "block_2"})
	processor.ProcessRange(0, []rangeprocessor.Block{"block_0", "block_1", "block_2"})
	processor.ProcessRange(0, []rangeprocessor.Block{"block_0", "block_1", "block_2"})

	minHeight, maxHeight = processor.GetActiveRange()
	fmt.Printf("minHeight: %d | maxHeight: %d\n", minHeight, maxHeight)

	processor.ProcessRange(3, []rangeprocessor.Block{"block_3", "block_4", "block_5"})
	processor.ProcessRange(3, []rangeprocessor.Block{"block_3", "block_4", "block_5"})
	processor.ProcessRange(3, []rangeprocessor.Block{"block_3", "block_4", "block_5"})

	minHeight, maxHeight = processor.GetActiveRange()
	fmt.Printf("minHeight: %d | maxHeight: %d\n", minHeight, maxHeight)

	processor.ProcessRange(6, []rangeprocessor.Block{"block_6", "block_7", "block_8"})
	processor.ProcessRange(6, []rangeprocessor.Block{"block_6", "block_7", "block_8"})
	processor.ProcessRange(6, []rangeprocessor.Block{"block_6", "block_7", "block_8"})

	minHeight, maxHeight = processor.GetActiveRange()
	fmt.Printf("minHeight: %d | maxHeight: %d\n", minHeight, maxHeight)
}
