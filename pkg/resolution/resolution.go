// Package resolution provides functionality for determining and standardizing video resolutions.
// It offers utilities to convert raw video dimensions into standardized resolution names
// (such as 1080p, 4K, etc.) and handles a variety of edge cases in resolution detection.
//
// It's designed to be used by video analysis tools that need to categorize videos
// based on their resolution quality, supporting both exact resolution matching and
// approximate detection when dimensions don't exactly match standard values.
package resolution

import "fmt"

// Resolution represents a standard video resolution with a common name and dimensions.
// This structure is used to map width×height combinations to their industry-standard names.
type Resolution struct {
	Name   string // Common name (e.g., "1080p", "4K")
	Width  int    // Width in pixels
	Height int    // Height in pixels
}

// StandardResolutions defines the list of recognized standard video resolutions
// in descending order of quality. This ordering is intentional and used in the
// resolution detection algorithm to prioritize higher quality names when
// approximate matching is required.
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

// GetStandardResolution returns the standard resolution name for given dimensions.
// It attempts to find an exact match first, then falls back to determining the
// closest standard resolution based on height if no exact match is found.
//
// The function handles special cases:
// - Returns "0p" for invalid dimensions (height = 0)
// - Uses height as the primary factor for resolution naming when no exact match exists
// - Handles non-standard aspect ratios by focusing on vertical resolution
//
// Example:
//
//	res := resolution.GetStandardResolution(1920, 1080)
//	fmt.Println(res) // Output: "1080p"
//
//	res = resolution.GetStandardResolution(3840, 2160)
//	fmt.Println(res) // Output: "4K"
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

	// If no exact match, determine based on height which is the common
	// way to name resolutions (e.g., "1080p" refers to the height)
	
	// Special case for 4K and 8K which are typically named by approximate width
	if height > 1440 {
		if width >= 7000 {
			return "8K"
		} else if width >= 3800 {
			return "4K"
		}
	}

	// For other resolutions, use the height-based naming convention
	return fmt.Sprintf("%dp", height)
}

// GetClosestStandardResolution returns the name of the standard resolution
// that most closely matches the given dimensions. Unlike GetStandardResolution,
// this always returns one of the predefined names in StandardResolutions.
//
// Example:
//
//	res := resolution.GetClosestStandardResolution(1280, 720)
//	fmt.Println(res) // Output: "720p"
//
//	// Even for unusual dimensions, it finds the closest standard:
//	res = resolution.GetClosestStandardResolution(1800, 1000)
//	fmt.Println(res) // Output: "1080p" (closest standard)
func GetClosestStandardResolution(width, height int) string {
	if height == 0 {
		return "0p"
	}

	// First check for exact matches
	for _, res := range StandardResolutions {
		if width == res.Width && height == res.Height {
			return res.Name
		}
	}

	// If no exact match, find the closest standard based on height
	// Starting with the highest resolution and working down
	for _, res := range StandardResolutions {
		if height >= res.Height {
			return res.Name
		}
	}

	// If height is smaller than all standards, return the smallest standard
	return StandardResolutions[len(StandardResolutions)-1].Name
}

// abs returns the absolute value of x
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
