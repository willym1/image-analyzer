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
Elapsed: 17.526ms
Image size: 400 x 456
Average RGB: 164 80 79

-----gradient.gif-----
ERROR: - File extension not supported.

-----green-apple.jpeg-----
Elapsed: 23.0251ms
Image size: 510 x 490
Average RGB: 146 182 93

-----louise-belcher.png-----
Elapsed: 33.0317ms
Image size: 343 x 603
Average RGB: 175 138 99

-----orange-full.jpg-----
Elapsed: 49.0464ms
Image size: 780 x 400
Average RGB: 242 140 27

-----orange.jpg-----
Elapsed: 30.0287ms
Image size: 520 x 520
Average RGB: 249 148 58

-----invalidimage-----
ERROR: - open ./gallery/invalidimage: The system cannot find the file specified.

-----skyline.jpg-----
Elapsed: 105.1003ms
Image size: 1300 x 796
Average RGB: 69 51 64
```