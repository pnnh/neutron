package convert

import (
	"testing"
)

func TestConvertString(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    string
		wantErr bool
	}{
		{"nil value", nil, "", true},
		{"string", "hello", "hello", false},
		{"int", 42, "42", false},
		{"int8", int8(8), "8", false},
		{"int16", int16(16), "16", false},
		{"int32", int32(32), "32", false},
		{"int64", int64(64), "64", false},
		{"uint", uint(7), "7", false},
		{"uint8", uint8(8), "8", false},
		{"uint16", uint16(16), "16", false},
		{"uint32", uint32(32), "32", false},
		{"uint64", uint64(64), "64", false},
		{"[]byte", []byte("bytes"), "bytes", false},
		{"float32", 3.14, "3.140000", false},
		{"float64", 3.14159, "3.141590", false},
		{"bool true", true, "true", false},
		{"bool false", false, "false", false},
		{"unsupported type", struct{}{}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertString() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ConvertString() = %v, want %v", got, tt.want)
			}
		})
	}
}
