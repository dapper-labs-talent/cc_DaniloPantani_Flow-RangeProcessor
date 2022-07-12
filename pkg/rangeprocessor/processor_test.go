package rangeprocessor

import (
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
			wantMinHeight: 10,
			wantMaxHeight: 42,
		},
		{
			name: "test empty range",
			fields: fields{
				activeRange:    0,
				lastHeight:     10,
				blockResponses: 7,
			},
			wantMinHeight: 10,
			wantMaxHeight: 10,
		},
		{
			name: "test empty height",
			fields: fields{
				activeRange:    33,
				lastHeight:     0,
				blockResponses: 7,
			},
			wantMinHeight: 0,
			wantMaxHeight: 32,
		},
		{
			name: "test empty block responses",
			fields: fields{
				activeRange:    33,
				lastHeight:     10,
				blockResponses: 0,
			},
			wantMinHeight: 10,
			wantMaxHeight: 42,
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
