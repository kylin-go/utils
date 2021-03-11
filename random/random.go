package random

import (
	"math/rand"
	"time"
)

// 生成随机字符串
// l   : 随机字符串长度
// char: 生成随机字符的因子
func RandStr(l int, char string) string {
	rand.Seed(time.Now().UnixNano())
	if char == "" {
		char = "abcdefghijklmnopqrstuvwxy0123456789"
	}
	var d []byte
	for i := 0; i < l; i++ {
		d = append(d, char[rand.Intn(len(char))])
	}
	return string(d)
}

func RandInt(i int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(i)
}

func RandIntRange(i, j int) int {
	if i > j {
		n := i
		i = j
		j = n
	}
	j += 1
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(j)
	cha := j - i
	for {
		if num >= i && num <= j {
			break
		}
		num += cha
	}
	return num
}

func RandFloat() float32 {
	rand.Seed(time.Now().UnixNano())
	return rand.Float32()
}
