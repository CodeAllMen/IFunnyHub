package utils

import (
	"math/rand"
	"strconv"
	"time"
)

func NumToString(num int) string {
	num_str := strconv.Itoa(num)
	if len(num_str) < 4 {
		return num_str
	}
	remainder := len(num_str) % 3
	num_str_new := num_str[:remainder] + ","
	num_str_bak := num_str[remainder:]
	lable := 0
	for i := range num_str_bak {
		if lable == 3 {
			num_str_new += ","
			lable = 0
		}
		num_str_new += num_str_bak[i : i+1]
		lable++
	}
	return num_str_new
}

func Rand(num int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num_rand := r.Int31n(int32(num))
	return int(num_rand)
}
