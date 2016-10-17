package imageio

import (
	"log"
	"testing"
)

func TestGetFFmpegLib(t *testing.T) {
	exe, err := GetFFmpegExe()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(exe)
}

func TestGetPlatform(t *testing.T) {
	plat := GetPlatform()
	log.Printf(plat)
}

func TestCheckIfFileExcutable(t *testing.T) {
	plat := GetPlatform()
	_, err := CheckIfFileExecutable(FNAME_PER_PLATFORM[plat])
	if err != nil {
		log.Fatal(err)
	}
}

func TestLoadImage(t *testing.T) {
	_, err := LoadImage("images/camera.png")
	if err != nil {
		log.Fatal(err)
	}
}
