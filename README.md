# Golang Image Analyzer

Small library for scanning images and retrieving its data.

#### Example

Calling Run() from `test/test.go`
```go
/* Test images */
func Run(logging bool) {
    images := []string{
        "apple.jpg",
        "gradient.gif",
        "green-apple.jpeg",
        "louise-belcher.png",
        "orange-full.jpg",
        "orange.jpg",
        "invalidimage",
        "skyline.jpg",
    }
    autoimage.NewImages(images, logging)
}
```
Prints
```
-----gradient.gif-----
ERROR: File extension not supported.

-----invalidimage-----
ERROR: open ./gallery/invalidimage: The system cannot find the file specified.

-----red-blocks.png-----
Elapsed: 30.5296ms
Image size: 250 x 250
Average RGB: 255 44 44
Valid pixels: 48939

-----apple.jpg-----
Elapsed: 69.0656ms
Image size: 400 x 456
Average RGB: 165 81 80
Valid pixels: 127336

-----orange-full.jpg-----
Elapsed: 106.1015ms
Image size: 780 x 400
Average RGB: 242 140 27
Valid pixels: 312000

-----louise-belcher.png-----
Elapsed: 111.1058ms
Image size: 343 x 603
Average RGB: 175 138 99
Valid pixels: 44549

-----green-apple.jpeg-----
Elapsed: 118.6132ms
Image size: 510 x 490
Average RGB: 147 182 94
Valid pixels: 170773

-----orange.jpg-----
Elapsed: 121.6174ms
Image size: 520 x 520
Average RGB: 249 149 58
Valid pixels: 168931

-----skyline.jpg-----
Elapsed: 180.1719ms
Image size: 1300 x 796
Average RGB: 69 51 64
Valid pixels: 1034632
```
