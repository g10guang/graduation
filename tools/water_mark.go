package tools

import (
	"bufio"
	"bytes"
	"github.com/sirupsen/logrus"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/png"
)

// FIXME jpeg 图片无法添加水印
func WaterMark(im draw.Image, label string) {
	col := color.RGBA{255, 255, 255, 255}
	point := fixed.Point26_6{fixed.Int26_6(10 * 64), fixed.Int26_6(10 * 64)}
	d := &font.Drawer{
		Dst:  im,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}

// TODO find better method
// 将 image.Image 转化为 draw.Image
func ConvertImage2Draw(im image.Image) draw.Image {
	var err error
	if im == nil {
		return nil
	}
	v, ok := im.(draw.Image)
	if ok {
		return v
	}

	buf := &bytes.Buffer{}
	if err = png.Encode(buf, im); err != nil {
		logrus.Errorf("png Encode Error: %s", err)
		return nil
	}
	reader := bufio.NewReader(buf)
	if im, err = png.Decode(reader); err != nil {
		logrus.Errorf("png Decode Error: %s", err)
		return nil
	}

	v, _ = im.(draw.Image)
	return v
}
