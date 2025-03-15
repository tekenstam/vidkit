#!/bin/bash
# Script to test VidKit's directory organization features
# This script demonstrates different organization patterns using the test videos
#
# Usage:
#   ./test_directory_organization.sh           # Run with metadata lookups (requires API keys)
#   ./test_directory_organization.sh --no-metadata  # Run without metadata lookups (for CI)

# Set up consistent coloring for test output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

# Setup exit code tracking
TESTS_FAILED=0

# Check for --no-metadata flag
NO_METADATA=""
if [[ "$*" == *"--no-metadata"* ]]; then
  NO_METADATA="--no-metadata"
  echo "Running in no-metadata mode (offline)"
fi

# Detect if running in CI environment
IN_CI=0
if [ -n "$CI" ] || [ -n "$GITHUB_ACTIONS" ]; then
  IN_CI=1
  echo "Detected CI environment, adapting tests accordingly"
fi

# Helper function to report test results
report_test_result() {
  local test_name="$1"
  local result="$2"
  local error_msg="$3"
  
  if [ "$result" -eq 0 ]; then
    echo -e "${GREEN}✅ Passed:${NC} $test_name"
  else
    echo -e "${RED}❌ Failed:${NC} $test_name"
    echo -e "${YELLOW}Error:${NC} $error_msg"
    TESTS_FAILED=1
  fi
}

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

# Helper function to test movie organization pattern
test_movie_organization() {
  local test_name="$1"
  local pattern="$2"
  local video_file="$3"
  
  echo -e "\n--- $test_name ---"
  echo "Testing pattern: $pattern"
  ./vidkit --preview $NO_METADATA --movie-dir "$pattern" "$video_file" 2>&1 | tee output.log
  
  # Check output for basic organization information
  if grep -q "would be organized" output.log || grep -q "test_results" output.log || grep -q "Resolution:" output.log; then
    report_test_result "$test_name" 0 ""
  else
    report_test_result "$test_name" 1 "Failed to organize using pattern: $pattern"
  fi
}

# Helper function to test TV show organization pattern
test_tv_organization() {
  local test_name="$1"
  local pattern="$2"
  local video_file="$3"
  
  echo -e "\n--- $test_name ---"
  echo "Testing pattern: $pattern"
  ./vidkit --preview $NO_METADATA --tv-dir "$pattern" "$video_file" 2>&1 | tee output.log
  
  # Check output for basic organization information
  if grep -q "would be organized" output.log || grep -q "test_results" output.log || grep -q "Resolution:" output.log; then
    report_test_result "$test_name" 0 ""
  else
    report_test_result "$test_name" 1 "Failed to organize using pattern: $pattern"
  fi
}

# Test organization patterns
echo -e "\n=== Testing Movie Organization ===\n"

# Find a test movie
MOVIE_FILE=""
for vid in test_videos/*.mp4; do
  if grep -q "Shawshank\|Inception\|Dark.Knight" <<< "$vid"; then
    MOVIE_FILE=$vid
    break
  fi
done

if [ -z "$MOVIE_FILE" ]; then
  # Just use any video if we can't find a specific movie
  MOVIE_FILE=$(ls test_videos/*.mp4 | head -1)
  if [ -z "$MOVIE_FILE" ]; then
    report_test_result "Find test movie" 1 "No test videos found. Cannot continue."
    exit 1
  fi
fi

echo "Using movie file: $MOVIE_FILE"

# Test basic movie organization
test_movie_organization "Basic Movie Organization" "test_results/Movies/{title} ({year})" "$MOVIE_FILE"

# Test alphabetical organization
test_movie_organization "Alphabetical Movie Organization" "test_results/Movies/{title[0]}/{title} ({year})" "$MOVIE_FILE"

# Test genre-based organization
test_movie_organization "Genre-Based Movie Organization" "test_results/Movies/By Genre/{genre}/{title} ({year})" "$MOVIE_FILE"

echo -e "\n=== Testing TV Show Organization ===\n"

# Find a test TV episode
TV_FILE=""
for vid in test_videos/*.mp4; do
  if grep -q "Breaking.Bad\|Game.of.Thrones\|S[0-9][0-9]E[0-9][0-9]" <<< "$vid"; then
    TV_FILE=$vid
    break
  fi
done

if [ -z "$TV_FILE" ]; then
  # Just use any video if we can't find a specific TV episode
  TV_FILE=$(ls test_videos/*.mp4 | head -1)
  if [ -z "$TV_FILE" ]; then
    report_test_result "Find test TV episode" 1 "No test videos found. Cannot continue."
    exit 1
  fi
fi

echo "Using TV file: $TV_FILE"

# Test basic TV show organization
test_tv_organization "Basic TV Show Organization" "test_results/TV Shows/{title}/Season {season:02d}" "$TV_FILE"

# Test network-based organization
test_tv_organization "Network-Based TV Show Organization" "test_results/TV Shows/{network}/{title}/Season {season}" "$TV_FILE"

# Test year-based organization
test_tv_organization "Year-Based TV Show Organization" "test_results/TV Shows/By Year/{year}/{title}/Season {season:02d}" "$TV_FILE"

# Skip provider-specific tests if running in no-metadata mode
if [ -z "$NO_METADATA" ]; then
  echo -e "\n=== Testing Multiple Provider Organization ===\n"

  # Test OMDb provider with movie organization
  echo -e "\n--- OMDb Movie Provider with Organization ---"
  ./vidkit --preview --movie-provider omdb --movie-dir "test_results/Movies (OMDb)/{genre}/{title} ({year})" "$MOVIE_FILE" 2>&1 | tee output.log
  if grep -q "would be organized" output.log || grep -q "test_results" output.log; then
    report_test_result "OMDb Provider Organization" 0 ""
  else
    report_test_result "OMDb Provider Organization" 1 "Failed to organize using OMDb provider"
  fi

  # Test TVDb provider with TV organization
  echo -e "\n--- TVDb TV Provider with Organization ---"
  ./vidkit --preview --tv-provider tvdb --tv-dir "test_results/TV Shows (TVDb)/{network}/{title}/S{season:02d}" "$TV_FILE" 2>&1 | tee output.log
  if grep -q "would be organized" output.log || grep -q "test_results" output.log; then
    report_test_result "TVDb Provider Organization" 0 ""
  else
    report_test_result "TVDb Provider Organization" 1 "Failed to organize using TVDb provider"
  fi
else
  echo -e "\n${YELLOW}=== Skipping Provider-Specific Tests (No Metadata Mode) ===${NC}\n"
fi

echo -e "\n=== Directory Organization Tests Complete ===\n"
if [ $TESTS_FAILED -eq 0 ]; then
  echo -e "${GREEN}All directory organization tests passed successfully!${NC}"
  echo "The above examples show how VidKit would organize files with different templates."
  echo "Use the --preview flag to see the changes without applying them."
  echo "Remove the --preview flag to actually create the directories and move the files."
  exit 0
else
  echo -e "${RED}Some directory organization tests failed!${NC}"
  exit 1
fi
