package analyzer

type Pixel struct {
    rgba []uint8
    white, black, tested bool
    state uint8
    X, Y int
}

func (pixel *Pixel) IsWhite() bool {
    return pixel.rgba[0] == 255 && pixel.rgba[1] == 255 && pixel.rgba[2] == 255
}

func (pixel *Pixel) IsBlack() bool {
    return pixel.rgba[0] == 0 && pixel.rgba[1] == 0 && pixel.rgba[2] == 0
}

/*
Test if pixel is valid
*/
func (pixel *Pixel) Test() {
    if !pixel.tested {
        pixel.tested = true
        pixel.state = states.valid
    
        if pixel.white || pixel.black || pixel.rgba[3] == 0 {
            pixel.state = states.invalid
        }
    }
}
