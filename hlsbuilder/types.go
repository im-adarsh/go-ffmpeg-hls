package hlsbuilder

import (
	"fmt"
	"os"
	"path/filepath"
)

var initCommand = "ffmpeg  -y"

const Separator = " "

type hlsStream struct {
	command              string
	hideBanner           bool
	inputFilePath        string
	outputDirectoryPath  string
	masterFilename       string
	masterFileVideoCodec string // TODO detect this
	videoFilters         []VideoFilterOptions
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

func (b *HLSStreamBuilder) MasterFileVideoCodec(masterFileVideoCodec string) *HLSStreamBuilder {
	b.hLSStream.masterFileVideoCodec = masterFileVideoCodec
	return b
}

func (b *HLSStreamBuilder) Build() (string, error) {
	videoFilters := ""
	for _, v := range b.hLSStream.videoFilters {
		dimension := fmt.Sprintf("%d_%d_%d", v.width, v.height, v.videoBitrate)
		videoFilters = videoFilters + Separator + v.GetFilterCommand() +
			Separator + fmt.Sprintf("-hls_segment_filename %s/%s_%s.ts", b.hLSStream.outputDirectoryPath, dimension, "%03d") +
			Separator + fmt.Sprintf("%s/%s.m3u8", b.hLSStream.outputDirectoryPath, dimension)
	}
	b.hLSStream.command = b.hLSStream.command + Separator + videoFilters
	return b.hLSStream.command, nil
}

func (b *HLSStreamBuilder) GenerateMasterPlaylist() error {

	lines := []string{"#EXTM3U", "#EXT-X-VERSION:3"}
	f, err := os.Create(filepath.Join(b.hLSStream.outputDirectoryPath, b.hLSStream.masterFilename))
	if err != nil {
		fmt.Println(err)
		return err
	}
	for _, v := range b.hLSStream.videoFilters {
		meta := fmt.Sprintf("#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%dx%d,CODECS=\"%+v\"", v.videoBitrate*1000, v.width, v.height, b.hLSStream.masterFileVideoCodec)
		dimension := fmt.Sprintf("%d_%d_%d", v.width, v.height, v.videoBitrate)
		segmentMaster := fmt.Sprintf("%s.m3u8", dimension)
		lines = append(lines, meta, segmentMaster)
	}
	for _, v := range lines {
		_, err = fmt.Fprintln(f, v)
		if err != nil {
			return err
		}
	}
	err = f.Close()
	if err != nil {
		return err
	}
	return nil
}
