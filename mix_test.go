package main

import (
	"fmt"
	"os"
	"testing"
)

func TestMix(t *testing.T) {
	bs, err := os.ReadFile(BasePath + "language.mix")
	HandleErr(err)
	temp := ParseMix(bs)
	for key, val := range temp {
		fmt.Println(key, len(val))
	}
}
