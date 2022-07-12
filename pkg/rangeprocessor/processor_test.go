package rangeprocessor

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRangeResponseProcessor_GetActiveRange(t *testing.T) {
	type fields struct {
		activeRange    uint64
		lastHeight     uint64
		blockResponses int64
	}
	tests := []struct {
		name          string
		fields        fields
		wantMinHeight uint64
		wantMaxHeight uint64
	}{
		{
			name: "test valid parameters",
			fields: fields{
				activeRange:    33,
				lastHeight:     10,
				blockResponses: 7,
			},
			wantMinHeight: 7,
			wantMaxHeight: 39,
		},
		{
			name: "test empty range",
			fields: fields{
				activeRange:    0,
				lastHeight:     10,
				blockResponses: 7,
			},
			wantMinHeight: 7,
			wantMaxHeight: 7,
		},
		{
			name: "test empty height",
			fields: fields{
				activeRange:    33,
				lastHeight:     0,
				blockResponses: 7,
			},
			wantMinHeight: 7,
			wantMaxHeight: 39,
		},
		{
			name: "test empty block responses",
			fields: fields{
				activeRange:    33,
				lastHeight:     10,
				blockResponses: 0,
			},
			wantMinHeight: 0,
			wantMaxHeight: 32,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New(tt.fields.activeRange, tt.fields.lastHeight, tt.fields.blockResponses)
			gotMinHeight, gotMaxHeight := r.GetActiveRange()
			require.Equal(t, tt.wantMinHeight, gotMinHeight)
			require.Equal(t, tt.wantMaxHeight, gotMaxHeight)
		})
	}
}

