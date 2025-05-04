/*
@author: sk
@date: 2025/5/4
*/
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"strings"
)

var (
	nameSet = map[string]bool{}
)

// 参考 https://www.zhihu.com/column/c_1899172031100069701

func main() {
	// 获取调色板
	bs, err := os.ReadFile("ra2.mix")
	HandleErr(err)
	res := ParseMix(bs)
	res = ParseMix(res["cache.mix"])
	pal := ParsePal(res["cameo.pal"])
	fmt.Println(len(pal))
	// 获取图标
	bs, err = os.ReadFile("language.mix")
	HandleErr(err)
	res = ParseMix(bs)
	res = ParseMix(res["cameo.mix"])
	for key, val := range res {
		if strings.HasSuffix(key, ".shp") {
			shp := ParseShp(val)
			SaveShp(shp, pal, key)
		}
	}
}

func SaveShp(shp *Shp, pal []*Color, key string) {
	for i, frame := range shp.Frames {
		if frame.Offset == 0 {
			continue // 复用上一帧的情况 这里直接舍弃
		}
		img := image.NewRGBA(image.Rect(0, 0, int(frame.Width), int(frame.Height)))
		buff := shp.Buffs[i]
		for y := 0; y < int(frame.Height); y++ {
			for x := 0; x < int(frame.Width); x++ {
				clr := pal[buff[x+y*int(frame.Width)]]
				img.Set(x+int(frame.X), y+int(frame.Y), color.RGBA{
					R: clr.Red,
					G: clr.Green,
					B: clr.Blue,
					A: 255,
				})
			}
		}
		file, err := os.Create(fmt.Sprintf("output/%s_%d.jpeg", key, i))
		HandleErr(err)
		err = jpeg.Encode(file, img, &jpeg.Options{
			Quality: 100,
		})
		HandleErr(err)
	}
}

func HandleByte(data []byte) {
	res := ParseMix(data)
	for key, val := range res {
		if strings.HasSuffix(key, ".mix") {
			HandleByte(val)
		}
		index := strings.LastIndex(key, ".")
		if index >= 0 {
			nameSet[key[index:]] = true
		}
		if strings.HasSuffix(key, ".shp") {
			ParseShp(val)
		}
	}
}
