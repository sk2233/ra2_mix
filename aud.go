package main

import (
	"bytes"
	"fmt"
	"math"
)

type AudHeader struct {
	SampleRate uint16
	DataSize   uint32 // 存储压缩过的内容大小
	OutputSize uint32 // 最终输出大小
	Flags      uint8  // 0x1 双通道还是单通道    0x2 采样是16位还是8位
	Format     uint8  // 1 WestwoodCompressed  99 ImaAdpcm
}

func (h *AudHeader) GetChannelCount() int {
	if h.Flags&0x1 == 0x1 {
		return 2
	} else {
		return 1
	}
}

type Aud struct {
	Header *AudHeader
	Data   []byte
}

func ParseAud(data []byte) *Aud {
	reader := bytes.NewReader(data)
	header := ReadAny[AudHeader](reader) // 读取头信息
	return &Aud{
		Header: header,
		Data:   ParseAudData(reader, header),
	}
}

func ParseAudData(reader *bytes.Reader, header *AudHeader) []byte {
	switch header.Format {
	case 1: // WestwoodCompressed
		return parseWestwoodAudData(reader, header)
	case 99: // ImaAdpcm
		return parseAdpcmAudData(reader, header)
	default:
		panic(fmt.Sprint("aud: unknown format ", header.Format))
	}
}

type Chunk struct {
	DataSize   uint16
	OutputSize uint16
	Mark       [4]byte
}

var (
	adpcmIndex  = 0
	adpcmSample = 0
)

var adpcmIndexAdjust = []int{-1, -1, -1, -1, 2, 4, 6, 8}
var adpcmDeltaMap = []int{
	7, 8, 9, 10, 11, 12, 13, 14, 16,
	17, 19, 21, 23, 25, 28, 31, 34, 37,
	41, 45, 50, 55, 60, 66, 73, 80, 88,
	97, 107, 118, 130, 143, 157, 173, 190, 209,
	230, 253, 279, 307, 337, 371, 408, 449, 494,
	544, 598, 658, 724, 796, 876, 963, 1060, 1166,
	1282, 1411, 1552, 1707, 1878, 2066, 2272, 2499, 2749,
	3024, 3327, 3660, 4026, 4428, 4871, 5358, 5894, 6484,
	7132, 7845, 8630, 9493, 10442, 11487, 12635, 13899, 15289,
	16818, 18500, 20350, 22385, 24623, 27086, 29794, 32767,
}

func decodeAdpcmAudSample(val byte) int16 {
	// 一个 byte 拆分为两份使用，每份 最高位是符号，低 3 位是偏移
	reversal := (val & 8) != 0
	val &= 7
	delta := adpcmDeltaMap[adpcmIndex]*int(val)/4 + adpcmDeltaMap[adpcmIndex]/8
	if reversal {
		delta = -delta
	}
	// 纠正本次采样
	adpcmSample += delta
	if adpcmSample > math.MaxInt16 {
		adpcmSample = math.MaxInt16
	}
	if adpcmSample < math.MinInt16 {
		adpcmSample = math.MinInt16
	}
	// 纠正索引
	adpcmIndex += adpcmIndexAdjust[val]
	if adpcmIndex < 0 {
		adpcmIndex = 0
	}
	if adpcmIndex > 88 {
		adpcmIndex = 88
	}
	return int16(adpcmSample)
}

func parseAdpcmAudData(reader *bytes.Reader, header *AudHeader) []byte {
	adpcmIndex = 0 // 初始化
	adpcmSample = 0
	lastSize := header.DataSize
	res := make([]byte, 0)
	for lastSize > 0 { // 只读取指定数目的字节
		chunk := ReadAny[Chunk](reader)
		bs := ReadBytes(reader, int(chunk.DataSize))
		for _, item := range bs {
			// 一个字节拆分为两部分进行解码
			temp := decodeAdpcmAudSample(item)
			res = append(res, byte(temp), byte(temp>>8))
			temp = decodeAdpcmAudSample(item >> 4)
			res = append(res, byte(temp), byte(temp>>8))
		}
		lastSize -= 8 + uint32(chunk.DataSize) // size(Chunk)+chunk.DataSize = 本次读取的字节数
	}
	return res
}

var westwoodStep2Map = []int{-2, -1, 0, 1}
var westwoodStep4Map = []int{-9, -8, -6, -5, -4, -3, -2, -1, 0, 1, 2, 3, 4, 5, 6, 8}

func toByte(val int) byte {
	if val < 0 {
		return 0
	} else if val > 255 {
		return 255
	} else {
		return byte(val)
	}
}

func decodeWestwoodAudSample(bs []byte) []byte {
	sample := 0x80
	idx := 0
	res := make([]byte, 0)
	for idx < len(bs) {
		// 命令 byte 分为两部分使用
		count := int(bs[idx]) & 0x3f
		mode := bs[idx] >> 6
		idx++
		switch mode {
		case 0:
			for count++; count > 0; count-- {
				code := bs[idx]
				res = append(res, toByte(sample+westwoodStep2Map[(code>>0)&0x03]))
				res = append(res, toByte(sample+westwoodStep2Map[(code>>2)&0x03]))
				res = append(res, toByte(sample+westwoodStep2Map[(code>>4)&0x03]))
				res = append(res, toByte(sample+westwoodStep2Map[(code>>6)&0x03]))
				idx++
			}
		case 1:
			for count++; count > 0; count-- {
				code := bs[idx]
				res = append(res, toByte(sample+westwoodStep4Map[(code>>0)&0x0F]))
				res = append(res, toByte(sample+westwoodStep4Map[(code>>4)&0x0F]))
				idx++
			}
		case 2:
			if (count & 0x20) != 0 { // 0x20 控制位
				sample = int(toByte(sample + count&0x1F)) // 只使用低6位
				res = append(res, toByte(sample))
			} else {
				for count++; count > 0; count-- { // 直接存储的采样
					res = append(res, bs[idx])
					idx++
				}
				sample = int(bs[idx-1])
			}
		default:
			for count++; count > 0; count-- {
				res = append(res, toByte(sample))
			}
		}
	}
	return res
}

func parseWestwoodAudData(reader *bytes.Reader, header *AudHeader) []byte {
	lastSize := header.DataSize
	res := make([]byte, 0)
	for lastSize > 0 { // 只读取指定数目的字节
		chunk := ReadAny[Chunk](reader)
		bs := ReadBytes(reader, int(chunk.DataSize))
		if chunk.DataSize == chunk.OutputSize {
			res = append(res, bs...) // 无压缩
		} else {
			res = append(res, decodeWestwoodAudSample(bs)...)
		}
		lastSize -= 8 + uint32(chunk.DataSize) // size(Chunk)+chunk.DataSize = 本次读取的字节数
	}
	return res
}
