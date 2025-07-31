/*
@author: sk
@date: 2025/5/4
*/
package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"slices"

	"golang.org/x/crypto/blowfish"
)

func Str2BigInt(val string) *big.Int {
	res := new(big.Int)
	res.SetString(val, 10)
	return res
}

func Byte2BigInt(val []byte) *big.Int {
	res := new(big.Int)
	slices.Reverse(val) // 使用小端序  这里破坏了 slice
	res.SetBytes(val)
	return res
}

func BigInt2Byte(val *big.Int) []byte {
	temp := val.FillBytes(make([]byte, (val.BitLen()+7)/8))
	slices.Reverse(temp) // 使用小端序
	return temp
}

func Pow(x, y, m *big.Int) *big.Int {
	return new(big.Int).Exp(x, y, m)
}

func DecryptCipherKey(data []byte) []byte {
	// pow 与 mod 是固定的
	pow := big.NewInt(65537)
	mod := Str2BigInt("681994811107118991598552881669230523074742337494683459234572860554038768387821901289207730765589")
	res := make([]byte, 0)
	// 先计算第一段
	base := Byte2BigInt(data[:40])
	temp := Pow(base, pow, mod)
	res = append(res, BigInt2Byte(temp)...)
	// 再计算第 2 段
	base = Byte2BigInt(data[40:])
	temp = Pow(base, pow, mod)
	res = append(res, BigInt2Byte(temp)...)
	return res
}

func DecryptData(cipher *blowfish.Cipher, data []byte, count int) []byte {
	res := make([]byte, count)
	for i := 0; i < count; i += 8 {
		cipher.Decrypt(res[i:], data[i:])
	}
	return res
}

func Align(val int) int { // 按 8 对齐 用于解码
	return (val + 8 - 1) / 8 * 8
}

type MixHeader struct {
	ID     uint32
	Offset uint32
	Size   uint32
}

func ParseMix(data []byte) map[string][]byte {
	flag := binary.LittleEndian.Uint32(data[:4])
	var count uint16
	var dataSize uint32
	var headerData, bodyData []byte
	if flag&0x20000 > 0 { // 加密了
		// 构建解码器
		key := DecryptCipherKey(data[4 : 4+80])
		cipher, err := blowfish.NewCipher(key)
		HandleErr(err)
		// 获取 count dataSize
		temp := DecryptData(cipher, data[4+80:], 8) // 只需要 6 个但是需要对齐到 8
		count = binary.LittleEndian.Uint16(temp[:2])
		dataSize = binary.LittleEndian.Uint32(temp[2:])
		// 获取 headerData bodyData
		headerData = DecryptData(cipher, data[4+80+8:], Align(int(count)*12)) // 每个子项目 12 byte
		headerData = append(temp[6:], headerData...)                          // temp 还有 2 个也要拿回来
		bodyData = data[4+80+Align(int(count)*12+6):]
	} else { // 没有加密
		count = binary.LittleEndian.Uint16(data[4:])
		dataSize = binary.LittleEndian.Uint32(data[4+2:])
		headerData = data[4+2+4 : 4+2+4+count*12]
		bodyData = data[4+2+4+count*12:]
	}
	fmt.Printf("mix data flag = %d , count = %d , dataSize = %d\n", flag, count, dataSize)
	// 获取 id -> name  可以缓存不用每次都读取
	id2Name := make(map[uint32]string)
	bs, err := os.ReadFile("res/mix_database.json")
	HandleErr(err)
	err = json.Unmarshal(bs, &id2Name)
	HandleErr(err)
	// 处理每一项
	reader := bytes.NewReader(headerData)
	res := make(map[string][]byte)
	for i := 0; i < int(count); i++ {
		header := ReadAny[MixHeader](reader)
		res[id2Name[header.ID]] = bodyData[header.Offset:min(int(header.Offset+header.Size), len(bodyData))]
	}
	return res
}
