package imageio

import (
	"bufio"
	"errors"
	"fmt"
	"image"
	"image/color"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"path/filepath"
	"image/png"
	"image/jpeg"
	"image/gif"
)

// Ensure we have our version of the binary freeimage lib.
func GetFFmpegExe() (string, error) {
	inPath, err := CheckIfFFmpegInPATH()
	if inPath && err == nil {
		return "ffmpeg", err
	}
	plat := GetPlatform()
	if localFile, ok := FnamePerPlatform[plat]; ok {
		if _, err := os.Stat(localFile); os.IsNotExist(err) {
			return localFile, GetRomoteFile(plat)
		} else {
			isExe, err := CheckIfFileExecutable(localFile)
			if isExe && err == nil {
				return localFile, err
			} else {
				return localFile, GetRomoteFile(plat)
			}
		}
	}
	return "", errors.New("Platform not exist")
}

// Check if ffmpeg is in System PATH
func CheckIfFFmpegInPATH() (bool, error) {
	return CheckFFmpegVersion("ffmpeg")
}

// Check if the exe file is excutable.
func CheckIfFileExecutable(filepath string) (bool, error) {
	return CheckFFmpegVersion(filepath)
}

func CheckFFmpegVersion(filepath string) (bool, error) {
	cmd := exec.Command(filepath, "-version")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return false, err
	}
	defer stdout.Close()
	cmd.Start()
	reader := bufio.NewReader(stdout)
	line, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	if strings.Contains(line, "ffmpeg") {
		return true, nil
	} else {
		return false, errors.New("Local file is not excutable.")
	}
	return true, nil
}

// Get a string that specifies the platform more specific than
// The result can be: linux32, linux64, win32,
// win64, osx32, osx64. Other platforms may be added in the future.
func GetPlatform() (platform string) {
	bitNum := 32 << uintptr(^uintptr(0) >> 63)
	switch runtime.GOOS {
	case "darwin":
		platform = "osx" + strconv.Itoa(bitNum)
		break
	case "windows":
		platform = "win" + strconv.Itoa(bitNum)
		break
	case "linux":
		platform = "linux" + strconv.Itoa(bitNum)
		break
	}
	return platform
}

// Load image.
func LoadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	switch filepath.Ext(path) {
	case ".png":
		return png.Decode(file)
	case ".jpg":
		return jpeg.Decode(file)
	case ".gif":
		return gif.Decode(file)
	default:
		return nil, errors.New("Unkown file format.")
	}
	return jpeg.Decode(file)
}

