package main

import (
	"bytes"
	"github.com/ebitengine/oto/v3"
	"testing"
	"time"
)

func TestAud(t *testing.T) {
	data := ReadData("ra2.mix", "local.mix/intro.aud")
	aud := ParseAud(data)

	ctx, ready, err := oto.NewContext(&oto.NewContextOptions{
		SampleRate:   int(aud.Header.SampleRate),
		ChannelCount: aud.Header.GetChannelCount(),
		Format:       oto.FormatSignedInt16LE,
		BufferSize:   time.Second,
	})
	HandleErr(err)
	<-ready                                            // 等待初始化完成
	player := ctx.NewPlayer(bytes.NewReader(aud.Data)) // 创建播放器并写入PCM数据
	player.Play()
	for player.IsPlaying() { // 等待播放完成
	}
}
