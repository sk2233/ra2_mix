package main

import (
	"image"
	"image/color"
)

func TestPcx() {
	data := ReadData("ra2.mix", "local.mix/logo.pcx")
	pcx := ParsePcx(data)

	w, h := int(pcx.Header.GetWidth()), int(pcx.Header.GetHeight())
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := (y*w + x) * 3
			img.SetRGBA(x, y, color.RGBA{R: pcx.Data[idx], G: pcx.Data[idx+1], B: pcx.Data[idx+2], A: 255})
		}
	}
	ShowImage(img)
}
