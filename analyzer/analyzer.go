package analyzer

import (
	"bytes"
    "fmt"
    "image"
    _ "image/jpeg"
    _ "image/png"
    "io"
    "io/ioutil"
    "math"
    "mime/multipart"
    "os"
    "sync"
)

var (
    wg sync.WaitGroup

    states = States{
        valid: uint8(1),
        invalid: uint8(2),
    }

    filterProfile = FilterProfile{
        contiguous: true,
    }
)

type States struct {
    valid, invalid uint8
}

type FilterProfile struct {
    contiguous bool
}

type ImageManager struct {
    Items []ImageManagerItem
}

/*
Receives a string of filenames then make new image data from each image
*/
func ImagesFromFilenames(filenames []string) ImageManager {
    var items []ImageManagerItem
    
    for _, filename := range filenames {
        item := ImageManagerItem{Filename: filename}

        // looks for the filename in gallery folder
        reader, err := os.Open(fmt.Sprintf("./gallery/%s", filename))
        
        if err == nil {
            item.Decode(reader)
        } else {
            item.Error = err
        }

        items = append(items, item)
    }
    
    return ImageManager{items}
}

/*
Decodes images from a multipart reader
*/
func ImagesFromParts(partReader *multipart.Reader) ImageManager {
    var items []ImageManagerItem
    var err error

    // loop all parts until EOF
    for err == nil {
        var part *multipart.Part
        part, err = partReader.NextPart()

        if err == nil {
            var byts []byte
            byts, err = ioutil.ReadAll(part)
            bytesReader := bytes.NewReader(byts)
    
            if err == nil {
                item := ImageManagerItem{Filename: part.FileName()}
                item.Decode(bytesReader)
                items = append(items, item)
            }
        }
    }

    return ImageManager{items}
}

func (manager ImageManager) ProcessItems() {
    ch := make(chan *ImageData, len(manager.Items))
    
    for _, item := range manager.Items {
        if item.ImageData != nil {
            wg.Add(1)
            go item.ImageData.Process(ch)
        }
    }
    
    wg.Wait()
    close(ch)
}

func (manager ImageManager) Log() {
    for _, item := range manager.Items {
        if item.ImageData != nil {
            size := item.ImageData.GetSize()
            rgba := item.ImageData.GetRgbaAvgs()
    
            fmt.Printf("-----%s-----\n", item.Filename)
            fmt.Printf("Elapsed: %s\n", item.ImageData.GetElapsed())
            fmt.Printf("Image size: %d x %d\n", size.X, size.Y)
            fmt.Printf("Average RGB: %v %v %v\n", math.Round(float64(rgba[0])), math.Round(float64(rgba[1])), math.Round(float64(rgba[2])))
            fmt.Printf("Valid pixels: %d\n", item.ImageData.GetValidPixels())

        } else {
            fmt.Printf("-----%s-----\n", item.Filename)
            fmt.Printf("ERROR: %s\n", item.Error)
        }
        fmt.Println("")
    }
}

type ImageManagerItem struct {
    ImageData imageData
    Error error

    Filename string
}

func (item *ImageManagerItem) Decode(reader io.Reader) {
    // determine how the file should be decoded from its extension
    img, _, err := image.Decode(reader)
    if err != nil {
        item.Error = err
    }
    
    if item.Error == nil {
        maxBounds := img.Bounds().Max
        item.ImageData = &ImageData{
            image: img,
            Size: maxBounds,
        }
    }
}
