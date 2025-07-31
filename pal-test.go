package main

import (
	"image"
	"image/color"
)

func TestPal() {
	data := ReadData("ra2.mix", "sidec02.mix/sidebar.pal")
	pal := ParsePal(data)

	img := image.NewRGBA(image.Rect(0, 0, 16*16, 16*16))
	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			clr := pal[x+y*16]
			for i := 0; i < 16; i++ {
				for j := 0; j < 16; j++ {
					img.SetRGBA(x*16+j, y*16+i, color.RGBA{R: clr.Red, G: clr.Green, B: clr.Blue, A: 255})
				}
			}
		}
	}
	ShowImage(img)
}
