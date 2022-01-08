package nrgba64

import (
	"context"
	"image"

	"encoding/binary"

	"bytes"

	"github.com/lucasew/imgcode/crypt"
	"github.com/lucasew/imgcode/utils"
)

// Encode encodes to grayscale image
func Encode(ctx context.Context, r []byte) (image.Image, error) {
	size := len(r) + 4
	side := utils.Size(size / 4)
	buf := bytes.NewBuffer([]byte{})
	buf.Grow(side * side)
	err := binary.Write(buf, binary.LittleEndian, int32(len(r)))
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(r)
	if err != nil {
		return nil, err
	}
	img := image.NewNRGBA(image.Rect(0, 0, side, side))
	copy(img.Pix, buf.Bytes())
	crypter := crypt.GetCrypterFromContext(ctx)
	if crypter != nil {
		err = crypt.InplaceEncrypt(*crypter, img.Pix)
		if err != nil {
			return nil, err
		}
	}
	return img, nil
}

// Decode decodes the image to data
func Decode(ctx context.Context, img image.Image) ([]byte, error) {
	pix := img.(*image.NRGBA).Pix
	crypter := crypt.GetCrypterFromContext(ctx)
	if crypter != nil {
		err := crypt.InplaceDecrypt(*crypter, pix)
		if err != nil {
			return nil, err
		}
	}
	bbuf := bytes.NewBuffer(pix)
	var size int32
	err := binary.Read(bbuf, binary.LittleEndian, &size)
	println(size)
	if err != nil {
		return nil, err
	}
	raw := bbuf.Bytes()
	return raw[0:size], nil
}
