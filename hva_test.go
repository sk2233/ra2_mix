package main

import (
	"fmt"
	"testing"
)

func TestHva(t *testing.T) {
	data := ReadData("ra2.mix", "local.mix/hmec.hva")
	hva := ParseHva(data)
	fmt.Println(hva)
}
