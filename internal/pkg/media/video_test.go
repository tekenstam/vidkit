package media

import (
	"testing"
)

func TestIsVideoFile(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		extensions []string
		want       bool
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
			name:       "Custom Extension Not in List",
			path:       "video.avi",
			extensions: []string{".mp4", ".mkv"},
			want:       false,
		},
		{
			name:       "Non-Video Extension",
			path:       "document.txt",
			extensions: nil,
			want:       false,
		},
		{
			name:       "No Extension",
			path:       "video",
			extensions: nil,
			want:       false,
		},
		{
			name:       "Mixed Case Extension",
			path:       "video.MP4",
			extensions: nil,
			want:       true,
		},
		{
			name:       "Custom Extension Without Dot",
			path:       "video.mp4",
			extensions: []string{"mp4"},
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsVideoFile(tt.path, tt.extensions); got != tt.want {
				t.Errorf("IsVideoFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatBitRate(t *testing.T) {
	tests := []struct {
		name    string
		bitRate string
		want    string
	}{
		{
			name:    "Empty BitRate",
			bitRate: "",
			want:    "N/A",
		},
		{
			name:    "Valid BitRate",
			bitRate: "1000000",
			want:    "1000.00 Kbps",
		},
		{
			name:    "Invalid BitRate",
			bitRate: "not_a_number",
			want:    "not_a_number bps",
		},
		{
			name:    "Zero BitRate",
			bitRate: "0",
			want:    "0.00 Kbps",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatBitRate(tt.bitRate); got != tt.want {
				t.Errorf("FormatBitRate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatSize(t *testing.T) {
	tests := []struct {
		name string
		size string
		want string
	}{
		{
			name: "Bytes",
			size: "500",
			want: "500.00 B",
		},
		{
			name: "Kilobytes",
			size: "1500",
			want: "1.46 KB",
		},
		{
			name: "Megabytes",
			size: "1500000",
			want: "1.43 MB",
		},
		{
			name: "Gigabytes",
			size: "1500000000",
			want: "1.40 GB",
		},
		{
			name: "Invalid Size",
			size: "not_a_number",
			want: "not_a_number bytes",
		},
		{
			name: "Zero Size",
			size: "0",
			want: "0.00 B",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatSize(tt.size); got != tt.want {
				t.Errorf("FormatSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatFrameRate(t *testing.T) {
	tests := []struct {
		name      string
		frameRate string
		want      string
	}{
		{
			name:      "Empty FrameRate",
			frameRate: "",
			want:      "N/A",
		},
		{
			name:      "Valid FrameRate",
			frameRate: "24000/1000",
			want:      "24.00 fps",
		},
		{
			name:      "Invalid Format",
			frameRate: "not_a_framerate",
			want:      "not_a_framerate",
		},
		{
			name:      "Zero Denominator",
			frameRate: "24000/0",
			want:      "24000/0",
		},
		{
			name:      "Invalid Numbers",
			frameRate: "invalid/invalid",
			want:      "invalid/invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatFrameRate(tt.frameRate); got != tt.want {
				t.Errorf("FormatFrameRate() = %v, want %v", got, tt.want)
			}
		})
	}
}
