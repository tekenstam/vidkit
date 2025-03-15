package resolution

import "fmt"

// Resolution represents a standard video resolution
type Resolution struct {
	Name   string
	Width  int
	Height int
}

var StandardResolutions = []Resolution{
	{"8K", 7680, 4320},
	{"4K", 3840, 2160},
	{"1440p", 2560, 1440},
	{"1080p", 1920, 1080},
	{"2K", 2048, 1080},
	{"720p", 1280, 720},
	{"480p", 640, 480},
	{"360p", 640, 360},
}

// GetStandardResolution returns the standard resolution name for given dimensions
func GetStandardResolution(width, height int) string {
	if height == 0 {
		return "0p"
	}

	// First check for exact matches
	for _, res := range StandardResolutions {
		if width == res.Width && height == res.Height {
			return res.Name
		}
	}

	// Check for standard heights with tolerance
	tolerance := 10
	for _, res := range StandardResolutions {
		heightMatch := abs(height-res.Height) <= tolerance
		if heightMatch {
			// For any resolution with matching height (including ultrawide),
			// use the standard name based on height
			return res.Name
		}
	}

	// For non-standard resolutions, use the vertical resolution
	return fmt.Sprintf("%dp", height)
}

// abs returns the absolute value of x
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
