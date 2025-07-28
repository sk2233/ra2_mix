package main

import (
	"bytes"
	"fmt"
	"io"
)

type VxlHeader struct {
	Magic     [16]byte // Voxel Animation
	Unused1   [4]byte
	LimbCount uint32
	Unused2   [4]byte
	BodySize  uint32
	Unused3   [770]byte
}

type LimbHeader struct {
	Name   [16]byte
	Unused [12]byte
}

type LimbFooter struct {
	Offset                             uint32
	Unused1                            [8]byte
	Scale                              float32
	Unused2                            [48]byte
	MinX, MinY, MinZ, MaxX, MaxY, MaxZ float32
	SizeX, SizeY, SizeZ                uint8
	NormalType                         uint8
}

type LimbItem struct {
	Color  uint8
	Normal uint8
}

type Limb struct {
	Items [][][]*LimbItem // x y z 对应位置可能为 nil
}

type Vxl struct {
	Header *VxlHeader
	// 一一对应的
	LimbHeaders []*LimbHeader
	Limbs       []*Limb
	LimbFooters []*LimbFooter
}

func ParseVxl(data []byte) *Vxl {
	reader := bytes.NewReader(data)
	header := ReadAny[VxlHeader](reader)  // 读取 VxlHeader
	limbHeaders := make([]*LimbHeader, 0) // 读取 LimbHeader
	for i := 0; i < int(header.LimbCount); i++ {
		limbHeaders = append(limbHeaders, ReadAny[LimbHeader](reader))
	} // 读取LimbFooter
	_, err := reader.Seek(int64(header.BodySize), io.SeekCurrent)
	HandleErr(err)
	limbFooters := make([]*LimbFooter, 0)
	for i := 0; i < int(header.LimbCount); i++ {
		limbFooters = append(limbFooters, ReadAny[LimbFooter](reader))
	}
	// 读取 Limb
	limbs := make([]*Limb, 0)
	for _, footer := range limbFooters {
		_, err = reader.Seek(int64(802+28*header.LimbCount+footer.Offset), io.SeekStart)
		HandleErr(err)
		limbs = append(limbs, ParseLimb(reader, footer))
	}
	return &Vxl{
		Header:      header,
		LimbHeaders: limbHeaders,
		Limbs:       limbs,
		LimbFooters: limbFooters,
	}
}

func ParseLimb(reader *bytes.Reader, footer *LimbFooter) *Limb {
	temp, err := reader.Seek(0, io.SeekCurrent)
	HandleErr(err)
	fmt.Println(temp)
	// 读取开始位置
	starts := make([][]uint32, footer.SizeY) // y x
	for y := 0; y < len(starts); y++ {
		starts[y] = make([]uint32, footer.SizeX)
		for x := 0; x < len(starts[y]); x++ {
			starts[y][x] = ReadU32(reader)
		}
	}
	temp, err = reader.Seek(0, io.SeekCurrent)
	HandleErr(err)
	fmt.Println(temp)
	// 获取数据开始位置，下面偏移使用
	dataStart, err := reader.Seek(4*int64(footer.SizeX)*int64(footer.SizeY), io.SeekCurrent)
	HandleErr(err)
	// 初始化容器
	items := make([][][]*LimbItem, footer.SizeX) // x y z
	for x := 0; x < len(items); x++ {
		items[x] = make([][]*LimbItem, footer.SizeY)
		for y := 0; y < len(items[x]); y++ {
			items[x][y] = make([]*LimbItem, footer.SizeZ)
		}
	}
	// 处理每一列
	for y := 0; y < len(starts); y++ {
		for x := 0; x < len(starts[y]); x++ {
			if starts[y][x] == 0xFFFFFFFF { // 跳过空列
				continue
			}
			_, err = reader.Seek(dataStart+int64(starts[y][x]), io.SeekStart)
			HandleErr(err)
			z := uint8(0)
			for z < footer.SizeZ {
				z += ReadU8(reader)     // 跳过多少位置
				count := ReadU8(reader) // 有多少个连续方块
				for i := 0; i < int(count); i++ {
					items[x][y][z] = &LimbItem{
						Color:  ReadU8(reader),
						Normal: ReadU8(reader),
					}
					z++
				}
				ReadU8(reader) // 这里还有一个重复的 count 直接扔掉
			}
		}
	}
	return &Limb{
		Items: items,
	}
}
