/*
@author: sk
@date: 2025/5/4
*/
package main

type Color struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

// 调色板 #0 是透明色不渲染
// 阴影帧中所有非 #0 使用 #1 渲染

func ParsePal(data []byte) []*Color {
	res := make([]*Color, 0)
	for i := 0; i < len(data); i += 3 {
		res = append(res, &Color{AdjustColor(data[i]),
			AdjustColor(data[i+1]), AdjustColor(data[i+2])})
	}
	return res
}

func AdjustColor(val uint8) uint8 {
	// 虽然存储是 24bit 但是是按 16bit 用的，若需要按 24bit 使用需要进行一定缩放
	return uint8(int(val) * 255 / 63)
}
