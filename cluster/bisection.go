/*
* @Author: Yajun
* @Date:   2021/11/12 16:44
 */

package cluster

import (
	"math"

	"github.com/yinyajun/golearn/utils"
	"gonum.org/v1/gonum/mat"
)

const (
	RatioCut = "ratio_cut"
	NCut     = "n_cut"
)

type SpecBisection struct {
	Similarities *mat.SymDense // 数据点之间的相似度矩阵
	MinCut       bool          // 是否是min-cut
	CutType      string        // ratioCut or nCut
	Verbose      bool          // 冗余模式
	Strict       bool          // 严格二分
	labels       []bool        // 二分结果
	major        bool          // 个数最多的是哪类
	done         bool
}

func NewSpecBisection(sim *mat.SymDense) *SpecBisection {
	return &SpecBisection{
		Similarities: sim,
		MinCut:       true,
		CutType:      RatioCut,
	}
}

func (c *SpecBisection) Fit(X mat.Symmetric) error {
	return c.partialFit(X)
}

func (c *SpecBisection) HasFitted() bool { return c.done }

func (c *SpecBisection) check() {
	if c.CutType != RatioCut && c.CutType != NCut {
		panic(ErrInvalidArgument)
	}
}

func (c *SpecBisection) partialFit(sim mat.Symmetric) error {
	var (
		fac  *SpecFactorize
		vec  *mat.VecDense
		tNum int // 标记为true的类中元素个数
	)

	switch c.CutType {
	case RatioCut:
		fac = NewSpecFactorize(c.Verbose, false)
	case NCut:
		fac = NewSpecFactorize(c.Verbose, true)
	default:
		panic(ErrInvalidArgument)
	}
	if err := fac.Fit(sim); err != nil {
		return err
	}

	k := utils.If(c.MinCut, 2, -2).(int)
	vec = fac.SmallNthEigenVector(k)

	c.labels = make([]bool, vec.Len())
	for i := 0; i < len(c.labels); i++ {
		if vec.AtVec(i) > 0 {
			c.labels[i] = true
			tNum++
			continue
		}
	}
	c.major = utils.If(tNum >= vec.Len()-tNum, true, false).(bool)

	if c.Strict {
		c.major = c.balance(tNum, true)
	}

	c.done = true
	return nil
}

func (c *SpecBisection) balance(num int, class bool) bool {
	var (
		n        = len(c.labels)
		another  = n - num
		majority = utils.If(num > another, class, !class).(bool)
		cnt      = utils.If(majority, (num-another)/2, (another-num)/2).(int)
	)

	if cnt == 0 {
		return majority
	}
	var (
		inc, dec float64
		delta    float64
		choose   int
	)

	for t := 0; t < cnt; t++ {
		for i := 0; i < n; i++ { // 遍历majority
			if c.labels[i] != majority {
				continue
			}
			inc, dec = 0, 0
			delta = utils.If(c.MinCut, math.Inf(1), math.Inf(-1)).(float64)
			for j := 0; j < n; j++ {
				if i == j {
					continue
				}
				if c.labels[j] == majority {
					inc += c.Similarities.At(i, j)
				} else {
					dec += c.Similarities.At(i, j)
				}
			}
			if c.MinCut {
				if inc-dec < delta {
					delta = inc - dec
					choose = i
				}
			} else {
				if inc-dec > delta {
					delta = inc - dec
					choose = i
				}
			}
		}
		c.labels[choose] = !majority
	}
	return majority
}

func (c *SpecBisection) Labels() []bool { return c.labels }
