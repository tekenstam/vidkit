// Package media provides functionality for video file analysis and manipulation.
// It contains tools for extracting technical metadata from video files using FFmpeg,
// detecting video codecs, resolutions, and other stream information.
//
// This package is central to VidKit's ability to analyze video files and is used
// throughout the application to determine video quality, format information, and 
// technical properties needed for intelligent file organization.
package media

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// VideoInfo represents the structure of ffprobe JSON output with comprehensive
// information about a video file's format and streams.
// 
// The structure closely follows FFmpeg's ffprobe JSON output format, containing
// both general format information (file size, duration, bit rate) and detailed
// stream information (video resolution, codec, audio channels, etc.).
type VideoInfo struct {
	Format struct {
		Filename   string `json:"filename"`   // Complete path to the video file
		FormatName string `json:"format_name"` // Container format (e.g., "mp4", "mkv")
		Duration   string `json:"duration"`   // Duration in seconds as a string
		Size       string `json:"size"`       // File size in bytes as a string
		BitRate    string `json:"bit_rate"`   // Total bit rate in bits/second
		ProbeScore int    `json:"probe_score"` // Confidence score of format detection (higher is better)
	} `json:"format"`
	Streams []struct {
		CodecType     string `json:"codec_type"`     // Type of stream ("video", "audio", "subtitle")
		CodecName     string `json:"codec_name"`     // Codec name (e.g., "h264", "aac")
		Width         int    `json:"width,omitempty"`  // Video width in pixels
		Height        int    `json:"height,omitempty"` // Video height in pixels
		BitRate       string `json:"bit_rate,omitempty"` // Stream bit rate in bits/second
		FrameRate     string `json:"r_frame_rate,omitempty"` // Frame rate as a fraction (e.g., "24000/1001")
		SampleRate    string `json:"sample_rate,omitempty"`  // Audio sample rate in Hz
		Channels      int    `json:"channels,omitempty"`     // Number of audio channels
		ChannelLayout string `json:"channel_layout,omitempty"` // Audio channel layout (e.g., "stereo")
	} `json:"streams"`
}

// GetVideoInfo retrieves detailed video file information using FFmpeg's ffprobe tool.
// It executes ffprobe with JSON output format and parses the results into a structured
// VideoInfo object that can be easily processed by the application.
//
// The function requires ffprobe to be installed and available in the system PATH.
// It captures comprehensive information about video and audio streams, including
// resolution, codec, bit rate, and other technical properties.
//
// Example:
//
//	info, err := media.GetVideoInfo("/path/to/video.mp4")
//	if err != nil {
//	    log.Fatalf("Failed to analyze video: %v", err)
//	}
//	fmt.Printf("Resolution: %dx%d\n", info.Streams[0].Width, info.Streams[0].Height)
func GetVideoInfo(filename string) (*VideoInfo, error) {
	// Execute ffprobe to analyze the video file
	// -v quiet: Suppress unnecessary output
	// -print_format json: Output in JSON format for easy parsing
	// -show_format: Include container format information
	// -show_streams: Include detailed stream information
	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		filename)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ffprobe failed: %v", err)
	}

	var info VideoInfo
	if err := json.Unmarshal(output, &info); err != nil {
		return nil, fmt.Errorf("failed to parse ffprobe output: %v", err)
	}

	return &info, nil
}

