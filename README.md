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
Elapsed: 17.0285ms
Image size: 400 x 456
Average RGB of apple.jpg: [164.48085 80.440674 78.99653]

ERROR: gradient.gif - File extension not supported.

Elapsed: 25.0413ms
Image size: 510 x 490
Average RGB of green-apple.jpeg: [146.45868 182.07582 92.9763]

Elapsed: 22.5255ms
Image size: 343 x 603
Average RGB of louise-belcher.png: [174.96744 137.9245 99.00828]

Elapsed: 48.058ms
Image size: 780 x 400
Average RGB of orange-full.jpg: [242.04344 140.11331 27.245275]

Elapsed: 29.0161ms
Image size: 520 x 520
Average RGB of orange.jpg: [248.61993 148.3836 58.055096]

ERROR: invalidimage - open ./gallery/invalidimage: The system cannot find the file specified.

Elapsed: 105.6415ms
Image size: 1300 x 796
Average RGB of skyline.jpg: [68.83822 51.30977 64.37151]
```
