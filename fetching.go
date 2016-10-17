package imageio

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	RemoteResourceUrl = "https://github.com/imageio/imageio-binaries/raw/master/ffmpeg"
)

var (
	FnamePerPlatform = map[string]string{
		"osx32":   "ffmpeg.osx",
		"osx64":   "ffmpeg.osx",
		"win32":   "ffmpeg.win32.exe",
		"win64":   "ffmpeg.win32.exe",
		"linux32": "ffmpeg.linux32",
		"linux64": "ffmpeg.linux64",
	}
)

// The results from individual calls to it.
type PassThru struct {
	io.Reader
	CurrentSize int64 // Current # of bytes transferred
	TotalSize   int64 // Total # of bytes transferred
}

// Read 'overrides' the underlying io.Reader's Read method.
// This is the one that will be called by io.Copy(). We simply
// use it to keep track of byte counts and then forward the call.
func (pt *PassThru) Read(p []byte) (int, error) {
	n, err := pt.Reader.Read(p)
	pt.CurrentSize += int64(n)
	if err == nil {
		totalNum := int(pt.TotalSize) / n
		currentNum := int(pt.CurrentSize) / n
		if currentNum != 0 {
			fmt.Printf("\rDownloading: %.0f%%", float64(currentNum) / float64(totalNum) * 100)
			if float64(currentNum) / float64(totalNum) == 1.0 {
				fmt.Println()
			}
			os.Stdout.Sync()
		}
	}
	return n, err
}

// Get a filename for the local version of a file from the web
func GetRomoteFile(fname string) error {
	url := RemoteResourceUrl + "/" + FnamePerPlatform[fname]
	return downloadFromUrl(url, FnamePerPlatform[fname])
}

// Load requested file, downloading it if needed or requested.
func downloadFromUrl(url string, filename string) error {
	log.Println("FFmpeg was not found on your computer; downloading it now from", url)
	tmpFilename := filename + ".cache"
	output, err := os.Create(tmpFilename)
	if err != nil {
		return err
	}
	response, err := http.Get(url)
	fileSize := response.ContentLength
	log.Printf("Total file size : %v bytes\n", fileSize)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	src := &PassThru{Reader: response.Body, TotalSize: fileSize}
	count, err := io.Copy(output, src)
	if err != nil {
		return err
	}
	if err := output.Close(); err != nil {
		return err
	}
	if count == fileSize {
		fmt.Println("Transferred", count, "bytes")
		err := os.Rename(tmpFilename, filename)
		if err != nil {
			return err
		}
	}
	return nil
}
