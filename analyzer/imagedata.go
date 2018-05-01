package analyzer

import (
	"errors"
	"fmt"
    "path/filepath"
    "image"
    "image/jpeg"
    "image/png"
    "os"
    "time"
)

type imageData interface {
    RGBAAvgs() []float32
    Size() image.Point

    Scan() [][]Pixel
    RGBAAt(int, int) []uint8
    CalcAverages() []float32
    ValidPixels() int
    Elapsed() time.Duration

    SetElapsed(start, end time.Time)
}

type ImageData struct {
    image image.Image
    size image.Point
    elapsed time.Duration

    pixels [][]Pixel
    borders []Pixel

    rgbaSums []uint
    rgbaAvgs []float32
    validPixels int
}

func newImageData(filename string, c chan ImageManagerItem) {
    defer wg.Done()

    start := time.Now()
    item := ImageManagerItem{filename: filename}

    // looks for the filename in gallery folder
    reader, err := os.Open(fmt.Sprintf("./gallery/%s", filename))
    if err != nil {
        item.error = err
        c <- item
        return
    }
    
    // determine how the file should be decoded from its extension
    var image image.Image
    switch filepath.Ext(filename) {
        case ".jpg", ".jpeg":
            image, _ = jpeg.Decode(reader)
        case ".png":
            image, _ = png.Decode(reader)
        default:
            item.error = errors.New("File extension not supported.")
            c <- item
            return
    }
    
    maxBounds := image.Bounds().Max
    imgd := &ImageData{
        image: image,
        size: maxBounds,
    }
    
    imgd.Scan()
    imgd.FilterPixels()
    imgd.CalcAverages()
    imgd.SetElapsed(start, time.Now()) // end timer
    
    item.imageData = imgd
    c <- item
}

func (imgd *ImageData) SetElapsed(start, end time.Time) {
    imgd.elapsed = end.Sub(start)
}

func (imgd *ImageData) RGBAAvgs() []float32 {
    return imgd.rgbaAvgs
}

func (imgd *ImageData) Size() image.Point {
    return imgd.size
}

func (imgd *ImageData) ValidPixels() int {
    return imgd.validPixels
}

func (imgd *ImageData) Elapsed() time.Duration {
    return imgd.elapsed
}

/*
Iterate each pixel in the image.
Each pixel's RGBA values are stored in a 3d array
and the sums of all RGBA colors are collected.
*/
func (imgd *ImageData) Scan() [][]Pixel {
    // reset props
    imgd.validPixels = 0
    imgd.rgbaSums = make([]uint, 4)
    imgd.pixels = make([][]Pixel, imgd.size.Y)

    for y := 0; y < imgd.size.Y; y++ {
        pixelRow := make([]Pixel, imgd.size.X)
        
        // assign each nested entry a pixel
        for x := 0; x < imgd.size.X; x++ {
            pixel := Pixel{
                rgba: imgd.RGBAAt(x, y),
                x: x,
                y: y,
            }
            pixel.white = pixel.IsWhite()
            pixel.black = pixel.IsBlack()
            
            pixelRow[x] = pixel
        }
        
        imgd.pixels[y] = pixelRow
    }

    return imgd.pixels
}

/*
Get the uint32 RGBA values then convert it to uint8 (0-255)
*/
func (imgd *ImageData) RGBAAt(x, y int) []uint8 {
    r, g, b, a := imgd.image.At(x, y).RGBA()
    return []uint8{
        uint8(r>>8),
        uint8(g>>8),
        uint8(b>>8),
        uint8(a>>8),
    }
}

/*
Iterate each pixel to check its validity, then collect data of the ones that are valid.
*/
func (imgd *ImageData) FilterPixels() {
    imgd.borders = []Pixel{}

    if filterProfile.contiguous {
        // will first gather invalid pixels on image's edges
        invalidPixels := []Pixel{}
        // run through the top and bottom image edge
        for x := 0; x < imgd.size.X; x++ {
            topPixel := &imgd.pixels[0][x]
            if !imgd.Validate(topPixel) {
                invalidPixels = append(invalidPixels, *topPixel)
            }
            btmPixel := &imgd.pixels[imgd.size.Y-1][x]
            if !imgd.Validate(btmPixel) {
                invalidPixels = append(invalidPixels, *btmPixel)
            }
        }
        // run through left and right image edge
        for y := 1; y < imgd.size.Y-1; y++ {
            leftPixel := &imgd.pixels[y][0]
            if !imgd.Validate(leftPixel) {
                invalidPixels = append(invalidPixels, *leftPixel)
            }
            rightPixel := &imgd.pixels[y][imgd.size.X-1]
            if !imgd.Validate(rightPixel) {
                invalidPixels = append(invalidPixels, *rightPixel)
            }
        }
        // check adjacent pixels of the invalid edge pixels
        imgd.FilterAdjacent(invalidPixels)
    
        // process data of valid pixels
        for _, pixelRow := range imgd.pixels {
            for _, pixel := range pixelRow {
                if pixel.state != states.invalid {
                    imgd.validPixels++
                    for i, v := range pixel.rgba {
                        imgd.rgbaSums[i] += uint(v)
                    }
                }
            }
        }

    } else {
        for _, pixelRow := range imgd.pixels {
            for _, pixel := range pixelRow {
                if imgd.Validate(&pixel) {
                    imgd.validPixels++
                    for i, v := range pixel.rgba {
                        imgd.rgbaSums[i] += uint(v)
                    }
                }
            }
        }
    }
}

/*
From an array of invalid pixels, check all adjacent pixels for its validity.
*/
func (imgd *ImageData) FilterAdjacent(invalidPixels []Pixel) {
    newInvalidPixels := []Pixel{}
    for _, p := range invalidPixels {
        if p.x > 0 {
            leftPixel := &imgd.pixels[p.y][p.x-1]
            if !imgd.Validate(leftPixel) {
                newInvalidPixels = append(newInvalidPixels, *leftPixel)
            }
        }
        if p.x < imgd.size.X - 1 {
            rightPixel := &imgd.pixels[p.y][p.x+1]
            if !imgd.Validate(rightPixel) {
                newInvalidPixels = append(newInvalidPixels, *rightPixel)
            }
        }
        if p.y > 0 {
            topPixel := &imgd.pixels[p.y-1][p.x]
            if !imgd.Validate(topPixel) {
                newInvalidPixels = append(newInvalidPixels, *topPixel)
            }
        }
        if p.y < imgd.size.Y - 1 {
            btmPixel := &imgd.pixels[p.y+1][p.x]
            if !imgd.Validate(btmPixel) {
                newInvalidPixels = append(newInvalidPixels, *btmPixel)
            }
        }
    }
    
    // call function recursively until there are no more invalid pixels to be found
    if len(newInvalidPixels) > 0 {
        imgd.FilterAdjacent(newInvalidPixels)
    }
}

func (imgd *ImageData) Validate(pixel *Pixel) bool {
    if !pixel.tested {
        pixel.Test()
        
        if pixel.state == states.valid && filterProfile.contiguous {
            imgd.borders = append(imgd.borders, *pixel)
        }

        return pixel.state != states.invalid
    }

    return true
}

/*
Calculate the average RGBA values from the total of valid pixels
*/
func (imgd *ImageData) CalcAverages() []float32 {
    // reset averages
    imgd.rgbaAvgs = make([]float32, 4)

    for i, v := range imgd.rgbaSums {
        imgd.rgbaAvgs[i] = float32(v) / float32(imgd.validPixels)
    }
    
    return imgd.rgbaAvgs
}
