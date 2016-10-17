package imageio

import (
	"errors"
	"image"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"fmt"
)

// Mp4 video convert options.
type Options struct {
	FPS         int
	Codec       string
	Pixelformat string
	Pixfmt      string
}

// Video struct.
type Video struct {
	Cmd       *exec.Cmd
	StdinWr   io.WriteCloser
	Dimension string
	ExePath   string
	Output    string
	Option    *Options
}

// New Instance.
func NewVideo(output string, op *Options) *Video {
	if op.Codec == "" {
		op.Codec = "libx264"
	}
	if op.Pixelformat == "" {
		op.Pixelformat = "yuv420p"
	}
	if op.Pixfmt == "" {
		op.Pixfmt = "rgba"
	}
	if op.FPS == 0 {
		op.FPS = 25
	}
	return &Video{Cmd: nil, StdinWr: nil, Dimension: "", Output: output, Option: op}
}

// Initialize FFmpeg thread.
func (m *Video) initialize() error {
	exe, err := GetFFmpegExe()
	if err != nil {
		return err
	}
	if m.Dimension == "" {
		m.Dimension = "512x512"
	}
	pix_fmt := m.Option.Pixfmt
	fps := strconv.Itoa(m.Option.FPS)
	codec := m.Option.Codec
	pixelformat := m.Option.Pixelformat
	outputfile := m.Output

	cmdstr := []string{"-y",
		"-f", "rawvideo",
		"-vcodec", "rawvideo",
		"-s", m.Dimension,
		"-pix_fmt", pix_fmt,
		"-r", fps,
		"-i", "-", "-an",
		"-vcodec", codec,
		"-pix_fmt", pixelformat,
		"-crf", "25",
		"-r", "50",
		"-v", "warning", outputfile}
	return m.execFFmpegCommands(exe, cmdstr)
}

// Write image for mp4.
func (m *Video) WriteImage(imagePath string) error {
	img, err := LoadImage(imagePath)
	if err != nil {
		return err
	}
	width, height, err := m.getImageDimension(imagePath)
	if err != nil {
		return err
	}
	dimension := fmt.Sprintf("%dx%d", width, height)
	if m.Dimension == "" {
		m.Dimension = dimension
		if err := m.initialize(); err != nil {
			return err
		}
	}
	if dimension != m.Dimension {
		return errors.New("All images in a movie should have same size.")
	}
	if img != nil && m.Cmd != nil && m.StdinWr != nil {
		imgstring := LoadImageBitmap(img)
		m.StdinWr.Write(imgstring)
	}
	return nil
}

// Get image Dimension
func (m *Video) getImageDimension(imagePath string) (int, int, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return 0, 0, err
	}
	image, _, err := image.DecodeConfig(file)
	return image.Width, image.Height, err
}

// Close ffmpeg
func (m *Video) Close() error {
	if m.Cmd == nil {
		return errors.New("FFmpeg command is nil.")
	}
	if m.StdinWr != nil {
		err := m.StdinWr.Close()
		if err != nil {
			return err
		}
		log.Print("Close the mp4 video.")
	}
	m.Cmd.Wait()
	m.Cmd = nil
	return nil
}

// Execute FFmpeg commands.
func (m *Video) execFFmpegCommands(commandName string, params []string) error {
	m.Cmd = exec.Command(commandName, params...)
	stdinWr, err := m.Cmd.StdinPipe()
	if err != nil {
		return err
	}
	m.StdinWr = stdinWr
	m.Cmd.Start()
	return nil
}
