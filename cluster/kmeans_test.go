/*
* @Author: Yajun
* @Date:   2021/11/10 22:14
 */

package cluster

import (
	"fmt"
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestKMeans_Fit(t *testing.T) {
	data := mat.NewDense(20, 2, nil)
	data.SetRow(0, []float64{2.59417075, 1.28887601})
	data.SetRow(1, []float64{1.92989795, -0.45664432})
	data.SetRow(2, []float64{-1.06288269, -2.40444104})
	data.SetRow(3, []float64{-1.65295978, -3.7959314})
	data.SetRow(4, []float64{3.34500752, 0.30877033})
	data.SetRow(5, []float64{-0.63380666, -4.41110015})
	data.SetRow(6, []float64{-1.16246674, -3.45067134})
	data.SetRow(7, []float64{-2.49421692, -3.08565762})
	data.SetRow(8, []float64{-1.03966052, -1.91926718})
	data.SetRow(9, []float64{3.3853239, 1.70518444})
	data.SetRow(10, []float64{2.41577648, 1.18270849})
	data.SetRow(11, []float64{-0.66843753, -1.67742351})
	data.SetRow(12, []float64{0.78536247, 0.17177281})
	data.SetRow(13, []float64{-0.14551976, -1.86339789})
	data.SetRow(14, []float64{-2.19893054, -2.49214272})
	data.SetRow(15, []float64{0.88416824, 2.47133519})
	data.SetRow(16, []float64{2.00691585, 1.24345064})
	data.SetRow(17, []float64{2.93187229, 1.74025265})
	data.SetRow(18, []float64{0.46027166, -3.95214669})
	data.SetRow(19, []float64{0.61090075, 1.6846158})

	k := NewKMeans(2)
	k.Verbose = true
	err := k.Fit(data)
	if err != nil {
		panic(err.Error())
	}

	clusters := []int{0, 1, 4, 9, 10, 12, 15, 16, 17, 19}
	var n int

	for _, i := range clusters {
		n += k.Labels()[i]
	}
	fmt.Println(k.Center(0))
	fmt.Println(k.Center(1))
	if n == 0 || n == 10 {
		return
	}
	t.Errorf("unexpected cluster result")

}
