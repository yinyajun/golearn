/*
* @Author: Yajun
* @Date:   2021/11/10 23:11
 */

package cluster

import (
	"github.com/yinyajun/golearn/utils"

	"gonum.org/v1/gonum/mat"
)

type SpecFactorize struct {
	Norm    bool       // 拉普拉斯矩阵是否归一化
	Verbose bool       // 冗余模式
	eVal    []float64  // 特征值(升序)
	eVec    *mat.Dense // 特征向量
	done    bool
}

func NewSpecFactorize(verbose, norm bool) *SpecFactorize {
	return &SpecFactorize{Verbose: verbose, Norm: norm}
}

func (r *SpecFactorize) Fit(X mat.Symmetric) error {
	return r.partialFit(X)
}

func (r *SpecFactorize) partialFit(adj mat.Symmetric) (err error) {
	var (
		dim = adj.Symmetric()
		es  = mat.EigenSym{}
		L   mat.Symmetric
	)
	r.eVec = mat.NewDense(dim, dim, nil)
	r.eVal = make([]float64, dim)

	if r.Norm {
		L = NormedLaplacianMatrix(adj)
	} else {
		L = LaplacianMatrix(adj)
	}

	ok := es.Factorize(L, true)
	if !ok {
		return ErrEigenFactorization
	}
	es.Values(r.eVal)
	es.VectorsTo(r.eVec)
	r.done = true
	return
}

func (r *SpecFactorize) HasFitted() bool { return r.done }

// SmallKEigenVectors 获得L=D-A的最小k个特征向量（k>0）
func (r *SpecFactorize) SmallKEigenVectors(k int) *mat.Dense {
	// 注意这里分解的矩阵是L'= A-D，是半负定的
	// k > 0 : L'最大k个，L的最小k个eigen vector
	// k < 0 : L'最小k个，L的最大k个eigen vector
	// L'的最大的k个eigen value对应的eigen vector(=> L的最小k个eigen vector)
	utils.Assert(r.HasFitted(), ErrFitHasNotDone)
	n := len(r.eVal)
	// todo: check k
	lo, hi := n-k, n
	if k < 0 {
		lo, hi = 0, -k
	}
	return r.eVec.Slice(0, n, lo, hi).(*mat.Dense)
}

// SmallNthEigenVector 获得L=D-A的第N小特征向量（k>0）
func (r *SpecFactorize) SmallNthEigenVector(k int) *mat.VecDense {
	utils.Assert(r.HasFitted(), ErrFitHasNotDone)
	// todo:check k
	n := len(r.eVal)
	t := n - k
	if k < 0 {
		t = -k - 1
	}
	return r.eVec.Slice(0, n, t, t+1).(*mat.Dense).ColView(0).(*mat.VecDense)
}

func (r *SpecFactorize) EigenValues() []float64 {
	utils.Assert(r.HasFitted(), ErrFitHasNotDone)
	return r.eVal
}

func (r *SpecFactorize) EigenVectors() *mat.Dense {
	utils.Assert(r.HasFitted(), ErrFitHasNotDone)
	return r.eVec
}

// DegreeMatrix 从无向图的邻接矩阵计算度矩阵
func DegreeMatrix(m mat.Symmetric) *mat.DiagDense {
	n := m.Symmetric()
	d := make([]float64, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			d[i] += m.At(i, j)
		}
	}
	return mat.NewDiagDense(n, d)
}

// LaplacianMatrix 计算了L'矩阵
// L = D - W (semi-positive definite)
// L' = W - D (semi-negative definite)
func LaplacianMatrix(m mat.Symmetric) *mat.SymDense {
	var (
		t float64
	)
	n := m.Symmetric()
	L := mat.NewSymDense(n, nil)
	L.CopySym(m)
	D := DegreeMatrix(m)

	for i := 0; i < n; i++ {
		t = L.At(i, i) - D.At(i, i)
		L.SetSym(i, i, t)
	}
	return L
}

// NormedLaplacianMatrix 计算 L' = D^(-1/2) L' D^(-1/2)
func NormedLaplacianMatrix(m mat.Symmetric) *mat.SymDense {
	n := m.Symmetric()
	L := mat.NewSymDense(n, nil)
	D := DegreeMatrix(m)

	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			if i == j {
				L.SetSym(i, i, -1)
				continue
			}
			L.SetSym(i, j, m.At(i, j)/(D.At(i, i)*D.At(j, j)))
		}
	}
	return L
}
