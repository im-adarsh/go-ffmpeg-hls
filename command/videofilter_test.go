package command

import (
	"fmt"
	"testing"
)

// -vf scale=w=640:h=360:force_original_aspect_ratio=decrease
// -c:a aac -ar 48000 -c:v h264 -profile:v main -crf 20 -sc_threshold 0
// -g 48 -keyint_min 48 -hls_time 4 -hls_playlist_type vod
// -b:v 800k -maxrate 856k -bufsize 1200k -b:a 96k
// -hls_segment_filename beach/360p_%03d.ts beach/360p.m3u8
func TestHLSTranscoder_GetCommand(t *testing.T) {

	vf := NewVideoFilterBuilder().
		Dimension(640, 360).
		AudioCodec("aac").
		AudioSampleRate(48000).
		VideoCodec("h264").
		VideoProfile("main").
		Compression(20).
		SCThreshold(0).
		HlsTime(4).
		HlsPlaylistType("vod").
		VideoBitrate(800).
		Maxrate(856).
		BufferSize(1200).
		AudioBitrate(96).
		Build()

	got := vf.GetFilterCommand()
	expected := "-vf scale=w=640:h=360:force_original_aspect_ratio=decrease " +
		"-c:a aac -ar 48000 -c:v h264 -profile:v main -crf 20 -sc_threshold 0 " +
		"-g 48 -keyint_min 48 -hls_time 4 -hls_playlist_type vod" +
		"-b:v 800k -maxrate 856k -bufsize 1200k -b:a 96k"

	if got != expected {
		fmt.Printf("expected : %s, got : %s", expected, got)
	}

}
