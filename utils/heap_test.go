/*
* @Author: Yajun
* @Date:   2021/11/5 18:31
 */

package utils

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"gonum.org/v1/gonum/floats"
)

func TestKBiggest(t *testing.T) {
	size := 10
	k := 6
	nums := make([]float64, size)
	rand.Seed(time.Now().Unix())
	for i := 0; i < len(nums); i++ {
		nums[i] = rand.Float64()
		fmt.Println(i, nums[i])
	}

	res := KBiggest(nums, k)
	fmt.Println(res)

	idx := Range(0, size, 1)
	floats.Argsort(nums, idx)
	fmt.Println(idx[(size - k):])
}

func TestKSmallest(t *testing.T) {
	size := 10
	k := 4
	nums := make([]float64, size)
	rand.Seed(time.Now().Unix())
	for i := 0; i < len(nums); i++ {
		nums[i] = rand.Float64()
		fmt.Println(i, nums[i])
	}

	res := KSmallest(nums, k)
	fmt.Println(res)

	idx := Range(0, size, 1)
	floats.Argsort(nums, idx)
	fmt.Println(idx[:k])
}
