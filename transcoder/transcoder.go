package transcoder

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os/exec"

	"github.com/im-adarsh/go-ffmpeg-hls/hlsbuilder"
	"github.com/pkg/errors"
)

type hlsTranscoder struct {
	inputFile            string
	outputDir            string
	masterFileName       string
	masterFileVideoCodec string
	videoFilters         []hlsbuilder.VideoFilterOptions
}

// hlsTranscoder builder pattern code
type hlsTranscoderBuilder struct {
	hlsTranscoder *hlsTranscoder
}

func NewHlsTranscoderBuilder() *hlsTranscoderBuilder {
	hlsTranscoder := &hlsTranscoder{}
	b := &hlsTranscoderBuilder{hlsTranscoder: hlsTranscoder}
	return b
}

func (b *hlsTranscoderBuilder) VideoFiltersOptions(videoFilters []hlsbuilder.VideoFilterOptions) *hlsTranscoderBuilder {
	b.hlsTranscoder.videoFilters = videoFilters
	return b
}

func (b *hlsTranscoderBuilder) MasterFileName(masterFileName string) *hlsTranscoderBuilder {
	b.hlsTranscoder.masterFileName = masterFileName
	return b
}

func (b *hlsTranscoderBuilder) MasterFileVideoCodec(masterFileVideoCodec string) *hlsTranscoderBuilder {
	b.hlsTranscoder.masterFileVideoCodec = masterFileVideoCodec
	return b
}

func (b *hlsTranscoderBuilder) InputFile(inputFile string) *hlsTranscoderBuilder {
	b.hlsTranscoder.inputFile = inputFile
	return b
}

func (b *hlsTranscoderBuilder) OutputDir(outputDir string) *hlsTranscoderBuilder {
	b.hlsTranscoder.outputDir = outputDir
	return b
}

func (b *hlsTranscoderBuilder) Run() (*hlsTranscoder, error) {

	builder := hlsbuilder.
		NewHLSStreamBuilder(b.hlsTranscoder.inputFile, b.hlsTranscoder.outputDir).
		HideBanner(true).
		MasterFileName("master.m3u8").
		MasterFileVideoCodec(b.hlsTranscoder.masterFileVideoCodec)

	if b.hlsTranscoder.masterFileName != "" {
		builder.MasterFileName(b.hlsTranscoder.masterFileName)
	}

	for _, v := range b.hlsTranscoder.videoFilters {
		builder.AppendVideoFilter(v)
	}

	cmdFfmpeg, err := builder.Build()
	if err != nil {
		return nil, errors.Wrap(FailedToGenerateCommand, "unable to prepare command")
	}
	log.Println(context.Background(), "ffmpeg", cmdFfmpeg, "input_file", b.hlsTranscoder.inputFile)

	cmd := exec.Command("bash", "-c", cmdFfmpeg)
	// create a pipe for the output of the script
	cmdReader, err := cmd.StderrPipe()
	if err != nil {
		return nil, errors.Wrap(FailedInitializeStdPipe, "")
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("\t > %s\n", scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		return nil, errors.Wrap(FailedToStartCommand, fmt.Sprintf("error starting cmd : %+v", err))
	}

	err = cmd.Wait()
	if err != nil {
		return nil, errors.Wrap(FailedToWaitCommand, fmt.Sprintf("error waiting cmd : %+v", err))
	}

	err = builder.GenerateMasterPlaylist()
	if err != nil {
		return nil, errors.Wrap(FailedToGenerateMasterFile, fmt.Sprintf("error generating master file : %+v", err))
	}

	return b.hlsTranscoder, nil
}
