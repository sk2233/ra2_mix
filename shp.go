/*
@author: sk
@date: 2025/5/4
*/
package main

import (
	"bytes"
)

type AnimHeader struct {
	Unused        uint16
	Width, Height uint16 // 画布大小
	Count         uint16 // 帧数
}

type FrameHeader struct {
	// 只有区域 (x,y) (w,h) 的数据进行存储
	X, Y          uint16
	Width, Height uint16
	Flag          uint8
	Unused        [11]byte
	Offset        uint32 // 具体数据偏移
}

type Shp struct {
	Anim   *AnimHeader
	Frames []*FrameHeader
	Buffs  [][]byte // 解码后的数据
}

func ParseShp(data []byte) *Shp {
	reader := bytes.NewReader(data)
	anim := ReadAny[AnimHeader](reader)
	frames := make([]*FrameHeader, 0)
	buffs := make([][]byte, 0)
	for i := 0; i < int(anim.Count); i++ {
		frame := ReadAny[FrameHeader](reader)
		frames = append(frames, frame)
		buffs = append(buffs, ParseFrameBuff(frame, data))
	}
	return &Shp{
		Anim:   anim,
		Frames: frames,
		Buffs:  buffs,
	}
}

func ParseFrameBuff(frame *FrameHeader, data []byte) []byte {
	if frame.Offset == 0 { // 空白帧 忽略 保持上一帧内容
		return make([]byte, 0)
	}
	if frame.Flag&0x02 > 0 { // 使用压缩了需要解压
		res := make([]byte, 0)
		reader := bytes.NewReader(data[frame.Offset:])
		for i := 0; i < int(frame.Height); i++ { // 逐行编解码
			count := ReadU16(reader) - 2 // 该行一共 count byte 这里已经占用了 2 byte
			for count > 0 {
				clr := ReadU8(reader)
				count--
				if clr == 0 { // 只对透明色进行了游标编码
					res = append(res, make([]uint8, ReadU8(reader))...)
					count--
				} else {
					res = append(res, clr)
				}
			}
			res = res[:(i+1)*int(frame.Width)] // 有可能会多数据，多余的截断
		}
		return res
	} else { // 没有压缩
		size := uint32(frame.Width) * uint32(frame.Height)
		return data[frame.Offset : frame.Offset+size]
	}
}
