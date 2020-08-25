package hlsbuilder

import (
	"fmt"
	"testing"
)

// ffmpeg  -y -i apple.mov -hide_banner -vf scale=640:360 -c:a aac -ar 48000 -c:v h264 -profile:v main -crf 20 -sc_threshold 0 -g 48 -keyint_min 48 -hls_time 4 -hls_playlist_type vod -b:v 800k -maxrate 856k -bufsize 1200k -b:a 96k -master_pl_name master.m3u8
func TestHLSTranscoder_Build(t *testing.T) {

	vf2 := getVideoFilter(100, 100, 0)
	vf1 := getVideoFilter(640, 360, 1)
	cmd, err := NewHLSStreamBuilder("input.mp4", ".").
		HideBanner(true).
		AppendVideoFilter(vf1).
		AppendVideoFilter(vf2).
		MasterFileName("master.m3u8").
		Build()
	if err != nil {
		t.Error()
	}

	fmt.Println(cmd)
	expected := "ffmpeg  -y -i input.mp4 -hide_banner  -vf scale=640:360 -c:a aac -ar 48000 -c:v h264 -profile:v main -crf 20 -sc_threshold 0 -g 48 -keyint_min 48 -hls_time 4 -hls_playlist_type vod -b:v 800k -maxrate 856k -bufsize 1200k -b:a 96k -hls_segment_filename ./640_360_800_%03d.ts ./640_360_800.m3u8 -vf scale=100:100 -c:a aac -ar 48000 -c:v h264 -profile:v main -crf 20 -sc_threshold 0 -g 48 -keyint_min 48 -hls_time 4 -hls_playlist_type vod -b:v 800k -maxrate 856k -bufsize 1200k -b:a 96k -hls_segment_filename ./100_100_800_%03d.ts ./100_100_800.m3u8"
	if cmd != expected {
		t.Errorf("expected : %s, got : %s", expected, cmd)
	}

}
