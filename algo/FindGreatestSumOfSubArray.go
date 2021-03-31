package main

import "fmt"

func FindGreatestSumOfSubArray(array []int) int {
	n := len(array)
	if n <= 0 {
		return 0
	}
	greatest := array[0]
	buf := make([][]int, n)
	for i, val := range array {
		buf[i] = []int{val}
	}
	left, right := 0, 1
	var tmp int
	for left <= right && left < n-1 {
		tmp = buf[left][right-left-1] + array[right]
		buf[left] = append(buf[left], tmp)
		if tmp > greatest {
			greatest = tmp
		}
		right++
		if right == n {
			left++
			right = left + 1
		}
	}
	return greatest
}

func main() {
	array := []int{-2, -8, -1, -5, -9}
	max := FindGreatestSumOfSubArray(array)
	fmt.Println(max)
}
