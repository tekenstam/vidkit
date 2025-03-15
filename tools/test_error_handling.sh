#!/bin/bash
# Test script for VidKit's error handling capabilities
# 
# Usage:
#   ./test_error_handling.sh           # Run with metadata lookups (requires API keys)
#   ./test_error_handling.sh --no-metadata  # Run without metadata lookups (for CI)

# Check for --no-metadata flag
NO_METADATA=""
if [[ "$*" == *"--no-metadata"* ]]; then
  NO_METADATA="--no-metadata"
  echo "Running in no-metadata mode (offline)"
fi

# Set up test environment
echo "=== Testing Error Handling ==="
mkdir -p test_results/errors
rm -rf test_results/errors/*

# Build VidKit if needed
if [ ! -f "./vidkit" ]; then
  echo "Building VidKit..."
  go build -o vidkit ./cmd/vidkit
fi

# Create a directory for test files
mkdir -p test_results/errors/input

# Test 1: Non-existent file
echo -e "\n--- Test 1: Non-existent file ---"
./vidkit --preview --batch $NO_METADATA test_results/errors/nonexistent_file.mp4 2>&1 | tee output.log
if ! grep -q "Error" output.log; then
  echo "❌ Failed: Should report error for non-existent file"
  exit 1
else
  echo "✅ Passed: Correctly reported error for non-existent file"
fi

# Test 2: Invalid file format
echo -e "\n--- Test 2: Invalid file format ---"
# Create an invalid "video" file (text file with mp4 extension)
echo "This is not a valid video file" > test_results/errors/input/invalid.mp4
./vidkit --preview --batch $NO_METADATA test_results/errors/input/invalid.mp4 2>&1 | tee output.log
if ! grep -q "Error\|Failed\|Invalid" output.log; then
  echo "❌ Failed: Should report error for invalid file format"
  exit 1
else
  echo "✅ Passed: Correctly reported error for invalid file format"
fi

# Test 3: Read-only directory
echo -e "\n--- Test 3: Read-only output directory ---"
# Create a read-only directory to test write permission errors
mkdir -p test_results/errors/readonly
chmod 555 test_results/errors/readonly
if [ -d "test_videos" ] && [ -n "$(ls -A test_videos/*.mp4 2>/dev/null)" ]; then
  # Find a test video
  TEST_VIDEO=$(ls test_videos/*.mp4 | head -1)
  ./vidkit --preview --batch $NO_METADATA --movie-directory-template "test_results/errors/readonly" "$TEST_VIDEO" 2>&1 | tee output.log
  if ! grep -q "Error\|Permission\|Failed" output.log; then
    echo "❌ Failed: Should report error for read-only directory"
  else
    echo "✅ Passed: Correctly reported error for read-only directory"
  fi
else
  echo "⚠️ Skipping read-only test: No test videos found"
fi
# Reset directory permissions
chmod 755 test_results/errors/readonly

# Test 4: Invalid metadata provider
echo -e "\n--- Test 4: Invalid metadata provider ---"
if [ -d "test_videos" ] && [ -n "$(ls -A test_videos/*.mp4 2>/dev/null)" ]; then
  # Find a test video
  TEST_VIDEO=$(ls test_videos/*.mp4 | head -1)
  ./vidkit --preview --batch $NO_METADATA --movie-provider invalid_provider "$TEST_VIDEO" 2>&1 | tee output.log
  if ! grep -q "Error\|Invalid\|provider" output.log; then
    echo "❌ Failed: Should report error for invalid provider"
    exit 1
  else
    echo "✅ Passed: Correctly reported error for invalid provider"
  fi
else
  echo "⚠️ Skipping provider test: No test videos found"
fi

# Test 5: Missing API key when needed
echo -e "\n--- Test 5: Missing API key when needed ---"
if [ -d "test_videos" ] && [ -n "$(ls -A test_videos/*.mp4 2>/dev/null)" ]; then
  # Find a test movie
  TEST_MOVIE=
  for vid in test_videos/*.mp4; do
    if grep -q "Shawshank\|Inception\|Dark.Knight" <<< "$vid"; then
      TEST_MOVIE=$vid
      break
    fi
  done
  
  if [ -z "$TEST_MOVIE" ]; then
    # Just use any video if we can't find a specific movie
    TEST_MOVIE=$(ls test_videos/*.mp4 | head -1)
  fi
  
  # Temporarily move any config file to ensure API key is not provided
  if [ -f "$HOME/.config/vidkit/config.json" ]; then
    mv "$HOME/.config/vidkit/config.json" "$HOME/.config/vidkit/config.json.bak"
  fi
  
  # Try to access TMDb without API key
  ./vidkit --preview --batch $NO_METADATA --movie-provider tmdb --no-metadata=false "$TEST_MOVIE" 2>&1 | tee output.log
  if ! grep -q "API key\|Error\|Missing" output.log; then
    echo "❌ Failed: Should report error for missing API key"
  else
    echo "✅ Passed: Correctly reported error for missing API key"
  fi
  
  # Restore config file if it existed
  if [ -f "$HOME/.config/vidkit/config.json.bak" ]; then
    mv "$HOME/.config/vidkit/config.json.bak" "$HOME/.config/vidkit/config.json"
  fi
else
  echo "⚠️ Skipping API key test: No test videos found"
fi

# Test 6: Invalid format string
echo -e "\n--- Test 6: Invalid format string ---"
if [ -d "test_videos" ] && [ -n "$(ls -A test_videos/*.mp4 2>/dev/null)" ]; then
  # Find a test video
  TEST_VIDEO=$(ls test_videos/*.mp4 | head -1)
  ./vidkit --preview --batch $NO_METADATA --movie-filename-template "{invalid_variable}" "$TEST_VIDEO" 2>&1 | tee output.log
  # Check if it handles invalid variables in the format string
  if grep -q "panic\|crash" output.log; then
    echo "❌ Failed: Program crashed on invalid format string"
    exit 1
  else
    echo "✅ Passed: Gracefully handled invalid format string"
  fi
else
  echo "⚠️ Skipping format string test: No test videos found"
fi

# Test 7: File with unusual characters in name
echo -e "\n--- Test 7: File with unusual characters in name ---"
# Create a file with special characters
SPECIAL_CHAR_FILE="test_results/errors/input/file with @#%^&!$ characters.mp4"
# Create a simple video file
ffmpeg -hide_banner -loglevel error -f lavfi -i testsrc=duration=1:size=640x480:rate=30 \
  -c:v libx264 -preset ultrafast -t 1 "$SPECIAL_CHAR_FILE"

./vidkit --preview --batch $NO_METADATA "$SPECIAL_CHAR_FILE" 2>&1 | tee output.log
if grep -q "panic\|crash" output.log; then
  echo "❌ Failed: Program crashed on filename with special characters"
  exit 1
else
  echo "✅ Passed: Gracefully handled filename with special characters"
fi

# Test 8: Very large or zero-length filename
echo -e "\n--- Test 8: Very large filename ---"
# Create a file with a very long name
LONG_NAME_FILE="test_results/errors/input/$(head /dev/urandom | tr -dc A-Za-z0-9 | head -c 200).mp4"
# Create a simple video file
ffmpeg -hide_banner -loglevel error -f lavfi -i testsrc=duration=1:size=640x480:rate=30 \
  -c:v libx264 -preset ultrafast -t 1 "$LONG_NAME_FILE"

./vidkit --preview --batch $NO_METADATA "$LONG_NAME_FILE" 2>&1 | tee output.log
if grep -q "panic\|crash" output.log; then
  echo "❌ Failed: Program crashed on very long filename"
  exit 1
else
  echo "✅ Passed: Gracefully handled very long filename"
fi

echo -e "\n=== Error Handling Tests Completed Successfully ==="
