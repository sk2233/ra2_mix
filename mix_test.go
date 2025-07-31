package main

import (
	"fmt"
	"os"
	"testing"
)

func TestMix(t *testing.T) {
	bs, err := os.ReadFile(BasePath + "ra2.mix")
	HandleErr(err)
	temp := ParseMix(bs)
	for key, val := range temp {
		fmt.Println(key, len(val))
	}
}
