/*
* @Author: Yajun
* @Date:   2021/11/10 18:24
 */

package matrix

import (
	"gonum.org/v1/gonum/mat"
)

func DenseMean(m *mat.Dense, axis int) *mat.VecDense {
	r, c := m.Dims()
	var res *mat.VecDense
	if axis == 0 {
		res = mat.NewVecDense(c, nil)
		for i := 0; i < r; i++ {
			res.AddVec(res, m.RowView(i))
		}
		res.ScaleVec(1/float64(r), res)
		return res
	}
	if axis == 1 {
		res = mat.NewVecDense(r, nil)
		for i := 0; i < c; i++ {
			res.AddVec(res, m.ColView(i))
		}
		res.ScaleVec(1/float64(c), res)
		return res
	}
	panic(ErrInvalidArgument("axis", axis))
}

func DenseAddVector(m *mat.Dense, vec mat.Vector, axis int) *mat.Dense {
	r, c := m.Dims()

	if axis == 0 {
		var row *mat.VecDense
		for i := 0; i < r; i++ {
			row = m.RowView(i).(*mat.VecDense)
			row.AddVec(row, vec)
		}
		return m
	}

	if axis == 1 {
		var col *mat.VecDense
		for i := 0; i < c; i++ {
			col = m.ColView(i).(*mat.VecDense)
			col.AddVec(col, vec)
		}
		return m
	}
	panic(ErrInvalidArgument("axis", axis))
}

func DenseSubVector(m *mat.Dense, vec mat.Vector, axis int) *mat.Dense {
	r, c := m.Dims()

	if axis == 0 {
		var row *mat.VecDense
		for i := 0; i < r; i++ {
			row = m.RowView(i).(*mat.VecDense)
			row.SubVec(row, vec)
		}
		return m
	}

	if axis == 1 {
		var col *mat.VecDense
		for i := 0; i < c; i++ {
			col = m.ColView(i).(*mat.VecDense)
			col.SubVec(col, vec)
		}
		return m
	}
	panic(ErrInvalidArgument("axis", axis))
}

func DenseSubScala(m *mat.Dense, num float64) *mat.Dense {
	m.Apply(func(i, j int, v float64) float64 { return v + num }, m)
	return m
}
