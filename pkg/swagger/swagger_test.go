package swagger

import "testing"

func TestParseFromCollection(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "default test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ParseFromCollection()
		})
	}
}
