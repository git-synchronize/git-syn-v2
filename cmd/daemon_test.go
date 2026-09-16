/*
Copyright © 2024-2026 Lucas Ramage <lucas.ramage@infinite-omicron.com>
*/
package cmd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveInterval(t *testing.T) {
	tests := []struct {
		name           string
		flagInterval   time.Duration
		configInterval time.Duration
		once           bool
		want           time.Duration
		wantErr        bool
	}{
		{
			name:           "flag takes priority over config",
			flagInterval:   time.Minute,
			configInterval: 5 * time.Minute,
			want:           time.Minute,
		},
		{
			name:           "falls back to config when flag unset",
			flagInterval:   0,
			configInterval: 5 * time.Minute,
			want:           5 * time.Minute,
		},
		{
			name:           "zero interval is fine when once is set",
			flagInterval:   0,
			configInterval: 0,
			once:           true,
			want:           0,
		},
		{
			name:           "zero interval without once is an error",
			flagInterval:   0,
			configInterval: 0,
			once:           false,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveInterval(tt.flagInterval, tt.configInterval, tt.once)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
