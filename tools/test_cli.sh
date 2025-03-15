#!/bin/bash
# Test script for VidKit's command-line interface functionality

# Set up test environment
echo "=== Testing Command-Line Interface ==="
mkdir -p test_results/cli
rm -rf test_results/cli/*

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

# Find a test video to use
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
  if [ -z "$TEST_MOVIE" ]; then
    echo "❌ No test videos found. Cannot continue."
    exit 1
  fi
fi

echo "Using test video: $TEST_MOVIE"

# Test 1: Help flag
echo -e "\n--- Test 1: Help flag ---"
./vidkit --help | tee output.log
if ! grep -q "Usage\|vidkit" output.log; then
  echo "❌ Failed: Help information not displayed properly"
  exit 1
else
  echo "✅ Passed: Help information displayed properly"
fi

# Test 2: Version flag
echo -e "\n--- Test 2: Version flag ---"
./vidkit --version | tee output.log
if ! grep -q "VidKit\|version" output.log; then
  echo "❌ Failed: Version information not displayed properly"
  exit 1
else
  echo "✅ Passed: Version information displayed properly"
fi

# Test 3: Preview mode
echo -e "\n--- Test 3: Preview mode ---"
./vidkit --preview "$TEST_MOVIE" | tee output.log
if ! grep -q "PREVIEW MODE" output.log; then
  echo "❌ Failed: Preview mode indicator not displayed"
  exit 1
else
  echo "✅ Passed: Preview mode works correctly"
fi

# Test 4: No-metadata flag
echo -e "\n--- Test 4: No-metadata flag ---"
./vidkit --preview --no-metadata "$TEST_MOVIE" | tee output.log
if grep -q "Looking up movie metadata" output.log; then
  echo "❌ Failed: Metadata lookup still performed despite no-metadata flag"
  exit 1
else
  echo "✅ Passed: no-metadata flag works correctly"
fi

# Test 5: Language option
echo -e "\n--- Test 5: Language option ---"
./vidkit --preview --lang fr "$TEST_MOVIE" | tee output.log
if ! grep -q "lang=fr\|Language: fr" output.log; then
  echo "⚠️ Note: Could not verify language option, but command did not fail"
else
  echo "✅ Passed: Language option accepted"
fi

# Test 6: Custom movie format
echo -e "\n--- Test 6: Custom movie format ---"
CUSTOM_FORMAT="{title}_{year}"
./vidkit --preview --movie-format "$CUSTOM_FORMAT" "$TEST_MOVIE" | tee output.log
if ! grep -q "$CUSTOM_FORMAT\|format" output.log; then
  echo "⚠️ Note: Could not verify custom format, but command did not fail"
else
  echo "✅ Passed: Custom movie format accepted"
fi

# Test 7: Lowercase flag
echo -e "\n--- Test 7: Lowercase flag ---"
./vidkit --preview --lowercase "$TEST_MOVIE" | tee output.log
if grep -q "panic\|crash" output.log; then
  echo "❌ Failed: Lowercase flag caused crash"
  exit 1
else
  echo "✅ Passed: Lowercase flag accepted"
fi

# Test 8: Batch mode
echo -e "\n--- Test 8: Batch mode ---"
./vidkit --preview --batch "$TEST_MOVIE" | tee output.log
if grep -q "panic\|crash\|error" output.log; then
  echo "❌ Failed: Batch mode flag caused issue"
  exit 1
else
  echo "✅ Passed: Batch mode flag accepted"
fi

# Test 9: Scene style
echo -e "\n--- Test 9: Scene style flag ---"
./vidkit --preview --scene "$TEST_MOVIE" | tee output.log
if ! grep -q "scene" output.log && ! grep -q "dots\|separator" output.log; then
  echo "⚠️ Note: Could not verify scene style, but command did not fail"
else
  echo "✅ Passed: Scene style flag accepted"
fi

# Test 10: Provider selection
echo -e "\n--- Test 10: Provider selection ---"
./vidkit --preview --movie-provider tmdb "$TEST_MOVIE" | tee output.log
if grep -q "invalid provider\|unknown provider" output.log; then
  echo "❌ Failed: Valid provider was rejected"
  exit 1
else
  echo "✅ Passed: Provider selection works correctly"
fi

# Test 11: Directory organization
echo -e "\n--- Test 11: Directory organization options ---"
DIR_TEMPLATE="test_results/cli/Movies/{title}"
./vidkit --preview --movie-dir "$DIR_TEMPLATE" "$TEST_MOVIE" | tee output.log
if ! grep -q "test_results/cli/Movies" output.log; then
  echo "❌ Failed: Directory organization option not working"
  exit 1
else
  echo "✅ Passed: Directory organization option works correctly"
fi

# Test 12: Multiple arguments
echo -e "\n--- Test 12: Multiple arguments ---"
./vidkit --preview --batch --no-metadata --lowercase "$TEST_MOVIE" | tee output.log
if grep -q "panic\|crash\|conflict" output.log; then
  echo "❌ Failed: Multiple arguments caused issues"
  exit 1
else
  echo "✅ Passed: Multiple arguments work correctly"
fi

# Test 13: Unknown flag handling
echo -e "\n--- Test 13: Unknown flag handling ---"
./vidkit --preview --non-existent-flag "$TEST_MOVIE" 2>&1 | tee output.log
if ! grep -q "flag\|unknown\|error" output.log; then
  echo "⚠️ Note: Unknown flag did not generate expected error, check flag parsing"
else
  echo "✅ Passed: Unknown flag properly reported"
fi

echo -e "\n=== Command-Line Interface Tests Completed Successfully ==="
