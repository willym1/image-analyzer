package test

import "github.com/willym1/image-analyzer/analyzer"

/* Test images */
func Run() {
    filenames := []string{
        "apple.jpg",
		"gradient.gif",
		"green-apple.jpeg",
		"invalidimage",
        "louise-belcher.png",
        "orange-full.jpg",
        "orange.jpg",
        "red-blocks.png",
        "skyline.jpg",
    }
    manager := analyzer.ImagesFromFilenames(filenames)
    manager.ProcessItems()
    manager.Log()
}
