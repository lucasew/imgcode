package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"

	codec "github.com/lucasew/imgcode/codecs/nrgba"
	video "github.com/lucasew/imgcode/codecs/video"
	"github.com/lucasew/imgcode/crypt"
	"github.com/lucasew/imgcode/utils"
)

var fileFrom = ""
var fileTo = ""

var appContext = context.Background()

func help() {
	println(`
		imgcode: Code and decode data from images and videos
		parameters
			- command: videnc, viddec imgenc imgdec
			- from: source file
			- to: destination file
			video encodings will be into avi and image encodings into png
			Google Photos tips:
			- data lose for images larger than 4000x4000
	`)
}

func main() {
	var passwd string
	flag.StringVar(&passwd, "p", "", "password to {en,de}crypt the datastream")
	flag.Parse()
	if flag.NArg() < 3 {
		help()
		utils.Check(fmt.Errorf("invalid input"))
	}
	if passwd != "" {
		crypter := crypt.NewCrypterFromPassword(passwd)
		appContext = crypt.ContextWithCrypter(appContext, crypter)
	}
	fileFrom = flag.Arg(1)
	fileTo = flag.Arg(2)
	switch flag.Arg(0) {
	case "videnc":
		videncode()
	case "viddec":
		viddecode()
	case "imgenc":
		imgencode()
	case "imgdec":
		imgdecode()
	default:
		help()
		utils.Check(fmt.Errorf("invalid command"))
	}
}

func videncode() {
	r, err := os.Open(fileFrom)
	utils.Check(err)
	defer r.Close()
	utils.Check(err)
	err = video.Encode(appContext, r, fileTo)
	utils.Check(err)
}

func viddecode() {
	w, err := os.Create(fileTo)
	utils.Check(err)
	err = video.Decode(appContext, w, fileFrom)
	utils.Check(err)
}

func imgencode() {
	r, err := os.Open(fileFrom)
	utils.Check(err)
	bytes, err := io.ReadAll(r)
	utils.Check(err)
	img, err := codec.Encode(appContext, bytes)
	utils.Check(err)
	outf, err := os.Create(fileTo)
	utils.Check(err)
	defer outf.Close()
	err = png.Encode(outf, img)
	utils.Check(err)
}

func imgdecode() {
	f, err := os.Open(fileFrom)
	utils.Check(err)
	defer f.Close()
	img, _, err := image.Decode(f)
	utils.Check(err)
	outf, err := os.Create(fileTo)
	utils.Check(err)
	defer outf.Close()
	raw, err := codec.Decode(appContext, img)
	utils.Check(err)
	r := bytes.NewBuffer(raw)
	utils.Check(err)
	outFile, err := os.Create(fileTo)
	utils.Check(err)
	_, err = io.Copy(outFile, r)
	utils.Check(err)
}
