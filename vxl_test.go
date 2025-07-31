package main

import (
	"fmt"
	"testing"
)

func buildMesh(vxl *Vxl) []float32 { // 3 pos  3  clr
	items := vxl.Limbs[0].Items
	for x := 0; x < len(items); x++ {
		for y := 0; y < len(items[x]); y++ {
			for z := 0; z < len(items[x][y]); z++ {
				if items[x][y][z] == nil {
					continue
				}
				// TODO 上下左右前后

			}
		}
	}
	return nil
}

func TestVxl(t *testing.T) {
	data := ReadData("ra2.mix", "local.mix/harv.vxl")
	vxl := ParseVxl(data)

	mesh := buildMesh(vxl)
	fmt.Println(mesh)
}
