package main

import (
	"fmt"
	"image"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/image/draw"
)

func main() {
	infile, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer infile.Close()

	if strings.HasSuffix(os.Args[1], ".gif") {
		println("its a gif!")
		displayGif(infile)
		return
	}

	src, _, err := image.Decode(infile)
	if err != nil {
		panic(err)
	}

	terminalWidth, terminalHeight, err := terminal.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}

	imgWidth := src.Bounds().Max.X
	imgHeight := src.Bounds().Max.Y

	scale := math.Min(float64(terminalWidth)/float64(imgWidth), float64((terminalHeight-1)*2)/float64(imgHeight))

	rect := image.Rect(0, 0, int(float64(imgWidth)*scale), int(float64(imgHeight)*scale))

	img := image.NewRGBA(rect)
	draw.NearestNeighbor.Scale(img, rect, src, src.Bounds(), draw.Over, nil)

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y += 2 {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			r1, g1, b1, _ := simplify(img.At(x, y).RGBA())
			r2, g2, b2, _ := simplify(img.At(x, y+1).RGBA())

			fmt.Printf("\x1b[48;2;%d;%d;%dm\x1b[38;2;%d;%d;%dm▄\x1b[0m", r1, g1, b1, r2, g2, b2)
		}
		fmt.Print("\n")
	}
}

func displayGif(infile *os.File) {
	gif, err := gif.DecodeAll(infile)
	if err != nil {
		panic(err)
	}

	for i, img := range gif.Image {

		displayImage(img)
		if i < len(gif.Image)-1 {
			fmt.Print("\033[0;0H")
		}

		time.Sleep(time.Duration(gif.Delay[i]*10) * time.Millisecond)
	}

	// _, terminalHeight, err := terminal.GetSize(int(os.Stdin.Fd()))
	// if err != nil {
	// 	panic(err)
	// }

}

func displayImage(src image.Image) {
	terminalWidth, terminalHeight, err := terminal.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}

	imgWidth := src.Bounds().Max.X
	imgHeight := src.Bounds().Max.Y

	scale := math.Min(float64(terminalWidth)/float64(imgWidth), float64((terminalHeight-1)*2)/float64(imgHeight))

	rect := image.Rect(0, 0, int(float64(imgWidth)*scale), int(float64(imgHeight)*scale))

	img := image.NewRGBA(rect)
	draw.NearestNeighbor.Scale(img, rect, src, src.Bounds(), draw.Over, nil)

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y += 2 {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			r1, g1, b1, _ := simplify(img.At(x, y).RGBA())
			r2, g2, b2, _ := simplify(img.At(x, y+1).RGBA())

			fmt.Printf("\x1b[48;2;%d;%d;%dm\x1b[38;2;%d;%d;%dm▄\x1b[0m", r1, g1, b1, r2, g2, b2)
		}
		fmt.Print("\n")
	}
}

func simplify(bigR, bigG, bigB, bigA uint32) (r, g, b, a uint8) {
	return uint8(bigR / 0x101), uint8(bigG / 0x101), uint8(bigB / 0x101), uint8(bigA / 0x101)
}
