package core

import (
	"math"
	"math/rand"
	"strconv"
	"time"
)

func GenerateRandomString(t string) string {
	n := strconv.Itoa(int(math.Abs(float64(time.Now().Unix()) * rand.Float64() * 1e4)))

	o := 0

	for _, char := range n {
		digit, _ := strconv.Atoi(string(char))
		o += digit
	}

	o += len(n)
	oString := strconv.Itoa(o)
	for len(oString) < 3 {
		oString = "0" + oString
	}

	return t + n + oString
}