func TestRangeResponseProcessor_ProcessRange(t *testing.T) {
	t.Run("test new node sync", func(t *testing.T) {
		r := New(3, 4, 0)

		gotMinHeight, gotMaxHeight := r.GetActiveRange()
		require.Equal(t, uint64(0), gotMinHeight)
		require.Equal(t, uint64(2), gotMaxHeight)

		r.ProcessRange(0, []Block{"block_0", "block_1", "block_2"})
		r.ProcessRange(0, []Block{"block_0", "block_1", "block_2"})
		r.ProcessRange(0, []Block{"block_0", "block_1", "block_2"})
		r.ProcessRange(0, []Block{"block_0", "block_1", "block_2"})

		gotMinHeight, gotMaxHeight = r.GetActiveRange()
		require.Equal(t, uint64(3), gotMinHeight)
		require.Equal(t, uint64(5), gotMaxHeight)

		r.ProcessRange(3, []Block{"block_3", "block_4", "block_5"})
		r.ProcessRange(3, []Block{"block_3", "block_4", "block_5"})
		r.ProcessRange(3, []Block{"block_3", "block_4", "block_5"})
		r.ProcessRange(3, []Block{"block_3", "block_4", "block_5"})

		gotMinHeight, gotMaxHeight = r.GetActiveRange()
		require.Equal(t, uint64(6), gotMinHeight)
		require.Equal(t, uint64(8), gotMaxHeight)

		r.ProcessRange(6, []Block{"block_6", "block_7", "block_8"})
		r.ProcessRange(6, []Block{"block_6", "block_7", "block_8"})
		r.ProcessRange(6, []Block{"block_6", "block_7", "block_8"})
		r.ProcessRange(6, []Block{"block_6", "block_7", "block_8"})

		gotMinHeight, gotMaxHeight = r.GetActiveRange()
		require.Equal(t, uint64(9), gotMinHeight)
		require.Equal(t, uint64(11), gotMaxHeight)
	})

	t.Run("test out of range height", func(t *testing.T) {
		r := New(3, 2, 0)

		gotMinHeight, gotMaxHeight := r.GetActiveRange()
		require.Equal(t, uint64(0), gotMinHeight)
		require.Equal(t, uint64(2), gotMaxHeight)

		r.ProcessRange(0, []Block{"block_0", "block_1", "block_2"})
		r.ProcessRange(0, []Block{"block_0", "block_1", "block_2"})

		gotMinHeight, gotMaxHeight = r.GetActiveRange()
		require.Equal(t, uint64(3), gotMinHeight)
		require.Equal(t, uint64(5), gotMaxHeight)

		r.ProcessRange(0, []Block{"block_3", "block_4", "block_5"})
		r.ProcessRange(0, []Block{"block_3", "block_4", "block_5"})

		gotMinHeight, gotMaxHeight = r.GetActiveRange()
		require.Equal(t, uint64(3), gotMinHeight)
		require.Equal(t, uint64(5), gotMaxHeight)

		r.ProcessRange(6, []Block{"block_3", "block_4", "block_5"})
		r.ProcessRange(6, []Block{"block_3", "block_4", "block_5"})

		gotMinHeight, gotMaxHeight = r.GetActiveRange()
		require.Equal(t, uint64(3), gotMinHeight)
		require.Equal(t, uint64(5), gotMaxHeight)

		r.ProcessRange(5, []Block{"block_5", "block_6", "block_7"})
		r.ProcessRange(5, []Block{"block_5", "block_6", "block_7"})

		gotMinHeight, gotMaxHeight = r.GetActiveRange()
		require.Equal(t, uint64(6), gotMinHeight)
		require.Equal(t, uint64(8), gotMaxHeight)
	})

	t.Run("test with already existing blocks", func(t *testing.T) {
		r := New(3, 3, 3)
		r.blockCounter.Store(0, "block_0")
		r.blockCounter.Store(1, "block_1")
		r.blockCounter.Store(2, "block_2")

		gotMinHeight, gotMaxHeight := r.GetActiveRange()
		require.Equal(t, uint64(3), gotMinHeight)
		require.Equal(t, uint64(5), gotMaxHeight)

		r.ProcessRange(3, []Block{"block_3", "block_4"})
		r.ProcessRange(3, []Block{"block_3", "block_4"})
		r.ProcessRange(3, []Block{"block_3", "block_4"})

		gotMinHeight, gotMaxHeight = r.GetActiveRange()
		require.Equal(t, uint64(5), gotMinHeight)
		require.Equal(t, uint64(7), gotMaxHeight)
	})

	t.Run("test block number greater than the range", func(t *testing.T) {
		r := New(2, 4, 0)

		gotMinHeight, gotMaxHeight := r.GetActiveRange()
		require.Equal(t, uint64(0), gotMinHeight)
		require.Equal(t, uint64(1), gotMaxHeight)

		r.ProcessRange(0, []Block{"block_0", "block_1", "block_2", "block_3", "block_4"})
		r.ProcessRange(0, []Block{"block_0", "block_1", "block_2", "block_3", "block_4"})
		r.ProcessRange(0, []Block{"block_0", "block_1", "block_2", "block_3", "block_4"})
		r.ProcessRange(0, []Block{"block_0", "block_1", "block_2", "block_3", "block_4"})

		gotMinHeight, gotMaxHeight = r.GetActiveRange()
		require.Equal(t, uint64(2), gotMinHeight)
		require.Equal(t, uint64(3), gotMaxHeight)

		r.ProcessRange(2, []Block{"block_5", "block_6", "block_7", "block_8", "block_9"})
		r.ProcessRange(2, []Block{"block_5", "block_6", "block_7", "block_8", "block_9"})
		r.ProcessRange(2, []Block{"block_5", "block_6", "block_7", "block_8", "block_9"})
		r.ProcessRange(2, []Block{"block_5", "block_6", "block_7", "block_8", "block_9"})

		gotMinHeight, gotMaxHeight = r.GetActiveRange()
		require.Equal(t, uint64(4), gotMinHeight)
		require.Equal(t, uint64(5), gotMaxHeight)

		r.ProcessRange(4, []Block{"block_10", "block_11", "block_12", "block_13", "block_14"})
		r.ProcessRange(4, []Block{"block_10", "block_11", "block_12", "block_13", "block_14"})
		r.ProcessRange(4, []Block{"block_10", "block_11", "block_12", "block_13", "block_14"})
		r.ProcessRange(4, []Block{"block_10", "block_11", "block_12", "block_13", "block_14"})

		gotMinHeight, gotMaxHeight = r.GetActiveRange()
		require.Equal(t, uint64(6), gotMinHeight)
		require.Equal(t, uint64(7), gotMaxHeight)
	})

	t.Run("test nodes with less blocks", func(t *testing.T) {
		r := New(3, 3, 0)

		gotMinHeight, gotMaxHeight := r.GetActiveRange()
		require.Equal(t, uint64(0), gotMinHeight)
		require.Equal(t, uint64(2), gotMaxHeight)

		r.ProcessRange(0, []Block{"block_0", "block_1", "block_2"})
		r.ProcessRange(0, []Block{"block_0", "block_1", "block_2"})
		r.ProcessRange(0, []Block{"block_0", "block_1", "block_2"})

		gotMinHeight, gotMaxHeight = r.GetActiveRange()
		require.Equal(t, uint64(3), gotMinHeight)
		require.Equal(t, uint64(5), gotMaxHeight)

		r.ProcessRange(3, []Block{"block_3"})
		r.ProcessRange(3, []Block{"block_3"})
		r.ProcessRange(3, []Block{"block_3"})

		gotMinHeight, gotMaxHeight = r.GetActiveRange()
		require.Equal(t, uint64(4), gotMinHeight)
		require.Equal(t, uint64(6), gotMaxHeight)

		r.ProcessRange(6, []Block{"block_4", "block_5"})
		r.ProcessRange(6, []Block{"block_4", "block_5"})
		r.ProcessRange(6, []Block{"block_4", "block_5"})

		gotMinHeight, gotMaxHeight = r.GetActiveRange()
		require.Equal(t, uint64(7), gotMinHeight)
		require.Equal(t, uint64(9), gotMaxHeight)
	})

	t.Run("test nodes with messy blocks", func(t *testing.T) {
		r := New(3, 4, 0)

		gotMinHeight, gotMaxHeight := r.GetActiveRange()
		require.Equal(t, uint64(0), gotMinHeight)
		require.Equal(t, uint64(2), gotMaxHeight)

		r.ProcessRange(0, []Block{"block_0", "block_1", "block_2"})
		r.ProcessRange(0, []Block{"block_0", "block_1", "block_2", "block_3", "block_4", "block_5"})
		r.ProcessRange(0, []Block{"block_0", "block_1", "block_2", "block_3", "block_4", "block_5", "block_6"})
		r.ProcessRange(0, []Block{"block_0", "block_1", "block_2"})
		r.ProcessRange(0, []Block{"block_0", "block_2"})

		gotMinHeight, gotMaxHeight = r.GetActiveRange()
		require.Equal(t, uint64(3), gotMinHeight)
		require.Equal(t, uint64(5), gotMaxHeight)

		r.ProcessRange(5, []Block{"block_5", "block_6", "block_7", "block_8"})
		r.ProcessRange(5, []Block{"block_5"})
		r.ProcessRange(5, []Block{"block_5", "block_6", "block_7", "block_8", "block_9", "block_10"})
		r.ProcessRange(5, []Block{"block_5", "block_6", "block_7", "block_8", "block_9"})
		r.ProcessRange(5, []Block{"block_5", "block_6", "block_7"})
		r.ProcessRange(5, []Block{"block_5", "block_6"})

		gotMinHeight, gotMaxHeight = r.GetActiveRange()
		require.Equal(t, uint64(6), gotMinHeight)
		require.Equal(t, uint64(8), gotMaxHeight)

		r.ProcessRange(8, []Block{"block_8", "block_9", "block_10", "block_11", "block_12"})
		r.ProcessRange(8, []Block{"block_8"})
		r.ProcessRange(8, []Block{"block_8", "block_9", "block_10"})
		r.ProcessRange(8, []Block{"block_8", "block_9", "block_10", "block_11"})

		gotMinHeight, gotMaxHeight = r.GetActiveRange()
		require.Equal(t, uint64(9), gotMinHeight)
		require.Equal(t, uint64(11), gotMaxHeight)
	})
}

