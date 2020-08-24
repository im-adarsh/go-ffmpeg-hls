package main

import (
	"fmt"
	"os/exec"

	"github.com/im-adarsh/go-ffmpeg-hls/hlsbuilder"
)

func main() {

	vf2 := getVideoFilter(100, -1, 0)
	vf1 := getVideoFilter(640, -1, 1)
	builder := hlsbuilder.NewHLSStreamBuilder("sample_input/input.mov", "./output").
		HideBanner(true).
		AppendVideoFilter(vf1).
		AppendVideoFilter(vf2).
		MasterFileName("master.m3u8")
	cmdFfmpeg, err := builder.Build()
	if err != nil {
		fmt.Println(err)
		return
	}

	cmd := exec.Command("bash", "-c", cmdFfmpeg)
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	builder.GenerateMasterPlaylist()

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
