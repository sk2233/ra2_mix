/*
@author: sk
@date: 2025/5/4
*/
package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	//TestShp()
	//TestPal()
	//TestPcx()
	//return
	bs, err := os.ReadFile("res/ra2.mix")
	HandleErr(err)
	dfs(bs, "")
}

func dfs(bs []byte, base string) {
	res := ParseMix(bs)
	for key, val := range res {
		path := base + key
		if strings.HasSuffix(key, ".mix") {
			dfs(val, path+"/")
		} else if strings.HasSuffix(key, ".vxl") {
			// 用来找各种素材
			fmt.Println(path)
		}
	}
}
