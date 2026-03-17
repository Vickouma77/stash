package models

import (
	"testing"

	"stash.io/internal/assert"
)

func TestUserModelExists(t *testing.T) {
	if testing.Short() {
		t.Skip("model: skipping integration test")
	}
	
	tests := []struct {
		name string
		userID int
		want bool
	}{
		{
			name: "Valid is",
			userID: 1,
			want: true,
		},
		{
			name: "Zero id",
			userID: 0,
			want: false,
		},
		{
			name: "Non-existen id",
			userID: 2,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)

			m := UserModel{db}

			exists, err := m.Exists(tt.userID)

			assert.Equal(t, exists, tt.want)
			assert.NilError(t, err)
		})
	}
}