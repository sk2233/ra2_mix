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
	ColorPlanes            uint8
	Width                  uint16
	Unused2                [60]byte
}

type Pcx struct {
	Header *PcxHeader
	Data   []uint8
	Colors []*Color // 256个
}

func ParsePcx(data []byte) *Pcx {
	header := ReadAny[PcxHeader](bytes.NewReader(data[:128])) // 解析文件头
	end := len(data) - 256*3
	temp := ParseData(data[128:end-1], header) // 解析主体
	colors := ParseColors(data[end:])          // 解析调色盘
	return &Pcx{Header: header, Data: temp, Colors: colors}
}

func ParseData(data []byte, header *PcxHeader) []uint8 {
	if header.CompressType != 1 { // 没有压缩
		return data
	}
	// RLE 压缩
	res := make([]uint8, 0)
	idx := 0
	for idx+1 < len(data) {
		if data[idx]&0xC0 == 0xC0 { // 高两位为  1
			count := int(data[idx] & 0x3F) // 那低 6位就是 数量
			idx++
			for i := 0; i < count; i++ {
				res = append(res, data[idx])
			}
		} else {
			res = append(res, data[idx])
		}
		idx++
	}
	return res
}
