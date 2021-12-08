/*
* @Author: Yajun
* @Date:   2021/11/5 23:00
 */

package matrix

import (
	"github.com/yinyajun/golearn/utils"
	"gonum.org/v1/gonum/mat"
	"math/rand"
	"testing"
)

func TestDistances_SysCartesian(t *testing.T) {
	points := mat.NewDense(10, 4, nil)
	for i := 0; i < 10; i++ {
		for j := 0; j < 4; j++ {
			points.Set(i, j, rand.Float64())
		}
	}

	dist := DistFunc(func(i, j int) float64 {
		return utils.InnerProduct(points.RowView(i), points.RowView(j))
	})

	d := &Distances{
		Dist:        dist,
		Filters:     nil,
		NGoroutines: 8,
	}
	d.Filters = append(d.Filters, &KNNFilter{Typ: AnyKNN, K: 4})

	res := d.SelfCartesian([]int{0, 1, 2, 3, 4, 5})

	showMatrix(res)
	res2 := SubCartesian(res, []int{0, 1, 5}, []int{2, 4, 3})
	showMatrix(res2)
}
