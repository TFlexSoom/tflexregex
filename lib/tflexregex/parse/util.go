package parse

import (
	"fmt"
	"strconv"
)

func atoi(bs []byte) int {
	stringified := string(bs)
	res, err := strconv.Atoi(stringified)
	if err != nil {
		panic(fmt.Sprintf("could not convert decimal: %s | err: %v", stringified, err))
	}

	return res
}

func utf8(bs []byte) rune {
	r := uint32(0)
	currentIndex := 0
	current := bs[0]
	if current >= byte(0b1111_0000) {
		r += uint32(current) << 24
		currentIndex += 1
		current = bs[currentIndex]
	}

	if current >= byte(0b1110_0000) {
		r += uint32(current) << 16
		currentIndex += 1
		current = bs[currentIndex]
	}

	if current >= byte(0b1100_0000) {
		r += uint32(current) << 8
		r += uint32(bs[currentIndex+1])
		return rune(r)
	}

	if current >= byte(0b1000_0000) {
		panic("invalid unicode character detected")
	}

	return rune(current)
}
