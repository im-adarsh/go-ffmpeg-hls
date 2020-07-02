package main

import (
	"fmt"

	"./command"
	"./hls"
)

func main() {

	vf2 := getVideoFilter(100, 100, 0)
	vf1 := getVideoFilter(640, 360, 1)
	builder := command.NewHLSStreamBuilder("sample_input/input.mov", "./output").
		HideBanner(true).
		AppendVideoFilter(vf1).
		AppendVideoFilter(vf2).
		MasterFileName("master.m3u8")
	cmd, err := builder.Build()
	if err != nil {
		fmt.Println(err)
	}

	tran := new(hls.HLSTranscoder)
	tran.NewHlsTranscoder(cmd)
	done := tran.Run(false)

	// This channel is used to wait for the process to end
	err = <-done
	if err != nil {
		fmt.Println(err)
	}

	builder.GenerateMasterPlaylist()

}

func getVideoFilter(width, height int, filterIndex int) command.VideoFilterOptions {
	vf := command.NewVideoFilterBuilder(width, height, filterIndex).
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
