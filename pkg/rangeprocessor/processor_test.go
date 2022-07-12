package rangeprocessor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
