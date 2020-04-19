package video

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/lucasew/imgcode/crypt"
)

const width = 1920
const height = 1080
const headerSize = 4

// const inspectableArtifacts = true
const inspectableArtifacts = false

// Decode decodes a byte stream from a video
func Decode(ctx context.Context, w io.Writer, fromFile string) error {
	tmpdir, err := ioutil.TempDir("", "imgdecode")
	if inspectableArtifacts {
		println(tmpdir)
	} else {
		defer os.RemoveAll(tmpdir) // Cleanup
	}
	if err != nil {
		return err
	}
	ffmpeg_flags := []string{}
	// input
	ffmpeg_flags = append(ffmpeg_flags, "-i", fromFile)
	// out file
	ffmpeg_flags = append(ffmpeg_flags, path.Join(tmpdir, "%04d.png"))

	cmd := exec.Command("ffmpeg", ffmpeg_flags...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	ith := 1
	for {
		fmt.Printf("Decoding frame %d...\n", ith)
		f, err := os.Open(path.Join(tmpdir, fmt.Sprintf("%04d.png", ith)))
		if os.IsNotExist(err) {
			break
		}
		if err != nil {
			return err
		}
		defer f.Close()
		img, _, err := image.Decode(f)
		if err != nil {
			return err
		}
		pixels := img.(*image.NRGBA).Pix
		buf := bytes.NewBuffer([]byte{})
		crypter := crypt.GetCrypterFromContext(ctx)
		if crypter != nil {
			err = crypt.InplaceDecrypt(*crypter, pixels)
			if err != nil {
				return err
			}
		}
		_, err = buf.Write(pixels)
		if err != nil {
			return err
		}
		var size int32
		binary.Read(buf, binary.LittleEndian, &size)
		stepWrite, err := w.Write(buf.Bytes()[0:size])
		if err != nil {
			return err
		}
		if inspectableArtifacts {
			fmt.Printf("Written %d bytes\n", stepWrite)
		}
		ith++
	}
	return nil
}

// Encode encodes a byte stream to a video
func Encode(ctx context.Context, r io.Reader, outfile string) error {
	tmpdir, err := ioutil.TempDir("", "imgcode")
	if inspectableArtifacts {
		println(tmpdir)
	} else {
		defer os.RemoveAll(tmpdir) // Cleanup
	}
	if err != nil {
		return err
	}
	size := uint64(0)
	ith := 1
	originEOF := false
	for {
		fmt.Printf("Coding frame %d...\n", ith)
		img := image.NewNRGBA(image.Rect(0, 0, width, height))
		fetched, err := r.Read(img.Pix[4:])
		if fetched < (width*height)-headerSize {
			originEOF = true
		}
		if err != nil && err != io.EOF {
			return err
		}
		sizebuf := bytes.NewBuffer([]byte{})
		err = binary.Write(sizebuf, binary.LittleEndian, int32(fetched))
		if err != nil {
			return err
		}
		_, err = sizebuf.Read(img.Pix[0:4])
		if err != nil {
			return err
		}
		imgname := path.Join(tmpdir, fmt.Sprintf("%04d.png", ith))
		f, err := os.Create(imgname)
		if err != nil {
			return err
		}
		defer f.Close()
		crypter := crypt.GetCrypterFromContext(ctx)
		if crypter != nil {
			err = crypt.InplaceEncrypt(*crypter, img.Pix)
			if err != nil {
				return err
			}
		}
		err = png.Encode(f, img)
		if err != nil {
			return err
		}
		err = f.Sync()
		if err != nil {
			return err
		}
		size += uint64(fetched)
		ith++
		if originEOF {
			break
		}
	}
	if inspectableArtifacts {
		println(size)
		println(ith)
	}
	ffmpeg_params := []string{}
	// input files
	ffmpeg_params = append(ffmpeg_params, "-i", path.Join(tmpdir, "%04d.png"))
	// overwrite
	ffmpeg_params = append(ffmpeg_params, "-y")
	// lossless
	ffmpeg_params = append(ffmpeg_params, "-vcodec", "rawvideo")
	// out file
	ffmpeg_params = append(ffmpeg_params, outfile+".avi")

	cmd := exec.Command("ffmpeg", ffmpeg_params...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return err
}
