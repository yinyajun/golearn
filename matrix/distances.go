/*
* @Author: Yajun
* @Date:   2021/12/7 14:23
 */

package matrix

import (
	"fmt"
	"github.com/yinyajun/golearn/utils"
	"gonum.org/v1/gonum/mat"
	"math"
)

const (
	AnyKNN = "any" // (i,j) in R || (j,i) in R
	AllKNN = "all" // (i,j) in R && (j,i) in R
)

type DistFunc func(i, j int) float64

type Distances struct {
	Dist        DistFunc
	Filters     []DistFilter
	NGoroutines int
}

// Cartesian returns {<x, y> | x in X && y in Y}
func (d *Distances) Cartesian(X, Y []int) *mat.Dense {
	var (
		limit = make(chan int, d.NGoroutines)
		res   = mat.NewDense(len(X), len(Y), nil)
	)
	// 计算距离矩阵
	for i, x := range X {
		for j, y := range Y {
			limit <- 1
			go func(i, j, x, y int) {
				res.Set(i, j, d.Dist(x, y))
				<-limit
			}(i, j, x, y)
		}
	}
	for i := 0; i < d.NGoroutines; i++ { // 确保最后一批goroutine完成job
		limit <- 1
	}
	close(limit)
	for _, f := range d.Filters {
		f.FilterDense(res)
	}
	return res
}

// SelfCartesian returns {<x, y>| x, y in Index}
func (d *Distances) SelfCartesian(Index []int) *mat.SymDense {
	var (
		limit = make(chan int, d.NGoroutines)
		res   = mat.NewSymDense(len(Index), nil)
	)
	// 计算距离矩阵
	for i, x := range Index {
		for j, y := range Index {
			if j <= i {
				continue
			}
			limit <- 1
			go func(i, j, x, y int) {
				res.SetSym(i, j, d.Dist(x, y))
				<-limit
			}(i, j, x, y)
		}
	}
	for i := 0; i < d.NGoroutines; i++ { // 确保最后一批goroutine完成job
		limit <- 1
	}
	close(limit)
	for _, f := range d.Filters {
		f.FilterSymmetric(res)
	}
	return res
}

// SubCartesian 从sim中提取子矩阵，特别注意s1,s2是sim的行列index的index
func SubCartesian(sim *mat.SymDense, s1, s2 []int) *mat.Dense {
	var (
		res = mat.NewDense(len(s1), len(s2), nil)
	)
	for i, x := range s1 {
		for j, y := range s2 {
			res.Set(i, j, sim.At(x, y))
		}
	}
	return res
}

type DistFilter interface {
	FilterSymmetric(*mat.SymDense)
	FilterDense(dense *mat.Dense)
}

func showMatrix(s mat.Matrix) {
	m, n := s.Dims()
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			fmt.Printf("%06f   ", s.At(i, j))
		}
		fmt.Println()
	}
	fmt.Println()
}

type KNNFilter struct {
	Typ string
	// any relationship: (i,j) in R || (j,i) in R
	// all relationship: (i,j) in R && (j,i) in R, De Morgan's Law: (p,q) not in R || (q,p) not in R
	K int
}

func (f *KNNFilter) FilterSymmetric(m *mat.SymDense) {
	var rmMark func(int, int)

	switch f.Typ {
	case AllKNN:
		rmMark = func(i, j int) {
			if m.At(i, j) < 0 {
				m.SetSym(i, j, 0)
			}
		}
	case AnyKNN:
		rmMark = func(i, j int) {
			if m.At(i, j) < 0 {
				m.SetSym(i, j, -m.At(i, j))
			} else {
				m.SetSym(i, j, 0)
			}
		}
	default:
		panic("Unsupported type for KNN filter")
	}

	var (
		r, c   = m.Dims()
		nums   = make([]float64, c)
		knnIdx []int
	)

	// 逐行标记(all:非K邻近, any:k邻近)
	for i := 0; i < r; i++ {
		// copy
		for j := 0; j < c; j++ {
			nums[j] = math.Abs(m.At(i, j))
		}
		//fmt.Println(i)
		//showVector(nums)

		if f.Typ == AllKNN {
			knnIdx = utils.KSmallest(nums, c-f.K)
		} else {
			knnIdx = utils.KBiggest(nums, f.K)
		}
		//fmt.Println(knnIdx)

		for _, j := range knnIdx {
			if m.At(i, j) > 0 {
				m.SetSym(i, j, -m.At(i, j)) // 符号标记
			}
		}
	}

	//showMatrix(m)

	// 清除标记
	for i := 0; i < r; i++ {
		for j := i + 1; j < c; j++ {
			rmMark(i, j)
		}
	}
}

func (f *KNNFilter) FilterDense(m *mat.Dense) {
	var rmMark func(int, int)

	switch f.Typ {
	case AllKNN:
		rmMark = func(i, j int) {
			if m.At(i, j) < 0 {
				m.Set(i, j, 0)
			}
		}
	case AnyKNN:
		rmMark = func(i, j int) {
			if m.At(i, j) < 0 {
				m.Set(i, j, -m.At(i, j))
			} else {
				m.Set(i, j, 0)
			}
		}
	default:
		panic("Unsupported type for KNN filter")
	}

	var (
		r, c   = m.Dims()
		nums   = make([]float64, c)
		knnIdx []int
	)

	// 逐行标记(all:非K邻近, any:k邻近)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			nums[j] = math.Abs(m.At(i, j))
		}
		//fmt.Println(i)
		//showVector(nums)

		if f.Typ == AllKNN {
			knnIdx = utils.KSmallest(nums, c-f.K)
		} else {
			knnIdx = utils.KBiggest(nums, f.K)
		}
		//fmt.Println(knnIdx)

		for _, j := range knnIdx {
			if m.At(i, j) > 0 {
				m.Set(i, j, -m.At(i, j)) // 符号标记
			}
		}
	}
	//showMatrix(m)

	// 清除标记
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			rmMark(i, j)
		}
	}
}
