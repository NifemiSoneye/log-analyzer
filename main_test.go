package main

import "testing"

func TestParseLine(t *testing.T) {
	tests := []struct {
		name         string
		line         string
		wantLevel    string
		wantEndpoint string
		wantOk       bool
	}{
		{
			name:         "well-formed line",
			line:         "10:15:02 ERROR /api/orders 500",
			wantLevel:    "ERROR",
			wantEndpoint: "/api/orders",
			wantOk:       true,
		},
		{
			name:         "truncated line",
			line:         "10:15:04",
			wantLevel:    "",
			wantEndpoint: "",
			wantOk:       false,
		},
		{
			name:         "empty line",
			line:         "",
			wantLevel:    "",
			wantEndpoint: "",
			wantOk:       false,
		},
	}

	for _, tt := range tests {
		entry, ok := parseLine(tt.line)

		if ok != tt.wantOk {
			t.Errorf("%s: expected ok=%v, got %v", tt.name, tt.wantOk, ok)
		}
		if entry.Level != tt.wantLevel {
			t.Errorf("%s: expected Level=%s, got %s", tt.name, tt.wantLevel, entry.Level)
		}
		if entry.Endpoint != tt.wantEndpoint {
			t.Errorf("%s: expected Endpoint=%s, got %s", tt.name, tt.wantEndpoint, entry.Endpoint)
		}
	}
}

func BenchmarkParseLine(b *testing.B) {
    line := "10:15:02 ERROR /api/orders 500"
    for i := 0; i < b.N; i++ {
        parseLine(line)
    }
}