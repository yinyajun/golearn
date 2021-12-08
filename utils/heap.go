/*
* @Author: Yajun
* @Date:   2021/11/5 17:18
 */

package utils

// cur: n
// left child: 2*n+1; right child: 2*n+2
// parent: (n-1)/2

// K biggest放在index[:K]
func repartition(arr []float64, k int) []int {
	var (
		n     = len(arr)
		less  = func(i, j int) bool { return arr[i] < arr[j] }
		great = func(i, j int) bool { return arr[i] > arr[j] }
		index = Range(0, n, 1)
	)
	if k <= n/2 {
		// 用小顶堆求k biggest放在index[:K]
		Heapify(index[:k], great)
		for i := k; i < n; i++ {
			if great(index[i], index[0]) {
				index[0], index[i] = index[i], index[0]
				sink(index, 0, k-1, great)
			}
		}
	} else {
		// 用大顶堆求n-K smallest放在index[K:]
		if k == len(index) { // note: 必须保证index[K]合法
			return index
		}
		Heapify(index[k:], less)
		for i := 0; i < k; i++ {
			if less(index[i], index[k]) {
				index[k], index[i] = index[i], index[k]
				sink(index[k:], 0, n-k-1, less) // note 根节点的位置确定一个heap
			}
		}
	}
	return index
}

// KBiggest return K biggest element Index
func KBiggest(arr []float64, k int) []int {
	if k <= 0 {
		return []int{}
	}
	if k > len(arr) {
		k = len(arr)
	}
	return repartition(arr, k)[:k]
}

// KSmallest return K smallest element Index
// ( n-K biggest element Index)
func KSmallest(arr []float64, k int) []int {
	if k <= 0 {
		return []int{}
	}
	if k > len(arr) {
		k = len(arr)
	}
	return repartition(arr, len(arr)-k)[len(arr)-k:]
}

func Heapify(index []int, less func(i, j int) bool) {
	n := len(index)
	for i := (n - 1 - 1) / 2; i >= 0; i-- {
		sink(index, i, n-1, less)
	}
}

// 大顶堆的sink
func sink(index []int, lo, hi int, less func(i, j int) bool) {
	root := lo
	var m int

	for 2*root+1 <= hi {
		// find max child
		m = 2*root + 1
		if m+1 <= hi && less(index[m], index[m+1]) {
			m++
		}
		if !less(index[root], index[m]) { // parent >= max child
			break
		}
		index[m], index[root] = index[root], index[m]
		root = m
	}
}
