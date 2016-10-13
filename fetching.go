package imageio

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	REMOTE_RESOURCE_URL = "https://github.com/imageio/imageio-binaries/raw/master/ffmpeg"
	FNAME_PER_PLATFORM = map[string]string{
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
	currentSize int64 // Current # of bytes transferred
	totalSize   int64 // Total # of bytes transferred
}

// Read 'overrides' the underlying io.Reader's Read method.
// This is the one that will be called by io.Copy(). We simply
// use it to keep track of byte counts and then forward the call.
func (pt *PassThru) Read(p []byte) (int, error) {
	n, err := pt.Reader.Read(p)
	pt.currentSize += int64(n)
	if err == nil {
		//time.Sleep(100 * time.Millisecond)
		totalNum := int(pt.totalSize) / n
		currentNum := int(pt.currentSize) / n
		//h := strings.Repeat("=", currentNum) + strings.Repeat(" ", totalNum - currentNum)
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
	// Get dirs to look for the resource
	url := REMOTE_RESOURCE_URL + "/" + FNAME_PER_PLATFORM[fname]
	return downloadFromUrl(url, FNAME_PER_PLATFORM[fname])
}

// Load requested file, downloading it if needed or requested.
func downloadFromUrl(url string, filename string) error {
	log.Println("FFmpeg was not found on your computer; downloading it now from", url)
	// Create local dictory.
	tmpFilename := filename + ".cache"
	output, err := os.Create(tmpFilename)
	if err != nil {
		return err
	}
	// Request from url
	response, err := http.Get(url)
	filesize := response.ContentLength
	log.Printf("Total file size : %v bytes\n", filesize)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	src := &PassThru{Reader: response.Body, totalSize: filesize}
	count, err := io.Copy(output, src)
	if err != nil {
		return err
	}
	if err := output.Close(); err != nil {
		return err
	}
	if count == filesize {
		fmt.Println("Transferred", count, "bytes")
		err := os.Rename(tmpFilename, filename)
		if err != nil {
			return err
		}
	}
	return nil
}
