package media

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// VideoInfo represents the structure of ffprobe JSON output
type VideoInfo struct {
	Format struct {
		Filename    string `json:"filename"`
		FormatName  string `json:"format_name"`
		Duration    string `json:"duration"`
		Size        string `json:"size"`
		BitRate     string `json:"bit_rate"`
		ProbeScore  int    `json:"probe_score"`
	} `json:"format"`
	Streams []struct {
		CodecType     string `json:"codec_type"`
		CodecName     string `json:"codec_name"`
		Width         int    `json:"width,omitempty"`
		Height        int    `json:"height,omitempty"`
		BitRate       string `json:"bit_rate,omitempty"`
		FrameRate     string `json:"r_frame_rate,omitempty"`
		SampleRate    string `json:"sample_rate,omitempty"`
		Channels      int    `json:"channels,omitempty"`
		ChannelLayout string `json:"channel_layout,omitempty"`
	} `json:"streams"`
}

// GetVideoInfo retrieves video file information using ffprobe
func GetVideoInfo(filename string) (*VideoInfo, error) {
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

// IsVideoFile checks if the file has a video extension
func IsVideoFile(path string, extensions []string) bool {
	if len(extensions) > 0 {
		ext := strings.ToLower(filepath.Ext(path))
		for _, validExt := range extensions {
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

// FormatBitRate formats bit rate in a human-readable format
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

// FormatSize formats file size in a human-readable format
func FormatSize(sizeStr string) string {
	size, err := strconv.ParseFloat(sizeStr, 64)
	if err != nil {
		return sizeStr + " bytes"
	}
	
	units := []string{"B", "KB", "MB", "GB", "TB"}
	unitIndex := 0
	
	for size >= 1024 && unitIndex < len(units)-1 {
		size /= 1024
		unitIndex++
	}
	
	return fmt.Sprintf("%.2f %s", size, units[unitIndex])
}

// FormatFrameRate formats frame rate in a human-readable format
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
