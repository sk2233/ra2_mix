/*
@author: sk
@date: 2025/5/4
*/
package main

import (
	"encoding/binary"
	"io"
	"math"
)

func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func ReadU8(reader io.Reader) uint8 {
	return ReadBytes(reader, 1)[0]
}

func ReadU16(reader io.Reader) uint16 {
	res := ReadBytes(reader, 2)
	return binary.LittleEndian.Uint16(res)
}

func ReadU32(reader io.Reader) uint32 {
	res := ReadBytes(reader, 4)
	return binary.LittleEndian.Uint32(res)
}

func ReadF32(reader io.Reader) float32 {
	temp := ReadU32(reader)
	return math.Float32frombits(temp)
}

func ReadAny[T any](reader io.Reader) *T {
	res := new(T)
	err := binary.Read(reader, binary.LittleEndian, res)
	HandleErr(err)
	return res
}

func ReadBytes(reader io.Reader, count int) []byte {
	res := make([]byte, count)
	_, err := reader.Read(res)
	HandleErr(err)
	return res
}

func ReadStr(reader io.Reader) string {
	bs := make([]byte, 0)
	for {
		temp := ReadU8(reader)
		if temp == 0 {
			return string(bs)
		} else {
			bs = append(bs, temp)
		}
	}
}
