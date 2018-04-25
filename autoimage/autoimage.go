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
                fmt.Printf("Valid pixels: %d\n", imgd.ValidPixels())
                fmt.Println("")
            }

        } else if l {
            fmt.Printf("-----%s-----\n", v)
            fmt.Printf("ERROR: %s\n\n", err)
        }
    }

    ai := Images{images: images, logging: l}
    return ai
}

type imageData interface {
    RGBAvgs() []float32
    Size() image.Point

    Scan() [][]Pixel
    RGBAt(int, int) []uint8
    CalcAverages() []float32
    ValidPixels() int
}

type ImageData struct {
    image image.Image
    size image.Point

    pixels [][]Pixel
    rgbSums []uint
    rgbAvgs []float32
    validPixels int

    ignoreEdges bool
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
        ignoreEdges: true,
    }
    
    imgd.Scan()
    imgd.FilterPixels()
    imgd.CalcAverages()

    return imgd, nil
}

func (imgd *ImageData) RGBAvgs() []float32 {
    return imgd.rgbAvgs
}

func (imgd *ImageData) Size() image.Point {
    return imgd.size
}

func (imgd *ImageData) ValidPixels() int {
    return imgd.validPixels
}

/*
Iterate each pixel in the image.
Each pixel's RGB values are stored in a 3d array
and the sums of all RGB colors are collected.
*/
func (imgd *ImageData) Scan() [][]Pixel {
    // reset props
    imgd.validPixels = 0
    imgd.rgbSums = make([]uint, 3)
    imgd.pixels = make([][]Pixel, imgd.size.Y)

    for y := 0; y < imgd.size.Y; y++ {
        pixelRow := make([]Pixel, imgd.size.X)
        
        // assign each nested entry a pixel
        for x := 0; x < imgd.size.X; x++ {
            pixel := Pixel{
                rgb: imgd.RGBAt(x, y)[:3],
                valid: true,
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
Iterate each pixel to check its validity, then collect data of the ones that are valid.
*/
func (imgd *ImageData) FilterPixels() {
    // will first gather invalid pixels on image's borders
    invalidPixels := []Pixel{}
    // run through the top and bottom image border
    for x := 0; x < imgd.size.X; x++ {
        topPixel := &imgd.pixels[0][x]
        if topPixel.Validate() {
            invalidPixels = append(invalidPixels, *topPixel)
        }
        btmPixel := &imgd.pixels[imgd.size.Y-1][x]
        if btmPixel.Validate() {
            invalidPixels = append(invalidPixels, *btmPixel)
        }
    }
    // run through left and right image border
    for y := 1; y < imgd.size.Y-1; y++ {
        leftPixel := &imgd.pixels[y][0]
        if leftPixel.Validate() {
            invalidPixels = append(invalidPixels, *leftPixel)
        }
        rightPixel := &imgd.pixels[y][imgd.size.X-1]
        if rightPixel.Validate() {
            invalidPixels = append(invalidPixels, *rightPixel)
        }
    }
    // check adjacent pixels of the invalid border pixels
    imgd.FilterAdjacent(invalidPixels)

    // process data of valid pixels
    for _, pixelRow := range imgd.pixels {
        for _, pixel := range pixelRow {
            if pixel.valid {
                imgd.validPixels++
                for i, v := range pixel.rgb {
                    imgd.rgbSums[i] += uint(v)
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
            if leftPixel.ValidateWith(p) {
                newInvalidPixels = append(newInvalidPixels, *leftPixel)
            }
        }
        if p.x < imgd.size.X - 1 {
            rightPixel := &imgd.pixels[p.y][p.x+1]
            if rightPixel.ValidateWith(p) {
                newInvalidPixels = append(newInvalidPixels, *rightPixel)
            }
        }
        if p.y > 0 {
            topPixel := &imgd.pixels[p.y-1][p.x]
            if topPixel.ValidateWith(p) {
                newInvalidPixels = append(newInvalidPixels, *topPixel)
            }
        }
        if p.y < imgd.size.Y - 1 {
            btmPixel := &imgd.pixels[p.y+1][p.x]
            if btmPixel.ValidateWith(p) {
                newInvalidPixels = append(newInvalidPixels, *btmPixel)
            }
        }
    }
    
    // call function recursively until there are no more invalid pixels to be found
    if len(newInvalidPixels) > 0 {
        imgd.FilterAdjacent(newInvalidPixels)
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

type Pixel struct {
    rgb []uint8
    black bool
    white bool
    valid bool
    x int
    y int
}

func (pixel *Pixel) IsWhite() bool {
    return pixel.rgb[0] == 255 && pixel.rgb[1] == 255 && pixel.rgb[2] == 255
}

func (pixel *Pixel) IsBlack() bool {
    return pixel.rgb[0] == 0 && pixel.rgb[1] == 0 && pixel.rgb[2] == 0
}

/*
Test if pixel is valid
*/
func (pixel *Pixel) Validate() bool {
    // invalidate black or white pixels
    if pixel.white || pixel.black {
        pixel.valid = false
        return true
    }
    
    return false
}

/*
Test if pixel is valid with a second one that's invalid
*/
func (pixel *Pixel) ValidateWith(withP Pixel) bool {
    // invalidate if it matches with the invalid pixel
    if pixel.valid && ((pixel.black && withP.black) || (pixel.white && withP.white)) {
        pixel.valid = false
        return true
    }

    return false
}
