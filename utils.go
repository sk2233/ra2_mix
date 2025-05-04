/*
@author: sk
@date: 2025/5/4
*/
package main

import (
	"bytes"
	"encoding/binary"
	"io"
)

func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func ReadU8(reader *bytes.Reader) uint8 {
	res, err := reader.ReadByte()
	HandleErr(err)
	return res
}

func ReadU16(reader io.Reader) uint16 {
	res := ReadBytes(reader, 2)
	return binary.LittleEndian.Uint16(res)
}

func ReadU32(reader io.Reader) uint32 {
	res := ReadBytes(reader, 4)
	return binary.LittleEndian.Uint32(res)
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
