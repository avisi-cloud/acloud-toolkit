package snapshots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSnapshotBatchID(t *testing.T) {
	first, err := NewSnapshotBatchID()
	require.NoError(t, err)
	second, err := NewSnapshotBatchID()
	require.NoError(t, err)

	assert.NoError(t, ValidateSnapshotBatchID(first))
	assert.Regexp(t, `^\d{8}-\d{6}-[a-z0-9]{5}$`, first)
	assert.NotEqual(t, first, second)
}

func TestValidateSnapshotBatchID(t *testing.T) {
	tests := []struct {
		name    string
		batchID string
		wantErr bool
	}{
		{name: "simple", batchID: "before-upgrade"},
		{name: "dots and underscores", batchID: "v1.2_backup"},
		{name: "empty", batchID: "", wantErr: true},
		{name: "spaces", batchID: "before upgrade", wantErr: true},
		{name: "leading dash", batchID: "-backup", wantErr: true},
		{name: "too long", batchID: "a234567890123456789012345678901234567890123456789012345678901234", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSnapshotBatchID(tt.batchID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
