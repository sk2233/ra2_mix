package main

import (
	"bytes"
	"hash/crc32"
	"strings"
)

type Dat struct {
	Summary string
	Items   map[uint32]string // id -> name
}

func ParseDat(data []byte) *Dat {
	summary := ReadStr(bytes.NewReader(data[:0x34]))
	items := make(map[uint32]string)
	reader := bytes.NewReader(data[0x34:])
	for reader.Len() > 0 {
		name := ReadStr(reader)
		items[HashName(name)] = name
	}
	return &Dat{Summary: summary, Items: items}
}

// 使用 CRC32 获取 id 值
func HashName(name string) uint32 {
	name = strings.ToUpper(name)
	pad := 0
	if len(name)%4 > 0 { // 对齐到 4
		pad = 4 - len(name)%4
	}
	padName := make([]byte, len(name)+pad) // 准备补齐
	copy(padName, name)
	nameLen := len(name) // 进行补齐
	roundLen := nameLen / 4 * 4
	if nameLen != roundLen {
		padName[nameLen] = byte(nameLen - roundLen)
		for p := 1; p < pad; p++ {
			padName[nameLen+p] = padName[roundLen]
		}
	}
	return crc32.ChecksumIEEE(padName)
}
