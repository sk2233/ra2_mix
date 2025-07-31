/*
@author: sk
@date: 2025/5/4
*/
package main

import (
	"encoding/binary"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"io"
	"math"
	"os"
	"strings"
)

const (
	BasePath = "/Users/wepie/Documents/github/ra2_mix/res/"
)

func ReadData(file, path string) []byte {
	bs, err := os.ReadFile(BasePath + file)
	HandleErr(err)
	items := strings.Split(path, "/")
	for _, item := range items {
		bs = ParseMix(bs)[item]
	}
	return bs
}

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

type ImageShow struct {
	option *ebiten.DrawImageOptions
	image  *ebiten.Image
}

func NewImageShow(img image.Image) *ImageShow {
	return &ImageShow{
		option: &ebiten.DrawImageOptions{},
		image:  ebiten.NewImageFromImage(img),
	}
}

func (i *ImageShow) Update() error {
	return nil
}

func (i *ImageShow) Draw(screen *ebiten.Image) {
	screen.DrawImage(i.image, i.option)
}

func (i *ImageShow) Layout(w, h int) (int, int) {
	return w, h
}

func ShowImage(img image.Image) {
	bound := img.Bounds()
	ebiten.SetWindowSize(bound.Dx(), bound.Dy())
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	err := ebiten.RunGame(NewImageShow(img))
	HandleErr(err)
}
