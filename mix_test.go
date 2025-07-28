package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
)

var (
	nameTypes = make(map[string]bool)
)

func TestParseMix(t *testing.T) {
	bs, err := os.ReadFile("/Users/wepie/Documents/github/ra2_mix/res/ra2.mix")
	HandleErr(err)
	dfsMix(bs)
	bs, err = os.ReadFile("/Users/wepie/Documents/github/ra2_mix/res/ra2md.mix")
	HandleErr(err)
	dfsMix(bs)
	fmt.Println(nameTypes)
}

func dfsMix(bs []byte) {
	confBs, err := os.ReadFile("/Users/wepie/Documents/github/ra2_mix/mix_database.json")
	HandleErr(err)
	cof := make(map[string]string)
	err = json.Unmarshal(confBs, &cof)
	HandleErr(err)
	res := ParseMix(bs)
	for key, val := range res {
		if strings.HasSuffix(key, ".mix") {
			dfsMix(val)
		}
		idx := strings.LastIndex(key, ".")
		if idx >= 0 {
			nameTypes[key[idx:]] = true
			// 1 2 4 8
			if key[idx:] == ".fnt" {
				//temp := ParseHva(val)
				//fmt.Println(temp)
				//os.WriteFile("/Users/wepie/Documents/github/ra2_mix/res/temp.hva", val, 0644)
			}
		} else {
			fmt.Println(key)
		}
	}
}
