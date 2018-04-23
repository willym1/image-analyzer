package test

import "github.com/willym1/image-analyzer/autoimage"

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