// Load Image Bitmap.
func LoadImageBitmap(imgfile image.Image) []uint8 {
	fmt.Printf("format %T\n", imgfile)
	srcBounds := imgfile.Bounds()
	srcMinX := srcBounds.Min.X
	srcMinY := srcBounds.Min.Y

	dstBounds := srcBounds.Sub(srcBounds.Min)
	dstW := dstBounds.Dx()
	dstH := dstBounds.Dy()
	dst := image.NewNRGBA(dstBounds)

	switch src := imgfile.(type) {
	case *image.NRGBA:
		rowSize := srcBounds.Dx() * 4
		parallel(dstH, func(partStart, partEnd int) {
			for dstY := partStart; dstY < partEnd; dstY++ {
				di := dst.PixOffset(0, dstY)
				si := src.PixOffset(srcMinX, srcMinY + dstY)
				copy(dst.Pix[di:di + rowSize], src.Pix[si:si + rowSize])
			}
		})
	case *image.NRGBA64:
		parallel(dstH, func(partStart, partEnd int) {
			for dstY := partStart; dstY < partEnd; dstY++ {
				di := dst.PixOffset(0, dstY)
				si := src.PixOffset(srcMinX, srcMinY + dstY)
				for dstX := 0; dstX < dstW; dstX++ {

					dst.Pix[di + 0] = src.Pix[si + 0]
					dst.Pix[di + 1] = src.Pix[si + 2]
					dst.Pix[di + 2] = src.Pix[si + 4]
					dst.Pix[di + 3] = src.Pix[si + 6]

					di += 4
					si += 8

				}
			}
		})
	case *image.RGBA:
		parallel(dstH, func(partStart, partEnd int) {
			for dstY := partStart; dstY < partEnd; dstY++ {
				di := dst.PixOffset(0, dstY)
				si := src.PixOffset(srcMinX, srcMinY + dstY)
				for dstX := 0; dstX < dstW; dstX++ {

					a := src.Pix[si + 3]
					dst.Pix[di + 3] = a
					switch a {
					case 0:
						dst.Pix[di + 0] = 0
						dst.Pix[di + 1] = 0
						dst.Pix[di + 2] = 0
					case 0xff:
						dst.Pix[di + 0] = src.Pix[si + 0]
						dst.Pix[di + 1] = src.Pix[si + 1]
						dst.Pix[di + 2] = src.Pix[si + 2]
					default:
						var tmp uint16
						tmp = uint16(src.Pix[si + 0]) * 0xff / uint16(a)
						dst.Pix[di + 0] = uint8(tmp)
						tmp = uint16(src.Pix[si + 1]) * 0xff / uint16(a)
						dst.Pix[di + 1] = uint8(tmp)
						tmp = uint16(src.Pix[si + 2]) * 0xff / uint16(a)
						dst.Pix[di + 2] = uint8(tmp)
					}

					di += 4
					si += 4

				}
			}
		})
	case *image.RGBA64:
		parallel(dstH, func(partStart, partEnd int) {
			for dstY := partStart; dstY < partEnd; dstY++ {
				di := dst.PixOffset(0, dstY)
				si := src.PixOffset(srcMinX, srcMinY + dstY)
				for dstX := 0; dstX < dstW; dstX++ {

					a := src.Pix[si + 6]
					dst.Pix[di + 3] = a
					switch a {
					case 0:
						dst.Pix[di + 0] = 0
						dst.Pix[di + 1] = 0
						dst.Pix[di + 2] = 0
					case 0xff:
						dst.Pix[di + 0] = src.Pix[si + 0]
						dst.Pix[di + 1] = src.Pix[si + 2]
						dst.Pix[di + 2] = src.Pix[si + 4]
					default:
						var tmp uint16
						tmp = uint16(src.Pix[si + 0]) * 0xff / uint16(a)
						dst.Pix[di + 0] = uint8(tmp)
						tmp = uint16(src.Pix[si + 2]) * 0xff / uint16(a)
						dst.Pix[di + 1] = uint8(tmp)
						tmp = uint16(src.Pix[si + 4]) * 0xff / uint16(a)
						dst.Pix[di + 2] = uint8(tmp)
					}

					di += 4
					si += 8

				}
			}
		})

	case *image.Gray:
		parallel(dstH, func(partStart, partEnd int) {
			for dstY := partStart; dstY < partEnd; dstY++ {
				di := dst.PixOffset(0, dstY)
				si := src.PixOffset(srcMinX, srcMinY + dstY)
				for dstX := 0; dstX < dstW; dstX++ {

					c := src.Pix[si]
					dst.Pix[di + 0] = c
					dst.Pix[di + 1] = c
					dst.Pix[di + 2] = c
					dst.Pix[di + 3] = 0xff

					di += 4
					si += 1

				}
			}
		})

	case *image.Gray16:
		parallel(dstH, func(partStart, partEnd int) {
			for dstY := partStart; dstY < partEnd; dstY++ {
				di := dst.PixOffset(0, dstY)
				si := src.PixOffset(srcMinX, srcMinY + dstY)
				for dstX := 0; dstX < dstW; dstX++ {

					c := src.Pix[si]
					dst.Pix[di + 0] = c
					dst.Pix[di + 1] = c
					dst.Pix[di + 2] = c
					dst.Pix[di + 3] = 0xff

					di += 4
					si += 2

				}
			}
		})

	case *image.YCbCr:
		parallel(dstH, func(partStart, partEnd int) {
			for dstY := partStart; dstY < partEnd; dstY++ {
				di := dst.PixOffset(0, dstY)
				for dstX := 0; dstX < dstW; dstX++ {

					srcX := srcMinX + dstX
					srcY := srcMinY + dstY
					siy := src.YOffset(srcX, srcY)
					sic := src.COffset(srcX, srcY)
					r, g, b := color.YCbCrToRGB(src.Y[siy], src.Cb[sic], src.Cr[sic])
					dst.Pix[di + 0] = r
					dst.Pix[di + 1] = g
					dst.Pix[di + 2] = b
					dst.Pix[di + 3] = 0xff

					di += 4

				}
			}
		})

	case *image.Paletted:
		plen := len(src.Palette)
		pnew := make([]color.NRGBA, plen)
		for i := 0; i < plen; i++ {
			pnew[i] = color.NRGBAModel.Convert(src.Palette[i]).(color.NRGBA)
		}

		parallel(dstH, func(partStart, partEnd int) {
			for dstY := partStart; dstY < partEnd; dstY++ {
				di := dst.PixOffset(0, dstY)
				si := src.PixOffset(srcMinX, srcMinY + dstY)
				for dstX := 0; dstX < dstW; dstX++ {

					c := pnew[src.Pix[si]]
					dst.Pix[di + 0] = c.R
					dst.Pix[di + 1] = c.G
					dst.Pix[di + 2] = c.B
					dst.Pix[di + 3] = c.A

					di += 4
					si += 1

				}
			}
		})

	default:
		parallel(dstH, func(partStart, partEnd int) {
			for dstY := partStart; dstY < partEnd; dstY++ {
				di := dst.PixOffset(0, dstY)
				for dstX := 0; dstX < dstW; dstX++ {

					c := color.NRGBAModel.Convert(imgfile.At(srcMinX + dstX, srcMinY + dstY)).(color.NRGBA)
					dst.Pix[di + 0] = c.R
					dst.Pix[di + 1] = c.G
					dst.Pix[di + 2] = c.B
					dst.Pix[di + 3] = c.A

					di += 4

				}
			}
		})

	}
	return dst.Pix
}

var parallelizationEnabled = true

// if GOMAXPROCS = 1: no goroutines used
// if GOMAXPROCS > 1: spawn N=GOMAXPROCS workers in separate goroutines
func parallel(dataSize int, fn func(partStart, partEnd int)) {
	numGoroutines := 1
	partSize := dataSize

	if parallelizationEnabled {
		numProcs := runtime.GOMAXPROCS(0)
		if numProcs > 1 {
			numGoroutines = numProcs
			partSize = dataSize / (numGoroutines * 10)
			if partSize < 1 {
				partSize = 1
			}
		}
	}
	if numGoroutines == 1 {
		fn(0, dataSize)
	} else {
		var wg sync.WaitGroup
		wg.Add(numGoroutines)
		idx := uint64(0)

		for p := 0; p < numGoroutines; p++ {
			go func() {
				defer wg.Done()
				for {
					partStart := int(atomic.AddUint64(&idx, uint64(partSize))) - partSize
					if partStart >= dataSize {
						break
					}
					partEnd := partStart + partSize
					if partEnd > dataSize {
						partEnd = dataSize
					}
					fn(partStart, partEnd)
				}
			}()
		}
		wg.Wait()
	}
}
