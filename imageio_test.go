package imageio

import (
	"log"
	"testing"
	"time"
	"fmt"
)

func TestWritePNGImageTime(t *testing.T) {
	start := time.Now()
	mp4 := NewMp4("test.mp4", &Options{FPS:24})
	mp4.WriteImage("images/test.jpg")
	fmt.Println("cost: ", time.Now().Sub(start).Seconds())
}

func TestWriteJPGImageTime(t *testing.T) {
	start := time.Now()
	mp4 := NewMp4("test.mp4", &Options{FPS:24})
	mp4.WriteImage("images/camera.png")
	fmt.Println("cost: ", time.Now().Sub(start).Seconds())
}

func TestWriteImage(t *testing.T) {
	start := time.Now()
	mp4 := NewMp4("test.mp4", &Options{FPS:24})
	for i := 0; i < 100; i++ {
		err1 := mp4.WriteImage("images/camera.png")
		err2 := mp4.WriteImage("images/test.jpg")
		if err1 != nil && err2 != nil {
			log.Fatal(err1)
		}
	}
	fmt.Println("cost: ", time.Now().Sub(start).Seconds())
}