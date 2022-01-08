package main

import (
	"bytes"
	"context"
	"flag"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"os"

	codec "github.com/lucasew/imgcode/codecs/nrgba"
	video "github.com/lucasew/imgcode/codecs/video"
	"github.com/lucasew/imgcode/crypt"
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
		panic("invalid input")
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
		panic("invalid command")
	}
}

func videncode() {
	r, err := os.Open(fileFrom)
	if err != nil {
		panic(err)
	}
	defer r.Close()
	if err != nil {
		panic(err)
	}
	err = video.Encode(appContext, r, fileTo)
	if err != nil {
		panic(err)
	}
}

func viddecode() {
	w, err := os.Create(fileTo)
	if err != nil {
		panic(err)
	}
	err = video.Decode(appContext, w, fileFrom)
	if err != nil {
		panic(err)
	}
}

func imgencode() {
	r, err := os.Open(fileFrom)
	if err != nil {
		panic(err)
	}
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	img, err := codec.Encode(appContext, bytes)
	if err != nil {
		panic(err)
	}
	outf, err := os.Create(fileTo)
	if err != nil {
		panic(err)
	}
	defer outf.Close()
	err = png.Encode(outf, img)
	if err != nil {
		panic(err)
	}
}

func imgdecode() {
	f, err := os.Open(fileFrom)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}
	outf, err := os.Create(fileTo)
	if err != nil {
		panic(err)
	}
	defer outf.Close()
	raw, err := codec.Decode(appContext, img)
	if err != nil {
		panic(err)
	}
	r := bytes.NewBuffer(raw)
	if err != nil {
		panic(err)
	}
	outFile, err := os.Create(fileTo)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(outFile, r)
	if err != nil {
		panic(err)
	}
}
