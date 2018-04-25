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
-----apple.jpg-----
Elapsed: 30.0206ms
Image size: 400 x 456
Average RGB: 165 81 80
Valid pixels: 127336

-----gradient.gif-----
ERROR: File extension not supported.

-----green-apple.jpeg-----
Elapsed: 61.079ms
Image size: 510 x 490
Average RGB: 147 182 94
Valid pixels: 170773

-----invalidimage-----
ERROR: open ./gallery/invalidimage: The system cannot find the file specified.

-----louise-belcher.png-----
Elapsed: 67.064ms
Image size: 343 x 603
Average RGB: 175 138 99
Valid pixels: 44549

-----orange-full.jpg-----
Elapsed: 48.0498ms
Image size: 780 x 400
Average RGB: 242 140 27
Valid pixels: 312000

-----orange.jpg-----
Elapsed: 64.0615ms
Image size: 520 x 520
Average RGB: 249 149 58
Valid pixels: 168931

-----red-blocks.png-----
Elapsed: 6.509ms
Image size: 250 x 250
Average RGB: 255 44 44
Valid pixels: 48939

-----skyline.jpg-----
Elapsed: 122.6534ms
Image size: 1300 x 796
Average RGB: 69 51 64
Valid pixels: 1034632
```
