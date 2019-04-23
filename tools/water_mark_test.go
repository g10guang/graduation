package tools

import (
	"fmt"
	"image/jpeg"
	"os"
	"testing"
)

func TestWaterMark(t *testing.T) {
	fp, err := os.Open("/Users/g10guang/Public/output.jpeg")
	if err != nil {
		panic(err)
	}
	im, err := jpeg.Decode(fp)
	//im, err := png.Decode(fp)
	if err != nil {
		panic(err)
	}
	dim := ConvertImage2Draw(im)
	if dim == nil {
		fmt.Printf("convert ok")
		WaterMark(dim, "g10guang")
	} else {
		fmt.Printf("bad case")
	}

	// 输出到磁盘
	jpegW, err := os.OpenFile("./test.jpeg", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	if err = JpegCompress(im, jpegW); err != nil {
		panic(err)
	}

	pngW, err := os.OpenFile("./test.png", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer pngW.Close()
	if err = PngCompress(im, pngW); err != nil {
		panic(err)
	}
}
