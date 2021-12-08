/*
* @Author: Yajun
* @Date:   2021/11/12 20:17
 */

package cluster

import (
	"fmt"
	"github.com/yinyajun/golearn/utils"
	"math/rand"
	"testing"

	"github.com/yinyajun/golearn/matrix"
	"gonum.org/v1/gonum/mat"
)

func TestSpecBisection(t *testing.T) {
	points := mat.NewDense(300, 64, nil)
	for i := 0; i < 300; i++ {
		for j := 0; j < 64; j++ {
			points.Set(i, j, rand.Float64())
		}
	}
	dist := matrix.DistFunc(func(i, j int) float64 {
		return utils.CosineSim(points.RowView(i), points.RowView(j))
	})

	d := matrix.Distances{
		Dist:        dist,
		Filters:     nil,
		NGoroutines: 8,
	}

	sim := d.SelfCartesian(utils.Range(0, 30, 1))
	b := NewSpecBisection(sim)
	b.Strict = true
	err := b.Fit(sim)
	fmt.Println(err)
	fmt.Println(b.Labels())
	fmt.Println(b.major)
}

func BenchmarkSpecBisection(b *testing.B) {
	points := mat.NewDense(300, 64, nil)
	for i := 0; i < 300; i++ {
		for j := 0; j < 64; j++ {
			points.Set(i, j, rand.Float64())
		}
	}
	dist := matrix.DistFunc(func(i, j int) float64 {
		return utils.InnerProduct(points.RowView(i), points.RowView(j))
	})

	d := matrix.Distances{
		Dist:        dist,
		Filters:     nil,
		NGoroutines: 16,
	}
	sim := d.SelfCartesian(utils.Range(0, 50, 1))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := NewSpecBisection(sim)
		s.Strict = true
		_ = s.Fit(sim)
		//fmt.Println(err)
		//fmt.Println(s.Labels())
		//fmt.Println(s.major)
	}

}
