#!/bin/bash
# Test script for VidKit's resolution and video quality detection features

# Check for --no-metadata flag
NO_METADATA=""
if [[ "$*" == *"--no-metadata"* ]]; then
  NO_METADATA="--no-metadata"
  echo "Running in no-metadata mode (offline)"
fi

# Set up test environment
echo "=== Testing Resolution and Video Quality Detection ==="
mkdir -p test_results/resolution
rm -rf test_results/resolution/*

# Build VidKit if needed
if [ ! -f "./vidkit" ]; then
  echo "Building VidKit..."
  go build -o vidkit ./cmd/vidkit
fi

# Ensure we have test videos
if [ ! -d "test_videos" ] || [ -z "$(ls -A test_videos)" ]; then
  echo "Generating test videos..."
  ./tools/generate_test_videos.sh
fi

# Test different resolutions - create videos with specific resolutions
echo -e "\n--- Creating test videos with different resolutions ---"

# Create a 480p test video
echo "Creating 480p test video..."
ffmpeg -hide_banner -loglevel error -f lavfi -i testsrc=duration=1:size=720x480:rate=30 \
  -c:v libx264 -preset ultrafast -t 1 "test_results/resolution/480p_video.mp4"

# Create a 720p test video
echo "Creating 720p test video..."
ffmpeg -hide_banner -loglevel error -f lavfi -i testsrc=duration=1:size=1280x720:rate=30 \
  -c:v libx264 -preset ultrafast -t 1 "test_results/resolution/720p_video.mp4"

# Create a 1080p test video
echo "Creating 1080p test video..."
ffmpeg -hide_banner -loglevel error -f lavfi -i testsrc=duration=1:size=1920x1080:rate=30 \
  -c:v libx264 -preset ultrafast -t 1 "test_results/resolution/1080p_video.mp4"

# Create a video with non-standard resolution
echo "Creating non-standard resolution video..."
ffmpeg -hide_banner -loglevel error -f lavfi -i testsrc=duration=1:size=1440x900:rate=30 \
  -c:v libx264 -preset ultrafast -t 1 "test_results/resolution/custom_resolution.mp4"

# Create videos with different codecs
echo -e "\n--- Creating test videos with different codecs ---"

# Create H.264 video
echo "Creating H.264 video..."
ffmpeg -hide_banner -loglevel error -f lavfi -i testsrc=duration=1:size=1280x720:rate=30 \
  -c:v libx264 -preset ultrafast -t 1 "test_results/resolution/h264_video.mp4"

# Create MPEG4 video if available
echo "Creating MPEG4 video..."
ffmpeg -hide_banner -loglevel error -f lavfi -i testsrc=duration=1:size=1280x720:rate=30 \
  -c:v mpeg4 -q:v 2 -t 1 "test_results/resolution/mpeg4_video.mp4"

# Create videos with different frame rates
echo -e "\n--- Creating test videos with different frame rates ---"

# Create 24fps video
echo "Creating 24fps video..."
ffmpeg -hide_banner -loglevel error -f lavfi -i testsrc=duration=1:size=1280x720:rate=24 \
  -c:v libx264 -preset ultrafast -t 1 "test_results/resolution/24fps_video.mp4"

# Create 30fps video
echo "Creating 30fps video..."
ffmpeg -hide_banner -loglevel error -f lavfi -i testsrc=duration=1:size=1280x720:rate=30 \
  -c:v libx264 -preset ultrafast -t 1 "test_results/resolution/30fps_video.mp4"

# Create 60fps video
echo "Creating 60fps video..."
ffmpeg -hide_banner -loglevel error -f lavfi -i testsrc=duration=1:size=1280x720:rate=60 \
  -c:v libx264 -preset ultrafast -t 1 "test_results/resolution/60fps_video.mp4"

# Count test videos
NUM_VIDEOS=$(ls -l test_results/resolution/*.mp4 | wc -l)
echo -e "\nCreated $NUM_VIDEOS test videos for resolution and quality testing"

# Run tests for resolution detection
echo -e "\n--- Testing resolution detection ---"

# Test 480p detection
echo "Testing 480p detection..."
./vidkit --preview $NO_METADATA test_results/resolution/480p_video.mp4 | tee output.log
if ! grep -q "Resolution: 480p" output.log; then
  echo "❌ Failed to detect 480p resolution correctly"
  exit 1
else
  echo "✅ 480p resolution detected correctly"
fi

# Test 720p detection
echo -e "\nTesting 720p detection..."
./vidkit --preview $NO_METADATA test_results/resolution/720p_video.mp4 | tee output.log
if ! grep -q "Resolution: 720p" output.log; then
  echo "❌ Failed to detect 720p resolution correctly"
  exit 1
else
  echo "✅ 720p resolution detected correctly"
fi

# Test 1080p detection
echo -e "\nTesting 1080p detection..."
./vidkit --preview $NO_METADATA test_results/resolution/1080p_video.mp4 | tee output.log
if ! grep -q "Resolution: 1080p" output.log; then
  echo "❌ Failed to detect 1080p resolution correctly"
  exit 1
else
  echo "✅ 1080p resolution detected correctly"
fi

# Test custom resolution detection
echo -e "\nTesting custom resolution detection..."
./vidkit --preview $NO_METADATA test_results/resolution/custom_resolution.mp4 | tee output.log
if ! grep -q "Resolution:" output.log; then
  echo "❌ Failed to detect custom resolution"
  exit 1
else
  echo "✅ Custom resolution detected"
fi

# Test codec detection
echo -e "\n--- Testing codec detection ---"

# Test H.264 detection
echo "Testing H.264 detection..."
./vidkit --preview $NO_METADATA test_results/resolution/h264_video.mp4 | tee output.log
if ! grep -q "Codec: h264" output.log; then
  echo "❌ Failed to detect H.264 codec correctly"
  exit 1
else
  echo "✅ H.264 codec detected correctly"
fi

# Test frame rate detection
echo -e "\n--- Testing frame rate detection ---"

# Test 24fps detection
echo "Testing 24fps detection..."
./vidkit --preview $NO_METADATA test_results/resolution/24fps_video.mp4 | tee output.log
if ! grep -q "Frame Rate: 24.00 fps" output.log; then
  echo "❌ Failed to detect 24fps correctly"
  exit 1
else
  echo "✅ 24fps detected correctly"
fi

# Test 60fps detection
echo -e "\nTesting 60fps detection..."
./vidkit --preview $NO_METADATA test_results/resolution/60fps_video.mp4 | tee output.log
if ! grep -q "Frame Rate: 60.00 fps" output.log; then
  echo "❌ Failed to detect 60fps correctly"
  exit 1
else
  echo "✅ 60fps detected correctly"
fi

echo -e "\n=== Resolution and Video Quality Tests Completed Successfully ==="
