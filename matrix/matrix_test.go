/*
* @Author: Yajun
* @Date:   2021/11/10 19:21
 */

package matrix

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
	"testing"
)

func TestDenseMean(t *testing.T) {
	m := mat.NewDense(3, 4, nil)
	m.SetRow(0, []float64{1, 2, 3, 4})
	m.SetRow(1, []float64{2, 2, 3, 4})
	m.SetRow(2, []float64{3, 2, 3, 4})
	a := DenseMean(m, 0)
	b := DenseMean(m, 1)
	fmt.Println(a)
	fmt.Println(b)

	DenseSubVector(m, a, 0)
	for i := 0; i < 3; i++ {
		fmt.Println(m.RawRowView(i))
	}

}
