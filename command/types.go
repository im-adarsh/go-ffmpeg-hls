package command

import (
	"fmt"
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
	videoFilters := ""
	for _, v := range b.hLSStream.videoFilters {
		dimension := fmt.Sprintf("%d_%d_%d", v.width, v.height, v.videoBitrate)
		videoFilters = videoFilters + Separator + v.GetFilterCommand() +
			Separator + fmt.Sprintf("-hls_segment_filename %s/%s_%s.ts", b.hLSStream.outputDirectoryPath, dimension, "%03d") +
			Separator + fmt.Sprintf("%s/%s.m3u8", b.hLSStream.outputDirectoryPath, dimension)
	}
	b.hLSStream.command = b.hLSStream.command + Separator + videoFilters

	//segmentPrefix := filepath.Join(b.hLSStream.outputDirectoryPath, "v%v/fileSequence%d.ts")
	//segmentMasterPrefix := filepath.Join(b.hLSStream.outputDirectoryPath, "v%v/prog_index.m3u8")
	//b.hLSStream.command = b.hLSStream.command + Separator + fmt.Sprintf("-hls_segment_filename %s", fmt.Sprintf("\"%s\"", segmentPrefix))
	//b.hLSStream.command = b.hLSStream.command + Separator + segmentMasterPrefix

	//b.hLSStream.command = b.hLSStream.command + Separator + fmt.Sprintf("-hls_flags single_file")

	return b.hLSStream.command, nil
}
