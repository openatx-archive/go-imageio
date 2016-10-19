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
func (this *Video) initialize() error {
	exe, err := GetFFmpegExe()
	if err != nil {
		return err
	}
	if this.Dimension == "" {
		this.Dimension = "512x512"
	}
	pix_fmt := this.Option.Pixfmt
	fps := strconv.Itoa(this.Option.FPS)
	codec := this.Option.Codec
	pixelformat := this.Option.Pixelformat
	outputfile := this.Output

	cmdstr := []string{"-y",
		"-f", "rawvideo",
		"-vcodec", "rawvideo",
		"-s", this.Dimension,
		"-pix_fmt", pix_fmt,
		"-r", fps,
		"-i", "-", "-an",
		"-vcodec", codec,
		"-pix_fmt", pixelformat,
		"-crf", "25",
		"-r", "50",
		"-v", "warning", outputfile}
	return this.execFFmpegCommands(exe, cmdstr)
}

// Write image by file path.
func (this *Video) WriteImageFile(imagePath string) error {
	img, err := LoadImage(imagePath)
	if err != nil {
		return err
	}
	width, height, err := this.getImageDimension(imagePath)
	if err != nil {
		return err
	}
	dimension := fmt.Sprintf("%dx%d", width, height)
	if this.Dimension == "" {
		this.Dimension = dimension
		if err := this.initialize(); err != nil {
			return err
		}
	}
	if dimension != this.Dimension {
		return errors.New("All images in a movie should have same size.")
	}
	if img != nil && this.Cmd != nil && this.StdinWr != nil {
		imgstring := LoadImageBitmap(img)
		this.StdinWr.Write(imgstring)
	}
	return nil
}

// Write image by image.Image
func (this *Video) WriteImage(image image.Image) error {
	width := image.Bounds().Size().X
	height := image.Bounds().Size().Y
	dimension := fmt.Sprintf("%dx%d", width, height)
	if this.Dimension == "" {
		this.Dimension = dimension
		if err := this.initialize(); err != nil {
			return err
		}
	}
	if dimension != this.Dimension {
		return errors.New("All images in a movie should have same size.")
	}
	if image != nil && this.Cmd != nil && this.StdinWr != nil {
		imagestring := LoadImageBitmap(image)
		_, err := this.StdinWr.Write(imagestring)
		if err != nil {
			return err
		}
	}
	return nil
}

// Get image Dimension
func (this *Video) getImageDimension(imagePath string) (int, int, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return 0, 0, err
	}
	image, _, err := image.DecodeConfig(file)
	return image.Width, image.Height, err
}

// Close ffmpeg
func (this *Video) Close() error {
	if this.Cmd == nil {
		return errors.New("FFmpeg command is nil.")
	}
	if this.StdinWr != nil {
		err := this.StdinWr.Close()
		if err != nil {
			return err
		}
		log.Print("Close the mp4 video.")
	}
	this.Cmd.Wait()
	this.Cmd = nil
	return nil
}

// Execute FFmpeg commands.
func (this *Video) execFFmpegCommands(commandName string, params []string) error {
	this.Cmd = exec.Command(commandName, params...)
	stdinWr, err := this.Cmd.StdinPipe()
	if err != nil {
		return err
	}
	this.StdinWr = stdinWr
	this.Cmd.Start()
	return nil
}
