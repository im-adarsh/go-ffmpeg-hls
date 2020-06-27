package command

import "fmt"

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

func (b *HLSStreamBuilder) Build() (string, error) {

	for _, v := range b.hLSStream.videoFilters {
		filePrefix := fmt.Sprintf("%s/%d_%d_%d", b.hLSStream.outputDirectoryPath, v.width, v.height, v.videoBitrate)
		segmentFileName := filePrefix + "_%03d.ts"
		b.hLSStream.command = b.hLSStream.command + Separator + v.GetFilterCommand()
		b.hLSStream.command = b.hLSStream.command + Separator + fmt.Sprintf("-hls_segment_filename %s", segmentFileName)
		b.hLSStream.command = b.hLSStream.command + Separator + fmt.Sprintf("%s.m3u8", filePrefix)
	}

	//b.hLSStream.command = b.hLSStream.command + Separator + fmt.Sprintf("-var_stream_map \"v:0,a:0 v:1,a:1\"")
	b.hLSStream.command = b.hLSStream.command + Separator + fmt.Sprintf("-hls_flags single_file")
	b.hLSStream.command = b.hLSStream.command + Separator + fmt.Sprintf("-master_pl_name %s/%s", b.hLSStream.outputDirectoryPath, b.hLSStream.masterFilename)

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
