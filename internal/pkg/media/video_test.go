// Package media provides video file analysis functionality.
// This file contains tests for the video analysis capabilities.
package media

import (
	"testing"
)

// TestIsVideoFile verifies the video file detection functionality.
// These tests ensure VidKit correctly identifies video files based on their extensions.
//
// Test strategy:
// 1. Test standard video extensions (mp4, mkv) with default extensions list
// 2. Test custom extension lists for specialized use cases
// 3. Test non-video extensions to ensure proper rejection
// 4. Test edge cases like uppercase extensions and extensions without leading dots
func TestIsVideoFile(t *testing.T) {
	tests := []struct {
		name       string // Description of the test case
		path       string // Filepath to test
		extensions []string // Optional custom extensions list
		want       bool // Expected result
	}{
		{
			name:       "Standard MP4",
			path:       "video.mp4",
			extensions: nil,
			want:       true,
		},
		{
			name:       "Standard MKV",
			path:       "video.mkv",
			extensions: nil,
			want:       true,
		},
		{
			name:       "Custom Extension List",
			path:       "video.mp4",
			extensions: []string{".mp4", ".mkv"},
			want:       true,
		},
		{
			name:       "Case Insensitive",
			path:       "video.MP4",
			extensions: nil,
			want:       true,
		},
		{
			name:       "Non-Video Extension",
			path:       "document.pdf",
			extensions: nil,
			want:       false,
		},
		{
			name:       "Empty Extension",
			path:       "noextension",
			extensions: nil,
			want:       false,
		},
		{
			name:       "Extension Without Dot",
			path:       "video.mp4",
			extensions: []string{"mp4", "mkv"},
			want:       true,
		},
		{
			name:       "Not In Custom List",
			path:       "video.avi",
			extensions: []string{".mp4", ".mkv"},
			want:       false,
		},
	}

	// Run all the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsVideoFile(tt.path, tt.extensions); got != tt.want {
				t.Errorf("IsVideoFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestFormatFrameRate verifies the frame rate formatting functionality.
// These tests ensure VidKit correctly converts raw frame rate values from
// ffprobe (typically in fractional form) into human-readable formats.
//
// Test strategy:
// 1. Test common frame rates (24, 30, 60 fps)
// 2. Test fractional frame rates (e.g., 24000/1001 for 23.976 fps)
// 3. Test edge cases like empty strings and invalid formats
func TestFormatFrameRate(t *testing.T) {
	tests := []struct {
		name      string // Description of the test case
		frameRate string // Input frame rate (typically from ffprobe)
		want      string // Expected formatted output
	}{
		{
			name:      "24 fps",
			frameRate: "24/1",
			want:      "24.00 fps",
		},
		{
			name:      "30 fps",
			frameRate: "30/1",
			want:      "30.00 fps",
		},
		{
			name:      "23.976 fps (24000/1001)",
			frameRate: "24000/1001",
			want:      "23.98 fps",
		},
		{
			name:      "29.97 fps (30000/1001)",
			frameRate: "30000/1001",
			want:      "29.97 fps",
		},
		{
			name:      "60 fps",
			frameRate: "60/1",
			want:      "60.00 fps",
		},
		{
			name:      "Empty string",
			frameRate: "",
			want:      "N/A",
		},
		{
			name:      "Invalid format",
			frameRate: "not-a-framerate",
			want:      "not-a-framerate",
		},
	}

	// Run all the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatFrameRate(tt.frameRate); got != tt.want {
				t.Errorf("FormatFrameRate() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestFormatBitRate verifies the bit rate formatting functionality.
// These tests ensure VidKit correctly converts raw bit rate values from
// ffprobe (typically in bits per second) into human-readable formats.
//
// Test strategy:
// 1. Test various bit rates from very low to very high
// 2. Test edge cases like empty strings and invalid formats
func TestFormatBitRate(t *testing.T) {
	tests := []struct {
		name    string // Description of the test case
		bitRate string // Input bit rate (typically from ffprobe)
		want    string // Expected formatted output
	}{
		{
			name:    "Standard bit rate",
			bitRate: "1500000",
			want:    "1500.00 Kbps",
		},
		{
			name:    "High bit rate",
			bitRate: "15000000",
			want:    "15000.00 Kbps",
		},
		{
			name:    "Low bit rate",
			bitRate: "128000",
			want:    "128.00 Kbps",
		},
		{
			name:    "Empty string",
			bitRate: "",
			want:    "N/A",
		},
		{
			name:    "Invalid format",
			bitRate: "not-a-bitrate",
			want:    "not-a-bitrate bps",
		},
	}

	// Run all the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatBitRate(tt.bitRate); got != tt.want {
				t.Errorf("FormatBitRate() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestFormatSize verifies the file size formatting functionality.
// These tests ensure VidKit correctly converts raw file sizes from
// ffprobe (typically in bytes) into human-readable formats with appropriate units.
//
// Test strategy:
// 1. Test various file sizes from bytes to gigabytes
// 2. Test boundary conditions between units
// 3. Test edge cases like empty strings and invalid formats
func TestFormatSize(t *testing.T) {
	tests := []struct {
		name string // Description of the test case
		size string // Input file size in bytes (typically from ffprobe)
		want string // Expected formatted output
	}{
		{
			name: "Bytes",
			size: "100",
			want: "100 B",
		},
		{
			name: "Kilobytes",
			size: "1024",
			want: "1.00 KB",
		},
		{
			name: "Megabytes",
			size: "1048576", // 1024 * 1024
			want: "1.00 MB",
		},
		{
			name: "Gigabytes",
			size: "1073741824", // 1024 * 1024 * 1024
			want: "1.00 GB",
		},
		{
			name: "Boundary KB/MB",
			size: "1048575", // Just below 1MB
			want: "1024.00 KB",
		},
		{
			name: "Empty string",
			size: "",
			want: "N/A",
		},
		{
			name: "Invalid format",
			size: "not-a-size",
			want: "not-a-size bytes",
		},
	}

	// Run all the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatSize(tt.size); got != tt.want {
				t.Errorf("FormatSize() = %v, want %v", got, tt.want)
			}
		})
	}
}
