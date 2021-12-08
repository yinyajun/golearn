/*
* @Author: Yajun
* @Date:   2021/11/10 20:48
 */

package utils

import (
	"math"

	"gonum.org/v1/gonum/mat"
)

type Metric func(a, b mat.Vector) float64

func Euclidean(a, b mat.Vector) (s float64) {
	if a.Len() != b.Len() {
		return
	}
	for i := 0; i < a.Len(); i++ {
		t := b.AtVec(i) - a.AtVec(i)
		s += t * t
	}
	s = math.Sqrt(s)
	return
}

func EuclideanSquare(a, b mat.Vector) (s float64) {
	if a.Len() != b.Len() {
		return
	}
	for i := 0; i < a.Len(); i++ {
		t := b.AtVec(i) - a.AtVec(i)
		s += t * t
	}
	return
}

func InnerProduct(a, b mat.Vector) (s float64) {
	if a.Len() != b.Len() {
		return
	}
	for i := 0; i < a.Len(); i++ {
		s += a.AtVec(i) * b.AtVec(i)
	}
	return
}

func JaccardSim(a, b mat.Vector) (s float64) {
	if a.Len() != b.Len() {
		return
	}
	var m11 int
	for i := 0; i < a.Len(); i++ {
		if a.AtVec(i) == b.AtVec(i) && a.AtVec(i) == 1 {
			m11 += 1
		}
	}
	return float64(m11) / float64(a.Len()-m11)
}

func CosineSim(a, b mat.Vector) float64 {
	if a.Len() != b.Len() {
		return 0
	}
	return InnerProduct(a, b) / (vecModule(a) * vecModule(b))
}

func Hamming(a, b mat.Vector) (s float64) {
	var n int
	n = a.Len()
	if b.Len() < n {
		n = b.Len()
	}
	for i := 0; i < n; i++ {
		if a.AtVec(i) != b.AtVec(i) {
			s += 1
		}
	}
	return
}

func Manhattan(a, b mat.Vector) (s float64) {
	if a.Len() != b.Len() {
		return
	}
	for i := 0; i < a.Len(); i++ {
		s += math.Abs(b.AtVec(i) - a.AtVec(i))
	}
	return
}

// RuzickaSim is weighted Jaccard similarity
func RuzickaSim(a, b mat.Vector) (s float64) {
	if a.Len() != b.Len() {
		return
	}
	var d float64
	for i := 0; i < a.Len(); i++ {
		s += math.Min(a.AtVec(i), b.AtVec(i))
		d += math.Max(a.AtVec(i), b.AtVec(i))
	}
	s = s / d
	return
}

func Pearson(a, b mat.Vector) (s float64) {
	if a.Len() != b.Len() {
		return
	}
	var xy, x, y, x2, y2 float64
	for i := 0; i < a.Len(); i++ {
		xy += a.AtVec(i) * b.AtVec(i)
		x += a.AtVec(i)
		y += b.AtVec(i)
		x2 += a.AtVec(i) * a.AtVec(i)
		y2 += b.AtVec(i) * b.AtVec(i)
	}
	return (float64(a.Len())*xy - xy) /
		(math.Sqrt(float64(a.Len())*x2-x*x) * math.Sqrt(float64(a.Len())*y2-y*y))
}

func vecModule(a mat.Vector) float64 {
	var ans float64
	for i := 0; i < a.Len(); i++ {
		ans += a.AtVec(i) * a.AtVec(i)
	}
	return math.Sqrt(ans)
}
