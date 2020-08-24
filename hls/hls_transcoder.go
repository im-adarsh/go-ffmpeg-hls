package hls

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/im-adarsh/go-ffmpeg-hls/ffmpeg"
	"io"
	"os/exec"
)

type HLSTranscoder struct {
	command       string
	configuration ffmpeg.Configuration
	process       *exec.Cmd
	stdErrPipe    io.ReadCloser
	stdStdinPipe  io.WriteCloser
}

func (t *HLSTranscoder) NewHlsTranscoder(command string) error {
	var err error

	if command == "" {
		return errors.New("error on transcoder.Initialize: command missing")
	}
	t.command = command

	cfg := t.configuration
	if len(cfg.FfmpegBin) == 0 || len(cfg.FfprobeBin) == 0 {
		cfg, err = ffmpeg.Configure()
		if err != nil {
			return err
		}
	}
	t.configuration = cfg

	return nil
}

// Run Starts the transcoding process
func (t *HLSTranscoder) Run(progress bool) <-chan error {
	done := make(chan error)
	command := t.command

	proc := exec.Command(t.configuration.FfmpegBin, command)
	if progress {
		errStream, err := proc.StderrPipe()
		if err != nil {
			fmt.Println("Progress not available: " + err.Error())
		} else {
			t.stdErrPipe = errStream
		}
	}

	stdin, err := proc.StdinPipe()
	if nil != err {
		fmt.Println("Stdin not available: " + err.Error())
	}

	t.stdStdinPipe = stdin

	out := &bytes.Buffer{}
	if progress {
		proc.Stdout = out
	}

	err = proc.Start()

	t.SetProcess(proc)
	go func(err error, out *bytes.Buffer) {
		if err != nil {
			done <- fmt.Errorf("Failed Start FFMPEG (%s) with %s, message %s", command, err, out.String())
			close(done)
			return
		}
		err = proc.Wait()
		if err != nil {
			err = fmt.Errorf("Failed Finish FFMPEG (%s) with %s message %s", command, err, out.String())
		}
		done <- err
		close(done)
	}(err, out)

	return done
}

// SetProcess Set the transcoding process
func (t *HLSTranscoder) SetProcess(cmd *exec.Cmd) {
	t.process = cmd
}
