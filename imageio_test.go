package imageio

import (
	"log"
	"testing"
)

func TestWriteImage(t *testing.T) {
	mp4 := NewMp4("test.mp4", &Options{})
	for i := 0; i < 100; i++ {
		err := mp4.WriteImage("../images/camera.png")
		if err != nil {
			log.Fatal(err)
		}
	}
}
