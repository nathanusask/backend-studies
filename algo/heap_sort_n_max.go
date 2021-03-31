package main

import "fmt"

func maxHeapify(arr []int, cur int, bound int) {
	largest := cur
	left := 2*cur + 1
	right := 2*cur + 2
	if left < bound && arr[largest] < arr[left] {
		largest = left
	}
	if right < bound && arr[largest] < arr[right] {
		largest = right
	}

	if largest != cur {
		swap(&arr[largest], &arr[cur])
		maxHeapify(arr, largest, bound)
	}
}

func swap(a *int, b *int) {
	temp := *a
	*a = *b
	*b = temp
}

func main() {
	arr := []int{799, 10, 15, 2, 5, 6, 100, 51}
	size := len(arr)
	n := 3
	result := make([]int, n)

	for i := size/2 - 1; i >= 0; i-- {
		maxHeapify(arr, i, size)
	}

	for i := size - 1; i > size-1-n; i-- {
		swap(&arr[0], &arr[i])
		result[size-i-1] = arr[i]
		maxHeapify(arr, 0, i)
	}

	fmt.Println(result)
}
