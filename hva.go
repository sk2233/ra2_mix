package main

import (
	"bytes"
	"github.com/go-gl/mathgl/mgl32"
	"io"
)

type Hva struct {
	FrameCount uint32         // 动画帧数
	LimbCount  uint32         // 组件数目
	Transforms [][]mgl32.Mat4 // FrameCount , LimbCount
}

func ParseHva(data []byte) *Hva {
	reader := bytes.NewReader(data)
	_, err := reader.Seek(16, io.SeekStart)
	HandleErr(err) // 读取基本信息
	frameCount := ReadU32(reader)
	limbCount := ReadU32(reader)
	_, err = reader.Seek(int64(16*limbCount), io.SeekCurrent)
	HandleErr(err)
	// 读取变换矩阵信息
	idx := []int{0, 4, 8, 12, 1, 5, 9, 13, 2, 6, 10, 14} // 用于纠正顺序
	transforms := make([][]mgl32.Mat4, frameCount)
	for i := 0; i < len(transforms); i++ {
		transforms[i] = make([]mgl32.Mat4, limbCount)
		for j := 0; j < len(transforms[i]); j++ {
			transforms[i][j][3] = 0
			transforms[i][j][7] = 0
			transforms[i][j][11] = 0
			transforms[i][j][15] = 1
			for k := 0; k < 12; k++ {
				transforms[i][j][idx[k]] = ReadF32(reader)
			}
		}
	}
	return &Hva{
		FrameCount: frameCount,
		LimbCount:  limbCount,
		Transforms: transforms,
	}
}
