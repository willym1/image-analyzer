package analyzer

import (
    "image"
    "time"
)

type imageData interface {
    GetRgbaAvgs() []float32
    GetSize() image.Point
    GetValidPixels() int
    GetBorders() []Pixel
    GetElapsed() time.Duration
    
    Process(c chan *ImageData)
}

type ImageData struct {
    image image.Image
    Size image.Point
    Elapsed time.Duration

    pixels [][]Pixel
    Borders []Pixel

    rgbaSums []uint
    RgbaAvgs []float32
    ValidPixels int
}

func (imgd *ImageData) Process(c chan *ImageData) {
    defer wg.Done()

    start := time.Now() // start timer
    imgd.Scan()
    imgd.FilterPixels()
    imgd.CalcAverages()
    imgd.Elapsed = time.Now().Sub(start) // end timer
    
    c <- imgd
}

func (imgd *ImageData) GetRgbaAvgs() []float32 {
    return imgd.RgbaAvgs
}

func (imgd *ImageData) GetSize() image.Point {
    return imgd.Size
}

func (imgd *ImageData) GetValidPixels() int {
    return imgd.ValidPixels
}

func (imgd *ImageData) GetBorders() []Pixel {
    return imgd.Borders
}

func (imgd *ImageData) GetElapsed() time.Duration {
    return imgd.Elapsed
}

/*
Iterate each pixel in the image.
Each pixel's RGBA values are stored in a 3d array
and the sums of all RGBA colors are collected.
*/
func (imgd *ImageData) Scan() [][]Pixel {
    // reset props
    imgd.ValidPixels = 0
    imgd.rgbaSums = make([]uint, 4)
    imgd.pixels = make([][]Pixel, imgd.Size.Y)

    for y := 0; y < imgd.Size.Y; y++ {
        pixelRow := make([]Pixel, imgd.Size.X)
        
        // assign each nested entry a pixel
        for x := 0; x < imgd.Size.X; x++ {
            pixel := Pixel{
                rgba: imgd.RGBAAt(x, y),
                X: x,
                Y: y,
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
    imgd.Borders = []Pixel{}

    if filterProfile.contiguous {
        // will first gather invalid pixels on image's edges
        invalidPixels := []Pixel{}
        // run through the top and bottom image edge
        for x := 0; x < imgd.Size.X; x++ {
            topPixel := &imgd.pixels[0][x]
            if !imgd.Validate(topPixel) {
                invalidPixels = append(invalidPixels, *topPixel)
            }
            btmPixel := &imgd.pixels[imgd.Size.Y-1][x]
            if !imgd.Validate(btmPixel) {
                invalidPixels = append(invalidPixels, *btmPixel)
            }
        }
        // run through left and right image edge
        for y := 1; y < imgd.Size.Y-1; y++ {
            leftPixel := &imgd.pixels[y][0]
            if !imgd.Validate(leftPixel) {
                invalidPixels = append(invalidPixels, *leftPixel)
            }
            rightPixel := &imgd.pixels[y][imgd.Size.X-1]
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
                    imgd.ValidPixels++
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
                    imgd.ValidPixels++
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
        if p.X > 0 {
            leftPixel := &imgd.pixels[p.Y][p.X-1]
            if !imgd.Validate(leftPixel) {
                newInvalidPixels = append(newInvalidPixels, *leftPixel)
            }
        }
        if p.X < imgd.Size.X - 1 {
            rightPixel := &imgd.pixels[p.Y][p.X+1]
            if !imgd.Validate(rightPixel) {
                newInvalidPixels = append(newInvalidPixels, *rightPixel)
            }
        }
        if p.Y > 0 {
            topPixel := &imgd.pixels[p.Y-1][p.X]
            if !imgd.Validate(topPixel) {
                newInvalidPixels = append(newInvalidPixels, *topPixel)
            }
        }
        if p.Y < imgd.Size.Y - 1 {
            btmPixel := &imgd.pixels[p.Y+1][p.X]
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
            imgd.Borders = append(imgd.Borders, *pixel)
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
    imgd.RgbaAvgs = make([]float32, 4)

    for i, v := range imgd.rgbaSums {
        imgd.RgbaAvgs[i] = float32(v) / float32(imgd.ValidPixels)
    }
    
    return imgd.RgbaAvgs
}