// IsVideoFile checks if a file has a video extension, determining whether
// VidKit should process it as a video.
//
// The function supports two modes of operation:
// 1. With a custom list of extensions (extensions parameter)
// 2. With a default set of common video extensions
//
// The function normalizes extensions to lowercase and handles both formats:
// extensions with or without leading dots (e.g., both ".mp4" and "mp4" are valid).
//
// Example:
//
//	// Using default extensions
//	isVideo := media.IsVideoFile("movie.mp4", nil)
//	
//	// Using custom extensions
//	customExts := []string{".mp4", ".mkv", ".mov"}
//	isVideo := media.IsVideoFile("movie.avi", customExts)
func IsVideoFile(path string, extensions []string) bool {
	// If custom extensions are provided, check against those
	if len(extensions) > 0 {
		ext := strings.ToLower(filepath.Ext(path))
		for _, validExt := range extensions {
			// Handle extensions with or without leading dot
			if strings.HasPrefix(validExt, ".") {
				if ext == validExt {
					return true
				}
			} else if ext == "."+validExt {
				return true
			}
		}
		return false
	}

	// Default video extensions if none specified
	// This map provides O(1) lookup performance
	videoExts := map[string]bool{
		".mp4":  true,
		".mkv":  true,
		".avi":  true,
		".mov":  true,
		".wmv":  true,
		".flv":  true,
		".webm": true,
		".m4v":  true,
		".mpg":  true,
		".mpeg": true,
		".3gp":  true,
	}
	ext := strings.ToLower(filepath.Ext(path))
	return videoExts[ext]
}

// FormatBitRate formats a raw bit rate string into a human-readable format.
// It converts the raw bit rate (typically in bits per second) to a more
// readable format in Kbps (kilobits per second) with appropriate rounding.
//
// The function handles edge cases such as:
// - Empty bit rate values (returns "N/A")
// - Non-numeric bit rate values (appends " bps" to the original string)
//
// Example:
//
//	rawBitRate := "1500000"
//	formattedRate := media.FormatBitRate(rawBitRate)
//	// Result: "1500.00 Kbps"
func FormatBitRate(bitRate string) string {
	if bitRate == "" {
		return "N/A"
	}
	rate, err := strconv.ParseFloat(bitRate, 64)
	if err != nil {
		return bitRate + " bps"
	}
	rate /= 1000 // Convert to Kbps
	return fmt.Sprintf("%.2f Kbps", rate)
}

// FormatFileSize converts a file size in bytes to a human-readable string
// with appropriate units (B, KB, MB, GB).
//
// This function intelligently selects the most appropriate unit based on the
// file size, making it easier to present file sizes in user interfaces.
//
// Example:
//
//	size := "1048576"  // 1 MB in bytes
//	formatted := media.FormatFileSize(size)
//	// Result: "1.00 MB"
func FormatFileSize(sizeStr string) string {
	if sizeStr == "" {
		return "N/A"
	}
	size, err := strconv.ParseFloat(sizeStr, 64)
	if err != nil {
		return sizeStr + " bytes"
	}

	const unit = 1024.0
	if size < unit {
		return fmt.Sprintf("%.0f B", size)
	}
	div, exp := unit, 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", size/div, "KMG"[exp])
}

// FormatSize formats file size in a human-readable format.
// This is an alias for FormatFileSize for backward compatibility.
//
// Deprecated: Use FormatFileSize instead.
func FormatSize(sizeStr string) string {
	return FormatFileSize(sizeStr)
}

// FormatFrameRate formats a raw frame rate string into a human-readable format.
// It converts the raw frame rate (typically in the form of a fraction) to a more
// readable format in frames per second (FPS) with appropriate rounding.
//
// The function handles edge cases such as:
// - Empty frame rate values (returns "N/A")
// - Non-fraction frame rate values (returns the original string)
//
// Example:
//
//	rawFrameRate := "24000/1001"
//	formattedRate := media.FormatFrameRate(rawFrameRate)
//	// Result: "23.98 fps"
func FormatFrameRate(frameRate string) string {
	if frameRate == "" {
		return "N/A"
	}
	parts := strings.Split(frameRate, "/")
	if len(parts) != 2 {
		return frameRate
	}
	num, err1 := strconv.ParseFloat(parts[0], 64)
	den, err2 := strconv.ParseFloat(parts[1], 64)
	if err1 != nil || err2 != nil || den == 0 {
		return frameRate
	}
	return fmt.Sprintf("%.2f fps", num/den)
}
