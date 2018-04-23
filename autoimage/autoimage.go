package autoimage

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"time"
)

type Images struct {
	images []imageData
	logging bool
}

/*
Receives a string of filenames then
make new image data from each image
*/
func NewImages(fns []string, l bool) Images {
	var images []imageData
    for _, v := range fns {
		// log information if enabled
		var start time.Time
		if l {
			start = time.Now()
		}

		imgd, err := initImageData(v)
		if err == nil {
			images = append(images, imgd)
			
			// log information if enabled
			if l {
				n := time.Now()
				s := imgd.Size()
				rgb := imgd.RGBAvgs()

				fmt.Printf("-----%s-----\n", v)
				fmt.Printf("Elapsed: %s\n", n.Sub(start))
				fmt.Printf("Image size: %d x %d\n", s.X, s.Y)
				fmt.Printf("Average RGB: %v %v %v\n", math.Round(float64(rgb[0])), math.Round(float64(rgb[1])), math.Round(float64(rgb[2])))
				fmt.Println("")
			}

		} else if l {
			fmt.Printf("-----%s-----\n", v)
			fmt.Printf("ERROR: - %s\n\n", err)
		}
	}

	ai := Images{images: images, logging: l}
	return ai
}


type imageData interface {
	RGBAvgs() []float32
	Size() image.Point

	Scan() [][][]uint8
	RGBAt(int, int) []uint8
	CalcAverages() []float32
}

type ImageData struct {
	image image.Image
	size image.Point

	pixels int
	validPixels int

	colorTable [][][]uint8
	rgbSums []uint
	rgbAvgs []float32
}

func initImageData(filename string) (imageData, error) {
	// looks for the filename in gallery folder
	reader, err := os.Open(fmt.Sprintf("./gallery/%s", filename))
	if err != nil {
		return nil, err
	}
	
	// determine how the file should be decoded from its extension
	var image image.Image
	switch filepath.Ext(filename) {
		case ".jpg", ".jpeg":
			image, _ = jpeg.Decode(reader)
		case ".png":
			image, _ = png.Decode(reader)
		default:
			return nil, errors.New("File extension not supported.")
	}
	
	maxBounds := image.Bounds().Max
	imgd := &ImageData{
		image: image,
		size: maxBounds,
		pixels: maxBounds.X * maxBounds.Y,
	}
	
	imgd.Scan()
	imgd.CalcAverages()

	return imgd, nil
}

func (imgd *ImageData) RGBAvgs() []float32 {
	return imgd.rgbAvgs
}

func (imgd *ImageData) Size() image.Point {
	return imgd.size
}

/*
Iterate each pixel in the image.
Each pixel's RGB values are stored in a 3d array
and the sums of all RGB colors are collected.
*/
func (imgd *ImageData) Scan() [][][]uint8 {
	// reset props
	imgd.validPixels = 0
	imgd.rgbSums = make([]uint, 3)
	imgd.colorTable = make([][][]uint8, imgd.size.Y)

	for i_y := 0; i_y < imgd.size.Y; i_y++ {
		imgd.colorTable[i_y] = make([][]uint8, imgd.size.X)

		// assign each nested entry a color
		for i_x := 0; i_x < imgd.size.X; i_x++ {
			rgb := imgd.RGBAt(i_x, i_y)[:3]
			imgd.colorTable[i_y][i_x] = rgb

			// append each color to sum
			if Filter(rgb) {
				imgd.validPixels++
				for i, v := range rgb {
					imgd.rgbSums[i] += uint(v)
				}
			}
		}
	}

	return imgd.colorTable
}

/*
Get the uint32 RGB values then convert it to uint8 (0-255)
*/
func (imgd *ImageData) RGBAt(x, y int) []uint8 {
	r, g, b, _ := imgd.image.At(x, y).RGBA()
	return []uint8{
		uint8(r>>8),
		uint8(g>>8),
		uint8(b>>8),
	}
}

/*
Calculate the average RGB values from the total of valid pixels
*/
func (imgd *ImageData) CalcAverages() []float32 {
	// reset averages
	imgd.rgbAvgs = make([]float32, 3)

	for i, v := range imgd.rgbSums {
		imgd.rgbAvgs[i] = float32(v) / float32(imgd.validPixels)
	}
	
	return imgd.rgbAvgs
}

/*
Determine which colors can be processed
*/
func Filter(rgb []uint8) bool {
	// skip if pixel is white or black
	if (rgb[0] == 255 && rgb[1] == 255 && rgb[2] == 255) || (rgb[0] == 0 && rgb[1] == 0 && rgb[2] == 0) {
		return false
	}

	return true
}
