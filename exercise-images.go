package main

import (
	"golang.org/x/tour/pic"
	"image/color"
	"image"
)

type Image struct{
	w, h int
}

func (img Image) ColorModel() color.Model {
	return color.RGBAModel
}

func (img Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, img.w, img.h)
}

func (img Image) At(x, y int) color.Color {
	c := uint8(x^y)
	return color.RGBA{c, c, 255, 255}
}

func main() {
	m := Image{150, 150}
	pic.ShowImage(m)
}
