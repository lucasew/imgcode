package grayscale

import (
	"image"
	"image/color"

	"encoding/binary"

	"bytes"

	"github.com/lucasew/imgcode/utils"
)

// Encode encodes to grayscale image
func Encode(r []byte) (image.Image, error) {
	size := len(r)
	buf := bytes.NewBuffer([]byte{})
	err := binary.Write(buf, binary.LittleEndian, int32(len(r)))
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(r)
	if err != nil {
		return nil, err
	}
	bbuf := buf.Bytes()
	side := utils.Size(size)
	img := image.NewGray(image.Rect(0, 0, side, side))
	cur := 0
	for i := 0; i < side; i++ {
		for j := 0; j < side; j++ {
			if cur >= len(bbuf) {
				return img, nil
			}
			c := color.Gray{bbuf[cur]}
			img.Set(i, j, c)
			cur++
		}
	}
	return img, nil
}

// Decode decodes the image to data
func Decode(image image.Image) ([]byte, error) {
	bbuf := bytes.NewBuffer([]byte{})
	bbuf.Grow(image.Bounds().Dx() * image.Bounds().Dy())
	for i := 0; i < image.Bounds().Dx(); i++ {
		for j := 0; j < image.Bounds().Dy(); j++ {
			c := image.At(i, j).(color.Gray).Y
			bbuf.WriteByte(c)
		}
	}
	var size int32
	err := binary.Read(bbuf, binary.LittleEndian, &size)
	println(size)
	if err != nil {
		return nil, err
	}
	raw := bbuf.Bytes()
	return raw[0:size], nil
}
