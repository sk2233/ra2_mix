/*
@author: sk
@date: 2025/5/4
*/
package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	TestShp()
	//TestPal()
	//TestPcx()
	return
	vxls = make([]string, 0)
	pals = make([]string, 0)
	bs, err := os.ReadFile("res/ra2.mix")
	HandleErr(err)
	dfs(bs, "")
	sort.Strings(pals)
	fmt.Println(vxls, pals)
}

var (
	vxls = make([]string, 0)
	pals = make([]string, 0)
)

func dfs(bs []byte, base string) {
	res := ParseMix(bs)
	for key, val := range res {
		path := base + key
		if strings.HasSuffix(key, ".mix") {
			dfs(val, path+"/")
		} else if strings.HasSuffix(key, ".vxl") {
			vxls = append(vxls, key)
		} else if strings.HasSuffix(key, ".pal") {
			pals = append(pals, key)
		}
	}
}
