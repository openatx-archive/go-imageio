package imageio

import (
	"testing"
	"fmt"
)

var mp4 = NewVideo("test.mp4", &Options{FPS:24})

var imgjpg, _ = LoadImage("images/image720x720.jpg")

var imgpng, _ = LoadImage("images/image720x720.png")

func BenchmarkWriteJPEGImageFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mp4.WriteImageFile("images/image720x720.jpg")
	}
}

func BenchmarkWritePNGImageFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mp4.WriteImageFile("images/image720x720.png")
	}
}

func BenchmarkWriteJPEGImage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mp4.WriteImage(imgjpg)
	}
}

func BenchmarkWritePNGImage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mp4.WriteImage(imgpng)
	}
}

func BenchmarkDecodeJPEGImage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := LoadImage("images/image720x720.jpg")
		if err != nil {
			fmt.Println(err)
		}
	}
}

func BenchmarkDecodePNGImage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := LoadImage("images/image720x720.png")
		if err != nil {
			fmt.Println(err)
		}
	}
}
