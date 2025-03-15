package resolution

import (
	"testing"
)

func TestGetStandardResolution(t *testing.T) {
	tests := []struct {
		name     string
		width    int
		height   int
		expected string
	}{
		{"Exact 8K", 7680, 4320, "8K"},
		{"Exact 4K", 3840, 2160, "4K"},
		{"Exact 2K", 2048, 1080, "2K"},
		{"Exact 1440p", 2560, 1440, "1440p"},
		{"Exact 1080p", 1920, 1080, "1080p"},
		{"Exact 720p", 1280, 720, "720p"},
		{"Exact 480p", 640, 480, "480p"},
		{"Exact 360p", 640, 360, "360p"},
		{"Near 1080p", 1920, 1082, "1080p"},     // Within tolerance
		{"Near 720p", 1280, 715, "720p"},        // Within tolerance
		{"Near 1080p Low", 1920, 1075, "1080p"}, // Within tolerance
		{"Custom Resolution", 1600, 900, "900p"},
		{"Very Low Resolution", 320, 240, "240p"},
		{"Wide 1080p", 2560, 1080, "1080p"},
		{"Custom High Resolution", 3000, 2000, "2000p"}, // Non-standard resolution
		{"Zero dimensions", 0, 0, "0p"},
		{"Ultrawide 1440p", 3440, 1440, "1440p"},
		{"Near 4K", 3840, 2150, "4K"},
		{"Near 4K High", 3840, 2170, "4K"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetStandardResolution(tt.width, tt.height)
			if got != tt.expected {
				t.Errorf("GetStandardResolution(%d, %d) = %v, want %v",
					tt.width, tt.height, got, tt.expected)
			}
		})
	}
}

func TestAbs(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"Positive number", 5, 5},
		{"Negative number", -5, 5},
		{"Zero", 0, 0},
		{"Large positive", 1000, 1000},
		{"Large negative", -1000, 1000},
		{"Max int", int(^uint(0) >> 1), int(^uint(0) >> 1)},
		{"Min int + 1", -(int(^uint(0)>>1) - 1), int(^uint(0)>>1) - 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := abs(tt.input)
			if got != tt.expected {
				t.Errorf("abs(%d) = %v, want %v",
					tt.input, got, tt.expected)
			}
		})
	}
}
