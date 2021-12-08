/*
* @Author: Yajun
* @Date:   2021/11/10 20:25
 */

package utils

import (
	"math"
	"math/rand"
)

func Range(start, end, stride int) []int {
	n := int(math.Ceil(float64(end-start) / float64(stride)))
	ans := make([]int, n)
	for i, num := 0, start; i < n; i++ {
		ans[i] = num
		num += stride
	}
	return ans
}

type Sampler struct {
	cardinality int
	dist        []float64
	sum         float64
}

func NewSampler(c int) *Sampler {
	return &Sampler{
		cardinality: c,
		dist:        make([]float64, c),
	}
}

func (p *Sampler) Clear() {
	p.sum = 0
	for i := 0; i < len(p.dist); i++ {
		p.dist[i] = 0
	}
}

func (p *Sampler) Assign(i int, account float64) {
	if i >= p.cardinality {
		panic("[Sampler] assign i out of range")
	}
	p.dist[i] = account
	p.sum += account
}

func (p *Sampler) Sample() int {
	var (
		s, t float64
		k    int
	)
	k = 0
	t = rand.Float64() * p.sum
	for k = range p.dist {
		s += p.dist[k]
		if t <= s {
			break
		}
	}
	return k
}
