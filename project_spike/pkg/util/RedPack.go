package util

import (
	"fmt"
	"math/rand"
	"strconv"
)

func randFloats(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// Decimal 将一个浮点数保留两位小数
func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

func SplitRedPacket(amount float64, number int32) []float64 {
	packet := make([]float64, number)
	// 为每一个红包预分配0.01元，然后将剩余的金额放入红包池中
	pool := amount - float64(number)*0.01
	// 将红包转换成以分为单位的数字,微信抢红包发送的最多就也就是小数点后两位，所以到这里的时候，必定是一个整数了
	pool = pool * 100
	// 使用二倍均值法为每一个红包分配随机的金额
	var i int32 = 0
	for ; i < number; i++ {
		if i == number-1 {
			// 如果是最后一个红包，则将pool内的剩余金额全部分配
			p := Decimal(pool*0.01 + 0.01)
			packet[i] = p
			continue
		}
		avg := pool * 2.0 / float64(number-i)
		p := Decimal(randFloats(0.0, avg))
		pool = pool - p
		packet[i] = Decimal(p*0.01 + 0.01)
	}
	return packet
}
