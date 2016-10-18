package imageio

import (
	"log"
	"testing"
	"time"
	"fmt"
	"image/png"
	"image/jpeg"
	"os"
)

func TestWritePNGImageTime(t *testing.T) {
	start := time.Now()
	mp4 := NewVideo("test.mp4", &Options{FPS:24})
	mp4.WriteImage("images/test.jpg")
	fmt.Printf("Write PNG image cost: %4f \n", time.Now().Sub(start).Seconds())
}

func TestWriteJPGImageTime(t *testing.T) {
	start := time.Now()
	mp4 := NewVideo("test.mp4", &Options{FPS:24})
	mp4.WriteImage("images/camera.png")
	fmt.Printf("Write JPG image cost: %4f \n", time.Now().Sub(start).Seconds())
}

func TestDecodePNGImageTime(t *testing.T) {
	file, err := os.Open("image/camera.png")
	if err != nil {
		fmt.Println(err)
	}
	start := time.Now()
	if _, err := png.Decode(file); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Decode png image cost: %4f \n", time.Now().Sub(start).Seconds())
}

func TestDecodeJPGImageTime(t *testing.T) {
	file, err := os.Open("image/test.jpg")
	if err != nil {
		fmt.Println(err)
	}
	start := time.Now()
	if _, err := jpeg.Decode(file); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Decode jpg image cost: %4f \n", time.Now().Sub(start).Seconds())
}

func TestWriteImage(t *testing.T) {
	start := time.Now()
	mp4 := NewVideo("test.mp4", &Options{FPS:24})
	for i := 0; i < 100; i++ {
		err1 := mp4.WriteImage("images/camera.png")
		err2 := mp4.WriteImage("images/test.jpg")
		if err1 != nil && err2 != nil {
			log.Fatal(err1)
		}
	}
	fmt.Printf("Test write multi image cost: %4f \n", time.Now().Sub(start).Seconds())
}