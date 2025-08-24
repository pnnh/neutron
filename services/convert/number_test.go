package convert

import (
	"testing"
)

func TestToInt64(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    int64
		wantErr bool
	}{
		{"nil value", nil, 0, true},
		{"int", int(42), 42, false},
		{"int8", int8(8), 8, false},
		{"uint8", uint8(8), 8, false},
		{"int16", int16(16), 16, false},
		{"uint16", uint16(16), 16, false},
		{"int32", int32(32), 32, false},
		{"uint32", uint32(32), 32, false},
		{"int64", int64(64), 64, false},
		{"uint64 normal", uint64(64), 64, false},
		{"uint64 too large", uint64(^uint64(0)), 0, true},
		{"float32 normal", float32(12.0), 12, false},
		{"float32 negative", float32(-1.0), 0, true},
		{"float32 too large", float32(1e20), 0, true},
		{"float64 normal", float64(123.0), 123, false},
		{"float64 negative", float64(-1.0), 0, true},
		{"float64 too large", float64(1e20), 0, true},
		{"string valid", "12345", 12345, false},
		{"string invalid", "abc", 0, true},
		{"bool true", true, 1, false},
		{"bool false", false, 0, false},
		{"unsupported type", []int{1, 2, 3}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToInt64(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToInt64() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ToInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}
