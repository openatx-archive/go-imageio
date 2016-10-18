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

> System: Win 7 Memory:8G CPU; Core(TM) i5-4590

```
-- Decode --
TestDecodePNGImage      0.60       ms/op
TestDecodeJPGImage      0.18       ms/op
-- Write --
TestWritePNGImage       62.82      ms/op
TestWriteJPGImage       24.70      ms/op
```

### Reference

- [ffmpeg](https://www.ffmpeg.org/) 
- [imaging](https://github.com/disintegration/imaging) 

### LICENSE

Under LICENSE [MIT](https://github.com/openatx/go-stf/blob/master/LICENSE) 