package command

import "fmt"

// -vf scale=w=640:h=360:force_original_aspect_ratio=decrease
// -c:a aac -ar 48000 -c:v h264 -profile:v main -crf 20 -sc_threshold 0
// -g 48 -keyint_min 48 -hls_time 4 -hls_playlist_type vod
// -b:v 800k -maxrate 856k -bufsize 1200k -b:a 96k
// -hls_segment_filename beach/360p_%03d.ts beach/360p.m3u8

type VideoFilterOptions struct {
	command         string
	width           int
	height          int
	audioCodec      string
	videoCodec      string
	audioSampleRate int // hertz
	videoProfile    string
	compression     int
	threshold       int
	videoBitrate    int
	maxrate         int
	hlstime         int
	hlsPlaylistType string
	bufferSize      int
	audioBitrate    int
}

// VideoFilterOptions builder pattern code
type VideoFilterBuilder struct {
	videoFilter *VideoFilterOptions
}

func NewVideoFilterBuilder() *VideoFilterBuilder {
	videoFilter := &VideoFilterOptions{}
	videoFilter.command = ""
	b := &VideoFilterBuilder{videoFilter: videoFilter}
	return b
}

func (b *VideoFilterBuilder) Dimension(width, height int) *VideoFilterBuilder {
	b.videoFilter.width = width
	b.videoFilter.height = height
	if width > 0 && height > 0 {
		b.videoFilter.command = b.videoFilter.command + Separator + fmt.Sprintf("-vf scale=%d:%d", width, height)
	} else {
		b.videoFilter.command = b.videoFilter.command + Separator + fmt.Sprintf("-vf scale=%d:%s", width, "trunc(ow/a/2)*2")
	}
	return b
}

func (b *VideoFilterBuilder) AudioCodec(audioCodec string) *VideoFilterBuilder {
	b.videoFilter.audioCodec = audioCodec
	b.videoFilter.command = b.videoFilter.command + Separator + fmt.Sprintf("-c:a %s", audioCodec)
	return b
}

func (b *VideoFilterBuilder) VideoCodec(videoCodec string) *VideoFilterBuilder {
	b.videoFilter.videoCodec = videoCodec
	b.videoFilter.command = b.videoFilter.command + Separator + fmt.Sprintf("-c:v %s", videoCodec)
	return b
}

func (b *VideoFilterBuilder) AudioSampleRate(audioSampleRate int) *VideoFilterBuilder {
	b.videoFilter.audioSampleRate = audioSampleRate
	b.videoFilter.command = b.videoFilter.command + Separator + fmt.Sprintf("-ar %d", audioSampleRate)
	return b
}

func (b *VideoFilterBuilder) VideoProfile(videoProfile string) *VideoFilterBuilder {
	b.videoFilter.videoProfile = videoProfile
	b.videoFilter.command = b.videoFilter.command + Separator + fmt.Sprintf("-profile:v %s", videoProfile)
	return b
}

func (b *VideoFilterBuilder) Compression(compression int) *VideoFilterBuilder {
	b.videoFilter.compression = compression
	b.videoFilter.command = b.videoFilter.command + Separator + fmt.Sprintf("-crf %d", compression)
	return b
}

func (b *VideoFilterBuilder) SCThreshold(threshold int) *VideoFilterBuilder {
	b.videoFilter.threshold = threshold
	b.videoFilter.command = b.videoFilter.command + Separator + fmt.Sprintf("-sc_threshold %d -g 48 -keyint_min 48 -hls_playlist_type vod", threshold)
	return b
}

func (b *VideoFilterBuilder) HlsTime(hlsTime int) *VideoFilterBuilder {
	b.videoFilter.hlstime = hlsTime
	b.videoFilter.command = b.videoFilter.command + Separator + fmt.Sprintf("-hls_time %d", hlsTime)
	return b
}

func (b *VideoFilterBuilder) HlsPlaylistType(hlsPlaylistType string) *VideoFilterBuilder {
	b.videoFilter.hlsPlaylistType = hlsPlaylistType
	b.videoFilter.command = b.videoFilter.command + Separator + fmt.Sprintf("-hls_playlist_type %s", hlsPlaylistType)
	return b
}

func (b *VideoFilterBuilder) VideoBitrate(videoBitrate int) *VideoFilterBuilder {
	b.videoFilter.videoBitrate = videoBitrate
	b.videoFilter.command = b.videoFilter.command + Separator + fmt.Sprintf("-b:v %dk", videoBitrate)
	return b
}

func (b *VideoFilterBuilder) Maxrate(maxrate int) *VideoFilterBuilder {
	b.videoFilter.maxrate = maxrate
	b.videoFilter.command = b.videoFilter.command + Separator + fmt.Sprintf("-maxrate %dk", maxrate)
	return b
}

func (b *VideoFilterBuilder) BufferSize(bufferSize int) *VideoFilterBuilder {
	b.videoFilter.bufferSize = bufferSize
	b.videoFilter.command = b.videoFilter.command + Separator + fmt.Sprintf("-bufsize %dk", bufferSize)
	return b
}

func (b *VideoFilterBuilder) AudioBitrate(audioBitrate int) *VideoFilterBuilder {
	b.videoFilter.audioBitrate = audioBitrate
	b.videoFilter.command = b.videoFilter.command + Separator + fmt.Sprintf("-b:a %dk", audioBitrate)
	return b
}

func (b *VideoFilterBuilder) Build() *VideoFilterOptions {
	return b.videoFilter
}

func (b *VideoFilterOptions) GetFilterCommand() string {
	return b.command
}