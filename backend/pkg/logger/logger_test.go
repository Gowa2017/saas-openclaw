package logger

import (
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name  string
		level string
	}{
		{"debug level", "debug"},
		{"info level", "info"},
		{"warn level", "warn"},
		{"error level", "error"},
		{"invalid level defaults to info", "invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := New(tt.level)
			if err != nil {
				t.Fatalf("New(%s) error = %v", tt.level, err)
			}
			if logger == nil {
				t.Errorf("New(%s) returned nil logger", tt.level)
			}
		})
	}
}
