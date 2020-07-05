package vp

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/nfnt/resize"
	"golang.org/x/crypto/ssh/terminal"
)

const defualtRatio = 8.0 / 18.0
const maxUint = ^uint32(0)

// Image //
type Image struct {
	Input  image.Image
	Bounds image.Rectangle
}

// NewImageFromFile will open an image from a given filepath
func NewImageFromFile(filename string) (*Image, error) {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return NewImage(f)
}

// NewImage will return a new image from a []byte
func NewImage(bts []byte) (*Image, error) {
	contentType := http.DetectContentType(bts)
	buff := bytes.NewBuffer(bts)
	var img image.Image
	var err error

	switch contentType {
	case "image/jpeg":
		img, err = jpeg.Decode(buff)
		if err != nil {
			return nil, err
		}
	case "image/png":
		img, err = png.Decode(buff)
		if err != nil {
			return nil, err
		}
	default:
		return nil, err
	}

	bounds := img.Bounds()

	return &Image{
		Input:  img,
		Bounds: bounds,
	}, nil
}

// Print will print the image to StdOut
func (i Image) Print() error {
	w, h, err := terminal.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}

	// Remove one line for height to not cut off the top
	h = h - 1

	ws := float64(w)
	hs := float64(h) / defualtRatio
	wi := float64(i.Bounds.Dx())
	hi := float64(i.Bounds.Dy())

	rs := ws / hs
	ri := wi / hi

	var nw float64
	var nh float64
	if rs > ri {
		nw = wi * hs / hi
		nh = hs
	} else {
		nw = ws
		nh = hi * ws / wi
	}

	nh = nh * defualtRatio

	m := resize.Resize(uint(nw), uint(nh), i.Input, resize.Lanczos3)

	for y := 0; y < m.Bounds().Dy(); y++ {
		for x := 0; x < m.Bounds().Dx(); x++ {
			c := m.At(x, y)
			r, g, b, _ := c.RGBA()

			ur := uint8(r >> 8)
			ug := uint8(g >> 8)
			ub := uint8(b >> 8)

			fmt.Printf("\x1b[48;2;%d;%d;%dm \x1b[0m", ur, ug, ub)
		}
		fmt.Print("\n")
	}

	return nil
}