func TestBlock_Equals(t *testing.T) {
	tests := []struct {
		name string
		b    Block
		cmp  Block
		want bool
	}{
		{
			name: "equal",
			b:    "block_0",
			cmp:  "block_0",
			want: true,
		},
		{
			name: "not equal",
			b:    "block_0",
			cmp:  "block_1",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.b.Equals(tt.cmp)
			require.Equal(t, tt.want, got)
		})
	}
}

func BenchmarkPrimeNumbers(b *testing.B) {
	inputs := []int{0, 100, 1000, 10000, 50000, 100000, 500000}

	type fields struct {
		activeRange    uint64
		lastHeight     uint64
		blockResponses int64
	}
	tests := []struct {
		name        string
		fields      fields
		blockInputs []int
	}{
		{
			name: "test valid parameters",
			fields: fields{
				activeRange:    33,
				lastHeight:     10,
				blockResponses: 7,
			},
			blockInputs: inputs,
		},
		{
			name: "test empty range",
			fields: fields{
				activeRange:    0,
				lastHeight:     10,
				blockResponses: 7,
			},
			blockInputs: inputs,
		},
		{
			name: "test empty height",
			fields: fields{
				activeRange:    33,
				lastHeight:     0,
				blockResponses: 7,
			},
			blockInputs: inputs,
		},
		{
			name: "test empty block responses",
			fields: fields{
				activeRange:    33,
				lastHeight:     10,
				blockResponses: 0,
			},
			blockInputs: inputs,
		},
	}
	for _, tt := range tests {
		for _, inputs := range tt.blockInputs {
			b.Run(fmt.Sprintf("%s_%d", tt.name, inputs), func(b *testing.B) {
				r := New(tt.fields.activeRange, tt.fields.lastHeight, tt.fields.blockResponses)

				blocks := make([]Block, inputs)
				for i := 0; i < inputs; i++ {
					blocks[i] = Block("block_" + strconv.Itoa(i))
				}

				for i := 0; i < b.N; i++ {
					r.ProcessRange(uint64(i), blocks)
				}
			})
		}
	}
}
