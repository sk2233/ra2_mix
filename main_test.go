/*
@author: sk
@date: 2025/5/4
*/
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestInt(t *testing.T) {
	num := int32(-2233)
	fmt.Println(num, uint32(num))
}

func TestJson(t *testing.T) {
	m := make(map[int32]string)
	bs, err := os.ReadFile("/Users/bytedance/Documents/go/ra2mix/global_mix_database.json")
	HandleErr(err)
	json.Unmarshal(bs, &m)
	res := make(map[uint32]string)
	for k, v := range m {
		res[uint32(k)] = v
	}
	bs, err = json.Marshal(res)
	HandleErr(err)
	os.WriteFile("/Users/bytedance/Documents/go/ra2mix/mix_database.json", bs, 0644)
}

//func TestReadData(t *testing.T) {
//	file, err := os.Open("/Users/wepie/Documents/github/ra2_mix/res/global mix database.dat")
//	HandleErr(err)
//	count := ReadU32(file)
//	for i := 0; i < int(count); i++ {
//		temp := ReadStr(file)
//		fmt.Println(temp)
//	}
//}
