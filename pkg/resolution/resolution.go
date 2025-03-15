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
	{"2K", 2048, 1080},
	{"1440p", 2560, 1440},
	{"1080p", 1920, 1080},
	{"720p", 1280, 720},
	{"480p", 640, 480},
	{"360p", 640, 360},
}

// GetStandardResolution returns the standard resolution name for given dimensions
func GetStandardResolution(width, height int) string {
	// First check for exact matches
	for _, res := range StandardResolutions {
		if width == res.Width && height == res.Height {
			return fmt.Sprintf("%s", res.Name)
		}
	}

	// If no exact match, find the closest match based on vertical resolution
	for _, res := range StandardResolutions {
		if height >= res.Height-10 && height <= res.Height+10 {
			return fmt.Sprintf("%s", res.Name)
		}
	}

	// For very low resolutions, just return the vertical resolution
	if height < 360 {
		return fmt.Sprintf("%dp", height)
	}

	// For non-standard resolutions between standard ones,
	// find the closest standard resolution and mark it as custom
	var closestRes Resolution
	minDiff := 10000
	for _, res := range StandardResolutions {
		diff := abs(height - res.Height)
		if diff < minDiff {
			minDiff = diff
			closestRes = res
		}
	}

	// If within 20% of a standard resolution, use that name with a "~" prefix
	if float64(minDiff)/float64(closestRes.Height) <= 0.2 {
		return fmt.Sprintf("~%s", closestRes.Name)
	}

	// Otherwise just return the vertical resolution
	return fmt.Sprintf("%dp", height)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
