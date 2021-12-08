/*
* @Author: Yajun
* @Date:   2021/11/10 22:53
 */

package cluster

import (
	"github.com/yinyajun/golearn/utils"
	"gonum.org/v1/gonum/mat"
)

// SpecClustering 谱聚类
type SpecClustering struct {
	Similarities *mat.SymDense // 数据点之间的相似度矩阵
	NClusters    int           // 聚类数
	ReducedDim   int           // 谱分解后的维度（ReducedDim > 0为优化min—cut，ReducedDim < 0为max-cut）
	CutType      string        // ratioCut or nCut
	Verbose      bool          // 冗余模式
	KMeans       *KMeans       // kMeans实例（降维后用kMeans再聚类）
	done         bool
}

func NewSpectralClustering(sim *mat.SymDense, NClusters int) *SpecClustering {
	c := &SpecClustering{
		Similarities: sim,
		NClusters:    NClusters,
		ReducedDim:   3 * NClusters,
		CutType:      "ratio_cut",
		KMeans:       NewKMeans(NClusters),
	}
	return c
}

func (c *SpecClustering) Fit(X mat.Symmetric) error {
	c.check()
	return c.partialFit(X)
}

func (c *SpecClustering) partialFit(sim mat.Symmetric) error {
	var (
		err     error
		reduced *mat.Dense
		fac     *SpecFactorize
	)
	switch c.CutType {
	case RatioCut:
		fac = NewSpecFactorize(c.Verbose, false)
	case NCut:
		fac = NewSpecFactorize(c.Verbose, true)
	default:
		panic(ErrInvalidArgument)
	}
	if err = fac.partialFit(sim); err != nil {
		return err
	}
	reduced = fac.SmallKEigenVectors(c.ReducedDim)
	if err = c.KMeans.partialFit(reduced); err != nil {
		return err
	}
	c.done = true
	return nil
}

func (c *SpecClustering) check() {
	n := c.Similarities.Symmetric()

	if c.NClusters > n || c.NClusters <= 1 {
		panic(ErrInvalidArgument)
	}

	// 1< abs(c.ReducedDim) <=n
	if c.ReducedDim < -n || (-1 <= c.ReducedDim && c.ReducedDim <= 1) ||
		c.ReducedDim > n {
		panic(ErrInvalidArgument)
	}
}

func (c *SpecClustering) HasFitted() bool { return c.done }

func (c *SpecClustering) Centers(i int) mat.Vector {
	utils.Assert(c.HasFitted(), ErrFitHasNotDone)
	return c.KMeans.Center(i)
}

func (c *SpecClustering) Labels() []int {
	utils.Assert(c.HasFitted(), ErrFitHasNotDone)
	return c.KMeans.Labels()
}
