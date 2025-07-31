package main

import (
	"fmt"
	"testing"
)

func TestDat(t *testing.T) {
	data := ReadData("ra2.mix", "local mix database.dat")
	dat := ParseDat(data)
	fmt.Println(dat.Summary)
	for key, val := range dat.Items {
		fmt.Println(key, val)
	}
}
