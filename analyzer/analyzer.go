package analyzer

import (
    "fmt"
    "math"
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
    items []ImageManagerItem
    Logging bool
}

type ImageManagerItem struct {
    imageData
    error

    filename string
}

/*
Receives a string of filenames then
make new image data from each image
*/
func NewImages(filenames []string, l bool) ImageManager {
    var items []ImageManagerItem

    ch := make(chan ImageManagerItem, len(filenames))
    for _, filename := range filenames {
        wg.Add(1)
        go newImageData(filename, ch)
    }

    wg.Wait()
    close(ch)
    
    for item := range ch {
        items = append(items, item)
    }

    manager := ImageManager{items: items, Logging: l}
    if manager.Logging {
        manager.Log()
    }

    return manager
}

func (manager ImageManager) Log() {
    for _, v := range manager.items {
        if v.imageData != nil {
            size := v.imageData.Size()
            rgba := v.imageData.RGBAAvgs()
    
            fmt.Printf("-----%s-----\n", v.filename)
            fmt.Printf("Elapsed: %s\n", v.imageData.Elapsed())
            fmt.Printf("Image size: %d x %d\n", size.X, size.Y)
            fmt.Printf("Average RGB: %v %v %v\n", math.Round(float64(rgba[0])), math.Round(float64(rgba[1])), math.Round(float64(rgba[2])))
            fmt.Printf("Valid pixels: %d\n", v.imageData.ValidPixels())

        } else {
            fmt.Printf("-----%s-----\n", v.filename)
            fmt.Printf("ERROR: %s\n", v.error)
        }
        fmt.Println("")
    }
}
