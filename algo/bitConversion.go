package main

import (
	"fmt"
	"strconv"
)

func hexConversion(hex string) int64 {
	val, _ := strconv.ParseInt(hex, 16, 64)
	return val
}

func main() {
	var hex string
	for {
		fmt.Scan(&hex)
		hex = string([]byte(hex)[2:])
		fmt.Println(hexConversion(hex))
	}
	// hex := "0xAA"

}
