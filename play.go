package main

import (
	"fmt"

	"github.com/im-adarsh/go-ffmpeg-hls/transcoder"

	"github.com/im-adarsh/go-ffmpeg-hls/hlsbuilder"
)

func main() {

	vf2 := getVideoFilter(100, -1, 0)
	vf1 := getVideoFilter(640, -1, 1)

	t, err := transcoder.NewHlsTranscoderBuilder().
		InputFile("sample_input/input.mov").
		OutputDir("./output").
		VideoFiltersOptions([]hlsbuilder.VideoFilterOptions{vf1, vf2}).
		MasterFileName("master99.m3u8").
		Run()
	if err != nil {
		fmt.Println(err, "")
		return
	}
	_ = t
}

func getVideoFilter(width, height int, filterIndex int) hlsbuilder.VideoFilterOptions {
	vf := hlsbuilder.NewVideoFilterBuilder(width, height, filterIndex).
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
	return *vf
}
