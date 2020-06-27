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

func NewVideoFilterBuilder(width, height int) *VideoFilterBuilder {
	videoFilter := &VideoFilterOptions{}
	videoFilter.command = ""
	videoFilter.width = width
	videoFilter.height = height
	if width > 0 && height > 0 {
		videoFilter.command = videoFilter.command + fmt.Sprintf("-vf scale=%d:%d", width, height)
	} else {
		videoFilter.command = videoFilter.command + fmt.Sprintf("-vf scale=%d:%s", width, "trunc(ow/a/2)*2")
	}
	b := &VideoFilterBuilder{videoFilter: videoFilter}
	return b
}

func (b *VideoFilterBuilder) AudioCodec(audioCodec string) *VideoFilterBuilder {
	b.videoFilter.audioCodec = audioCodec
	b.videoFilter.command = b.videoFilter.command + fmt.Sprintf("%s -c:a %s", Separator, audioCodec)
	return b
}

func (b *VideoFilterBuilder) VideoCodec(videoCodec string) *VideoFilterBuilder {
	b.videoFilter.videoCodec = videoCodec
	b.videoFilter.command = b.videoFilter.command + fmt.Sprintf("%s -c:v %s", Separator, videoCodec)
	return b
}

func (b *VideoFilterBuilder) AudioSampleRate(audioSampleRate int) *VideoFilterBuilder {
	b.videoFilter.audioSampleRate = audioSampleRate
	b.videoFilter.command = b.videoFilter.command + fmt.Sprintf("%s -ar %d", Separator, audioSampleRate)
	return b
}

func (b *VideoFilterBuilder) VideoProfile(videoProfile string) *VideoFilterBuilder {
	b.videoFilter.videoProfile = videoProfile
	b.videoFilter.command = b.videoFilter.command + fmt.Sprintf("%s -profile:v %s", Separator, videoProfile)
	return b
}

func (b *VideoFilterBuilder) Compression(compression int) *VideoFilterBuilder {
	b.videoFilter.compression = compression
	b.videoFilter.command = b.videoFilter.command + fmt.Sprintf("%s -crf %d", Separator, compression)
	return b
}

func (b *VideoFilterBuilder) SCThreshold(threshold int) *VideoFilterBuilder {
	b.videoFilter.threshold = threshold
	b.videoFilter.command = b.videoFilter.command + fmt.Sprintf("%s -sc_threshold %d -g 48 -keyint_min 48", Separator, threshold)
	return b
}

func (b *VideoFilterBuilder) HlsTime(hlsTime int) *VideoFilterBuilder {
	b.videoFilter.hlstime = hlsTime
	b.videoFilter.command = b.videoFilter.command + fmt.Sprintf("%s -hls_time %d", Separator, hlsTime)
	return b
}

func (b *VideoFilterBuilder) HlsPlaylistType(hlsPlaylistType string) *VideoFilterBuilder {
	b.videoFilter.hlsPlaylistType = hlsPlaylistType
	b.videoFilter.command = b.videoFilter.command + fmt.Sprintf("%s -hls_playlist_type %s", Separator, hlsPlaylistType)
	return b
}

func (b *VideoFilterBuilder) VideoBitrate(videoBitrate int) *VideoFilterBuilder {
	b.videoFilter.videoBitrate = videoBitrate
	b.videoFilter.command = b.videoFilter.command + fmt.Sprintf("%s -b:v %dk", Separator, videoBitrate)
	return b
}

func (b *VideoFilterBuilder) Maxrate(maxrate int) *VideoFilterBuilder {
	b.videoFilter.maxrate = maxrate
	b.videoFilter.command = b.videoFilter.command + fmt.Sprintf("%s -maxrate %dk", Separator, maxrate)
	return b
}

func (b *VideoFilterBuilder) BufferSize(bufferSize int) *VideoFilterBuilder {
	b.videoFilter.bufferSize = bufferSize
	b.videoFilter.command = b.videoFilter.command + fmt.Sprintf("%s -bufsize %dk", Separator, bufferSize)
	return b
}

func (b *VideoFilterBuilder) AudioBitrate(audioBitrate int) *VideoFilterBuilder {
	b.videoFilter.audioBitrate = audioBitrate
	b.videoFilter.command = b.videoFilter.command + fmt.Sprintf("%s -b:a %dk", Separator, audioBitrate)
	return b
}

func (b *VideoFilterBuilder) Build() *VideoFilterOptions {
	return b.videoFilter
}

func (b *VideoFilterOptions) GetFilterCommand() string {
	return b.command
}
