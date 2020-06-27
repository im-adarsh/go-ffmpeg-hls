package command

import (
	"fmt"
	"path/filepath"
)

var initCommand = "ffmpeg  -y"

const Separator = " "

type hlsStream struct {
	command             string
	hideBanner          bool
	inputFilePath       string
	outputDirectoryPath string
	masterFilename      string
	videoFilters        []VideoFilterOptions
}

// hlsStream builder pattern code
type HLSStreamBuilder struct {
	hLSStream *hlsStream
}

func NewHLSStreamBuilder(inputFilePath string, outputDirectoryPath string) *HLSStreamBuilder {
	hLSStream := &hlsStream{command: initCommand}
	hLSStream.inputFilePath = inputFilePath
	hLSStream.outputDirectoryPath = outputDirectoryPath
	hLSStream.videoFilters = []VideoFilterOptions{}
	hLSStream.command = hLSStream.command + Separator + fmt.Sprintf("-i %s", inputFilePath)
	b := &HLSStreamBuilder{hLSStream: hLSStream}
	return b
}

func (b *HLSStreamBuilder) HideBanner(hideBanner bool) *HLSStreamBuilder {
	b.hLSStream.hideBanner = hideBanner
	if hideBanner {
		b.hLSStream.command = b.hLSStream.command + Separator + "-hide_banner"
	}
	return b
}

func (b *HLSStreamBuilder) AppendVideoFilter(vf VideoFilterOptions) *HLSStreamBuilder {
	b.hLSStream.videoFilters = append(b.hLSStream.videoFilters, vf)
	return b
}

func (b *HLSStreamBuilder) AppendOption(key, value string) *HLSStreamBuilder {
	b.hLSStream.command = b.hLSStream.command + Separator + fmt.Sprintf("-%s %s", key, value)
	return b
}

func (b *HLSStreamBuilder) MasterFileName(masterFileName string) *HLSStreamBuilder {
	b.hLSStream.masterFilename = masterFileName
	return b
}

// -map 0:v -map 0:a -map 0:v -map 0:a -f hls -var_stream_map "v:0,a:0 v:1,a:1"
func (b *HLSStreamBuilder) Build() (string, error) {

	maps := ""
	varStreamMaps := ""
	videoFilters := ""
	for _, v := range b.hLSStream.videoFilters {
		maps = maps + "-map 0" + Separator
		varStreamMaps = varStreamMaps + fmt.Sprintf("v:%d,a:%d", v.filterIndex, v.filterIndex) + Separator
		videoFilters = videoFilters + Separator + v.GetFilterCommand()
	}
	b.hLSStream.command = b.hLSStream.command + Separator + maps + videoFilters
	b.hLSStream.command = b.hLSStream.command + Separator + fmt.Sprintf("-var_stream_map \"%s\"", varStreamMaps)
	b.hLSStream.command = b.hLSStream.command + Separator + fmt.Sprintf("-master_pl_name %s", filepath.Join(b.hLSStream.outputDirectoryPath, b.hLSStream.masterFilename))

	segmentPrefix := filepath.Join(b.hLSStream.outputDirectoryPath, "v%v/fileSequence%d.ts")
	segmentMasterPrefix := filepath.Join(b.hLSStream.outputDirectoryPath, "v%v/prog_index.m3u8")
	b.hLSStream.command = b.hLSStream.command + Separator + fmt.Sprintf("-hls_segment_filename %s", fmt.Sprintf("\"%s\"", segmentPrefix))
	b.hLSStream.command = b.hLSStream.command + Separator + segmentMasterPrefix

	//b.hLSStream.command = b.hLSStream.command + Separator + fmt.Sprintf("-hls_flags single_file")

	return b.hLSStream.command, nil
}

/*
ffmpeg -threads 0 -i input.mov -r 24 -g 48 -keyint_min 48 -sc_threshold 0 -c:v libx264^
-s:v:0 960x540 -b:v:0 2400k -maxrate:v:0 2640k -bufsize:v:0 2400k^
-s:v:1 1920x1080 -b:v:1 5200k -maxrate:v:1 5720k -bufsize:v:1 5200k^
-s:v:2 1280x720 -b:v:2 3100k -maxrate:v:2 3410k -bufsize:v:2 3100k^
-s:v:3 640x360 -b:v:3 1200k -maxrate:v:3 1320k -bufsize:v:3 1200k^
-b:a 128k -ar 44100 -ac 2^
-map 0:v -map 0:v -map 0:v -map 0:v -map 0:a^
-f hls -var_stream_map "v:0,agroup:audio v:1,agroup:audio v:2,agroup:audio v:3,agroup:audio a:0,agroup:audio"^
-hls_flags single_file -hls_segment_type fmp4 -hls_list_size 0 -hls_time 6  -master_pl_name master.m3u8 -y TOS%v.m3u8

*/
