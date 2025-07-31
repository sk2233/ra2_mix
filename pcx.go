package main

import "bytes"

type PcxHeader struct {
	CompanyMark            uint8
	Version                uint8
	CompressType           uint8 // 压缩方式 一般是 RLE 游标压缩
	BitCount               uint8
	MinX, MinY, MaxX, MaxY uint16
	WDpi, HDpi             uint16
	Unused1                [49]byte
	ColorPlanes            uint8  // 1 下标 or 3 直接就是颜色
	PadWidth               uint16 // 每行有多少 算 pad
	Unused2                [60]byte
}

func (h *PcxHeader) GetHeight() uint16 {
	return h.MaxY - h.MinY + 1
}

func (h *PcxHeader) GetWidth() uint16 {
	return h.MaxX - h.MinX + 1
}

type Pcx struct {
	Header *PcxHeader
	Data   []uint8 // RGB RGB ...
}

func ParsePcx(data []byte) *Pcx {
	header := ReadAny[PcxHeader](bytes.NewReader(data[:128])) // 解析文件头
	w, h := int(header.GetWidth()), int(header.GetHeight())
	padW := int(header.PadWidth)
	if header.ColorPlanes == 3 { // 3 直接就是颜色
		temp := ParseData(data[128:], header)
		res := make([]uint8, 0)
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				idx := x + y*padW
				res = append(res, AdjustColor(temp[idx]), AdjustColor(temp[idx+padW*h]), AdjustColor(temp[idx+padW*h*2]))
			}
		}
		return &Pcx{Header: header, Data: res}
	} else { // 1 索引
		end := len(data) - 256*3
		colors := ParseColors(data[end:])        // 解析调色盘
		temp := ParseData(data[128:end], header) // 解析主体
		res := make([]uint8, 0)
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				idx := temp[x+y*padW]
				clr := colors[idx]
				res = append(res, clr.Red, clr.Green, clr.Blue)
			}
		}
		return &Pcx{Header: header, Data: res}
	}
}

func ParseData(data []byte, header *PcxHeader) []uint8 {
	size := int(header.PadWidth) * int(header.GetHeight()) * int(header.ColorPlanes)
	if header.CompressType != 1 { // 没有压缩
		return data[:size]
	}
	// RLE 压缩
	res := make([]uint8, 0)
	idx := 0
	for idx+1 < len(data) {
		if data[idx] < 0xC0 {
			res = append(res, data[idx])
		} else { // 高两位都是 1
			count := int(data[idx] & 0x3F) // 那低 6位就是 数量
			idx++
			for i := 0; i < count; i++ {
				res = append(res, data[idx])
			}
		}
		idx++
	}
	return res[:size] // 只要需要的 Width 是有 pad 的 Width
}
