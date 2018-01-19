package image

import (
    "fmt"
    "image"
    "os"
    "image/png"
)

type PngImage struct {
    Img     image.Image
}

func LoadImage(file string) (PngImage, error) {
    image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

    // Load the start image.
    sImg, err := os.Open(file)
    if err != nil {
        fmt.Printf("Error while trying to load image testImage.bmp: %v.\n", err)
        return PngImage{}, err
    }
    img, _, err := image.Decode(sImg)
    if err != nil {
        fmt.Printf("Error while trying to decode image testImage.bmp: %v.\n", err)
        return PngImage{}, err
    }

    return PngImage{img}, nil
}

func (png *PngImage) RangeX() int {
    return png.Img.Bounds().Max.X - png.Img.Bounds().Min.X
}

func (png *PngImage) RangeY() int {
    return png.Img.Bounds().Max.Y - png.Img.Bounds().Min.Y
}

func (png *PngImage) RGBAAt(x, y int, flippedY bool) (float32, float32, float32, float32) {

    coordY := y
    if flippedY {
        coordY = png.RangeY()-y
    }
    r,g,b,a := png.Img.At(x, coordY).RGBA()
    return float32(r/257), float32(g/257), float32(b/257), float32(a/257)
}


























