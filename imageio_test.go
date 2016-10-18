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
	mp4 := NewVideo("test.mp4", &Options{FPS:24})
	costList := make([]float64, 100)
	for i := 0; i < 100; i++ {
		start := time.Now()
		mp4.WriteImage("images/test.png")
		cost := time.Now().Sub(start).Seconds()
		costList = append(costList, cost)
	}
	var sum float64
	for _, cost := range costList {
		sum += cost
	}
	fmt.Printf("Write PNG image cost: %4f \n", sum / 100.0)
	mp4.Close()
}

func TestWriteJPGImageTime(t *testing.T) {
	mp4 := NewVideo("test.mp4", &Options{FPS:24})
	costList := make([]float64, 100)
	for i := 0; i < 100; i++ {
		start := time.Now()
		mp4.WriteImage("images/test.jpg")
		cost := time.Now().Sub(start).Seconds()
		costList = append(costList, cost)
	}
	var sum float64
	for _, cost := range costList {
		sum += cost
	}
	fmt.Printf("Write JPG image cost: %4f \n", sum / 100.0)
	mp4.Close()
}

func TestDecodePNGImageTime(t *testing.T) {
	file, err := os.Open("images/test.png")
	if err != nil {
		fmt.Println(err)
	}
	costList := make([]float64, 100)
	for i := 0; i < 100; i++ {
		start := time.Now()
		if _, err := png.Decode(file); err != nil {
			fmt.Println(err)
		}
		cost := time.Now().Sub(start).Seconds()
		costList = append(costList, cost)
	}
	var sum float64
	for _, cost := range costList {
		sum += cost
	}
	fmt.Printf("Decode png image cost: %4f \n", sum / 100.0)
}

func TestDecodeJPGImageTime(t *testing.T) {
	file, err := os.Open("images/test.jpg")
	if err != nil {
		fmt.Println(err)
	}
	costList := make([]float64, 100)
	for i := 0; i < 100; i++ {
		start := time.Now()
		if _, err := jpeg.Decode(file); err != nil {
			fmt.Println(err)
		}
		cost := time.Now().Sub(start).Seconds()
		costList = append(costList, cost)
	}
	var sum float64
	for _, cost := range costList {
		sum += cost
	}
	fmt.Printf("Decode jpeg image cost: %4f \n", sum / 100.0)
}

func TestWriteImage(t *testing.T) {
	start := time.Now()
	mp4 := NewVideo("test.mp4", &Options{FPS:24})
	for i := 0; i < 100; i++ {
		err1 := mp4.WriteImage("images/test.png")
		err2 := mp4.WriteImage("images/test.jpg")
		if err1 != nil && err2 != nil {
			log.Fatal(err1)
		}
	}
	mp4.Close()
	fmt.Printf("Test write multi image cost: %4f \n", time.Now().Sub(start).Seconds())
}