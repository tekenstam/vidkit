#!/bin/bash
# Script to test VidKit's directory organization features
# This script demonstrates different organization patterns using the test videos

# Set up test environment
echo "Setting up test environment..."
mkdir -p test_results
rm -rf test_results/*

# Build VidKit if it doesn't exist
if [ ! -f "./vidkit" ]; then
  echo "Building VidKit..."
  go build -o vidkit ./cmd/vidkit
fi

# Ensure we have test videos
if [ ! -d "test_videos" ] || [ -z "$(ls -A test_videos)" ]; then
  echo "Generating test videos..."
  ./tools/generate_test_videos.sh
fi

# Test organization patterns

echo -e "\n=== Testing Movie Organization ===\n"

# Test basic movie organization
echo -e "\n--- Basic Movie Organization ---"
./vidkit --preview --movie-dir "test_results/Movies/{title} ({year})" test_videos/The.Shawshank.Redemption.1994.1080p.BluRay.x264.mp4

# Test alphabetical organization
echo -e "\n--- Alphabetical Movie Organization ---"
./vidkit --preview --movie-dir "test_results/Movies/{title[0]}/{title} ({year})" test_videos/The.Shawshank.Redemption.1994.1080p.BluRay.x264.mp4

# Test genre-based organization
echo -e "\n--- Genre-Based Movie Organization ---"
./vidkit --preview --movie-dir "test_results/Movies/By Genre/{genre}/{title} ({year})" test_videos/The.Shawshank.Redemption.1994.1080p.BluRay.x264.mp4

echo -e "\n=== Testing TV Show Organization ===\n"

# Test basic TV show organization
echo -e "\n--- Basic TV Show Organization ---"
./vidkit --preview --tv-dir "test_results/TV Shows/{title}/Season {season:02d}" test_videos/Breaking.Bad.S01E01.Pilot.mp4

# Test network-based organization
echo -e "\n--- Network-Based TV Show Organization ---"
./vidkit --preview --tv-dir "test_results/TV Shows/{network}/{title}/Season {season}" test_videos/Breaking.Bad.S01E01.Pilot.mp4

# Test year-based organization
echo -e "\n--- Year-Based TV Show Organization ---"
./vidkit --preview --tv-dir "test_results/TV Shows/By Year/{year}/{title}/Season {season:02d}" test_videos/Breaking.Bad.S01E01.Pilot.mp4

echo -e "\n=== Testing Multiple Provider Organization ===\n"

# Test OMDb provider with movie organization
echo -e "\n--- OMDb Movie Provider with Organization ---"
./vidkit --preview --movie-provider omdb --movie-dir "test_results/Movies (OMDb)/{genre}/{title} ({year})" test_videos/The.Shawshank.Redemption.1994.1080p.BluRay.x264.mp4

# Test TVDb provider with TV organization
echo -e "\n--- TVDb TV Provider with Organization ---"
./vidkit --preview --tv-provider tvdb --tv-dir "test_results/TV Shows (TVDb)/{network}/{title}/S{season:02d}" test_videos/Breaking.Bad.S01E01.Pilot.mp4

echo -e "\n=== Test Complete ===\n"
echo "The above examples show how VidKit would organize files with different templates."
echo "Use the --preview flag to see the changes without applying them."
echo "Remove the --preview flag to actually create the directories and move the files."
