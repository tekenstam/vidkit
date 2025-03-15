#!/bin/bash
# Script to generate test TV show files for VidKit testing
# This creates small, valid MP4 files with various naming conventions

# Create a test directory
mkdir -p test_videos

echo "Generating test TV show files..."

# Create a very small test video template (1 second)
ffmpeg -hide_banner -loglevel error -f lavfi -i testsrc=duration=1:size=320x240:rate=15 -c:v libx264 -preset ultrafast -t 1 "template.mp4"

# Function to create a test video by copying the template
create_test_video() {
  outfile="$1"
  cp template.mp4 "$outfile"
  echo "Created: $outfile"
}

# Clean up any existing files
rm -f test_videos/*.mp4

# Generate test files for each naming convention pattern

# Breaking Bad (2008)
create_test_video "test_videos/Breaking.Bad.S01E01.Pilot.mp4"
create_test_video "test_videos/Breaking.Bad.S01E02.Cats.in.the.Bag.mp4"
create_test_video "test_videos/Breaking.Bad.S01E03.And.the.Bags.in.the.River.mp4"
create_test_video "test_videos/Breaking.Bad.S01E04.Cancer.Man.mp4"
create_test_video "test_videos/Breaking.Bad.S01E05.Gray.Matter.mp4"
create_test_video "test_videos/Breaking.Bad.1x01.mp4"
create_test_video "test_videos/Breaking.Bad.Season.1.Episode.1.mp4"
create_test_video "test_videos/Breaking Bad (2008) S01E01.mp4"
create_test_video "test_videos/Breaking Bad S01E01 720p HEVC.mp4"
create_test_video "test_videos/Breaking.Bad.S01E01.1080p.WEB-DL.x264.mp4"

# Game of Thrones (2011)
create_test_video "test_videos/Game.of.Thrones.S01E01.Winter.Is.Coming.mp4"
create_test_video "test_videos/Game.of.Thrones.S01E02.The.Kingsroad.mp4"
create_test_video "test_videos/Game.of.Thrones.1x01.mp4"
create_test_video "test_videos/Game of Thrones (2011) S01E01.mp4"
create_test_video "test_videos/Game.of.Thrones.S01E01.1080p.BluRay.x264.mp4"

# Stranger Things (2016)
create_test_video "test_videos/Stranger.Things.S01E01.The.Vanishing.of.Will.Byers.mp4"
create_test_video "test_videos/Stranger.Things.S01E02.The.Weirdo.on.Maple.Street.mp4"
create_test_video "test_videos/Stranger Things (2016) S01E01.mp4"
create_test_video "test_videos/Stranger.Things.1x01.mp4"
create_test_video "test_videos/Stranger.Things.S01E01.2160p.NF.WEB-DL.mp4"

# The Office (2005)
create_test_video "test_videos/The.Office.S01E01.Pilot.mp4"
create_test_video "test_videos/The.Office.S01E02.Diversity.Day.mp4"
create_test_video "test_videos/The Office (2005) S01E01.mp4"
create_test_video "test_videos/The.Office.S01E01.720p.HDTV.x264.mp4"

# Friends (1994)
create_test_video "test_videos/Friends.S01E01.The.One.Where.Monica.Gets.a.Roommate.mp4"
create_test_video "test_videos/Friends.S01E02.The.One.with.the.Sonogram.at.the.End.mp4"
create_test_video "test_videos/Friends (1994) S01E01.mp4"
create_test_video "test_videos/Friends.S01E01.DVDRip.x264.mp4"

# The Office (US) (2005)
create_test_video "test_videos/The.Office.US.S01E01.Pilot.mp4"
create_test_video "test_videos/The.Office.S01E02.Diversity.Day.mp4"
create_test_video "test_videos/The Office (US) S01E01.mp4"
create_test_video "test_videos/The.Office.US.1x01.mp4"
create_test_video "test_videos/The.Office.S01E01.720p.NBC.WEB-DL.mp4"

# The Good Place (2016)
create_test_video "test_videos/The.Good.Place.S01E01.Everything.Is.Fine.mp4"
create_test_video "test_videos/The.Good.Place.S01E02.Flying.mp4"
create_test_video "test_videos/The Good Place (2016) - Season 1 Episode 1 - Everything Is Fine.mp4"
create_test_video "test_videos/The.Good.Place.1x01.mp4"
create_test_video "test_videos/The.Good.Place.S01E01.1080p.NBC.WEB-DL.mp4"

# Files with just show names (no episode info)
create_test_video "test_videos/Breaking Bad.mp4"
create_test_video "test_videos/Game of Thrones.mp4"
create_test_video "test_videos/Stranger Things.mp4"
create_test_video "test_videos/The Office.mp4"
create_test_video "test_videos/Friends.mp4"

# Special case formats
create_test_video "test_videos/[GROUP] Breaking Bad S01E05.mp4"
create_test_video "test_videos/Breaking_Bad_-_s01e01_-_Pilot.mp4"
create_test_video "test_videos/The Good Place (2016) - Season 1 Episode 1 - Everything Is Fine.mp4"
create_test_video "test_videos/Loki.S01E01.Glorious.Purpose.1080p.DSNP.WEB-DL.DDP5.1.Atmos.H.264-CMRG.mp4"
create_test_video "test_videos/Severance 1x01 Good News About Hell 2160p.mp4"

echo "Generating test movie files..."

# The Shawshank Redemption (1994)
create_test_video "test_videos/The.Shawshank.Redemption.1994.1080p.BluRay.x264.mp4"
create_test_video "test_videos/The Shawshank Redemption (1994).mp4"
create_test_video "test_videos/Shawshank.Redemption.1994.720p.mp4"

# Inception (2010)
create_test_video "test_videos/Inception.2010.1080p.BluRay.x264.mp4"
create_test_video "test_videos/Inception (2010).mp4"
create_test_video "test_videos/Inception.2010.720p.mp4"

# The Dark Knight (2008)
create_test_video "test_videos/The.Dark.Knight.2008.1080p.BluRay.x264.mp4"
create_test_video "test_videos/The Dark Knight (2008).mp4"
create_test_video "test_videos/The.Dark.Knight.2008.4K.HDR.mp4"

# Pulp Fiction (1994)
create_test_video "test_videos/Pulp.Fiction.1994.1080p.BluRay.x264.mp4"
create_test_video "test_videos/Pulp Fiction (1994).mp4"
create_test_video "test_videos/Pulp.Fiction.1994.720p.mp4"

# Interstellar (2014)
create_test_video "test_videos/Interstellar.2014.1080p.BluRay.x264.mp4"
create_test_video "test_videos/Interstellar (2014).mp4"
create_test_video "test_videos/Interstellar.2014.2160p.HDR.mp4"

echo "Test video generation complete!"

# Clean up template
rm template.mp4

echo "Done! Generated $(find test_videos -name "*.mp4" | wc -l) test files in test_videos directory."
echo "You can test them with: go run cmd/vidkit/main.go test_videos/<filename>"
