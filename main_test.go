package main

import (
	"reflect"
	"testing"
)

func partition(nums []int, parts int) [][]int {
	size := len(nums) % parts

	var result [][]int
	for len(nums) > 0 {
		var newOne []int
		copy(newOne, nums[0:size])
		result = append(result, newOne)
		nums = nums[size:]
	}
	return result
}

func TestPartition(t *testing.T) {
	parts := partition([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 4)
	if !reflect.DeepEqual(parts[0], []int{1, 2}) {
		t.Errorf("expected part 0 to be [%v], got [%v]", []int{1, 2}, parts[0])
	}
}
