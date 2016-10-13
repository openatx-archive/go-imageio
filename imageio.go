package imageio

import (
	"errors"
	"image"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
)

// Video interface.
type Video interface {
	WriteImage(imagePath string) (err error)
	Close() (err error)
}

// Mp4 video convert options.
type Options struct {
	FPS         int
	codec       string
	pixelformat string
	pixfmt      string
}

// Mp4 video struct.
type Mp4 struct {
	Video
	cmd       *exec.Cmd
	stdinWr   io.WriteCloser
	dimension string
	exePath   string
	output    string
	option    *Options
}

// New Instance.
func NewMp4(output string, op *Options) Video {
	// Set default option
	if op.codec == "" {
		op.codec = "libx264"
	}
	if op.pixelformat == "" {
		op.pixelformat = "yuv420p"
	}
	if op.pixfmt == "" {
		op.pixfmt = "rgba"
	}
	if op.FPS == 0 {
		op.FPS = 25
	}
	return &Mp4{cmd: nil, stdinWr: nil, dimension: "", output: output, option: op}
}

// Initialize FFmpeg thread.
func (m *Mp4) initialize() error {
	// to-do
	exe, err := GetFFmpegLib()
	if err != nil {
		return err
	}
	if m.dimension == "" {
		m.dimension = "512x512"
	}
	pix_fmt := m.option.pixfmt
	fps := strconv.Itoa(m.option.FPS)
	codec := m.option.codec
	pixelformat := m.option.pixelformat
	outputfile := m.output

	cmdstr := []string{"-y",
		"-f", "rawvideo",
		"-vcodec", "rawvideo",
		"-s", m.dimension,
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
func (m *Mp4) WriteImage(imagePath string) error {
	img, err := LoadImage(imagePath)
	if err != nil {
		return err
	}
	width, height, err := m.getImageDimension(imagePath)
	if err != nil {
		return err
	}
	dimension := strconv.Itoa(width) + "x" + strconv.Itoa(height)
	if m.dimension == "" {
		m.dimension = dimension
		if err := m.initialize(); err != nil {
			return err
		}
	}
	if dimension != m.dimension {
		return errors.New("All images in a movie should have same size.")
	}
	if img != nil && m.cmd != nil && m.stdinWr != nil {
		imgstring := LoadImageBitmap(img)
		m.stdinWr.Write(imgstring)
	}
	return nil
}

// Get image Dimension
func (m *Mp4) getImageDimension(imagePath string) (int, int, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return 0, 0, err
	}
	image, _, err := image.DecodeConfig(file)
	return image.Width, image.Height, err
}

// Close ffmpeg
func (m *Mp4) Close() error {
	if m.cmd == nil {
		return errors.New("FFmpeg command is nil.")
	}
	if m.stdinWr != nil {
		err := m.stdinWr.Close()
		if err != nil {
			return err
		}
		log.Print("Close the mp4 video.")
	}
	m.cmd.Wait()
	m.cmd = nil
	return nil
}

// Execute FFmpeg commands.
func (m *Mp4) execFFmpegCommands(commandName string, params []string) error {
	m.cmd = exec.Command(commandName, params...)
	stdinWr, err := m.cmd.StdinPipe()
	m.stdinWr = stdinWr
	if err != nil {
		return err
	}
	m.cmd.Start()
	return nil
}
