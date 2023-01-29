package util

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func Code() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	vcode = vcode + "_" + strconv.FormatInt(time.Now().Unix(), 10)
	return vcode
}
