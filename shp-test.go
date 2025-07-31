package main

import (
	"image"
	"image/color"
)

func TestShp() {
	data := ReadData("ra2.mix", "load.mix/ls800russia.shp")
	shp := ParseShp(data)
	data = ReadData("ra2md.mix", "loadmd.mix/mplsr.pal") // 不太好找对应的调色板
	pal := ParsePal(data)

	w, h := int(shp.Anim.Width), int(shp.Anim.Height)
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	frame := shp.Frames[0]
	buff := shp.Buffs[0]
	for y := 0; y < int(frame.Height); y++ {
		for x := 0; x < int(frame.Width); x++ {
			clr := pal[buff[x+y*int(frame.Width)]]
			img.SetRGBA(x+int(frame.X), y+int(frame.Y),
				color.RGBA{R: clr.Red, G: clr.Green, B: clr.Blue, A: 255})
		}
	}
	ShowImage(img)
}
