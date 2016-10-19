# go-imageio
> Golang wraps for image IO read and write.

### Usages

```go
mp4 := NewMp4("test.mp4", &Options{})
	for i := 0; i < 100; i++ {
		err := mp4.WriteImage("camera.png")
		if err != nil {
			log.Printf(err)
		}
	}
```

### Result
2016-10-19
> System: Win 7
> Memory: 8G
> CPU: Core(TM) i5-4570 3.20GHz

```
-- Write --
BenchmarkWriteJPEGImageFile                100          25.75 ms/op
BenchmarkWritePNGImageFile                  20          60.60 ms/op
BenchmarkWriteJPEGImage                    100          17.11 ms/op
BenchmarkWritePNGImage                     100          14.33 ms/op
-- Decode --
BenchmarkDecodeJPEGImage                   100          17.54 ms/op
BenchmarkDecodePNGImage                     20          57.90 ms/op
```

2016-10-19
> Raspberry 2

```
-- Write --
BenchmarkWriteJPEGImageFile                  2         648 ms/op
BenchmarkWritePNGImageFile                   1        1271 ms/op
BenchmarkWriteJPEGImage                      3         436 ms/op
BenchmarkWritePNGImage                       3         400 ms/op
-- Decode --
BenchmarkDecodeJPEGImage                     2         510 ms/op
BenchmarkDecodePNGImage                      1        1229 ms/op
```

### Reference

- [ffmpeg](https://www.ffmpeg.org/) 
- [imaging](https://github.com/disintegration/imaging) 

### LICENSE

Under LICENSE [MIT](https://github.com/openatx/go-stf/blob/master/LICENSE) 
