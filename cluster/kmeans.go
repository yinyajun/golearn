/*
* @Author: Yajun
* @Date:   2021/11/10 15:48
 */

package cluster

import (
	"log"
	"math"
	"math/rand"
	"runtime"
	"time"

	"github.com/yinyajun/golearn/matrix"
	"github.com/yinyajun/golearn/utils"
	"gonum.org/v1/gonum/mat"
)

const (
	Full  = "full"
	Elkan = "elkan"
)

type KMeans struct {
	NClusters   int    // 聚类数
	MaxIter     int    // 最大迭代次数
	NInit       int    // 聚类次数（因为kmeans可能陷入local minima， 多次聚类取最好的一次）
	Verbose     bool   // 冗余模式
	NGoroutines int    // 计算并发程度
	Algorithm   string // 采用算法 "full"原始EM方式；"elkan"
	centers     *mat.Dense
	labels      []int
	cost        float64
	unchanged   bool
	done        bool
}

func NewKMeans(NClusters int) *KMeans {
	m := &KMeans{
		NClusters:   NClusters,
		MaxIter:     50,
		NInit:       3,
		Verbose:     false,
		NGoroutines: runtime.NumCPU() / 2,
		Algorithm:   Full,
	}
	return m
}

func (m *KMeans) checkParams(points *mat.Dense) {
	if points.IsEmpty() {
		panic(ErrEmptyInput)
	}
	nSamples, _ := points.Dims()
	if m.NClusters <= 1 || m.NClusters >= nSamples {
		panic(ErrInvalidArgument)
	}
	if m.NGoroutines <= 1 || m.NGoroutines >= 100 {
		panic(ErrInvalidArgument)
	}
	if m.NInit < 1 || m.MaxIter < 1 {
		panic(ErrInvalidArgument)
	}
}

// Fit
// note that X.floatVal will be modified (subtract mean) during fit time
func (m *KMeans) Fit(X *mat.Dense) error {
	m.checkParams(X)
	return m.partialFit(X)
}

func (m *KMeans) partialFit(points *mat.Dense) error {
	var (
		seed        int64
		bestCenters *mat.Dense
		bestCost    = math.Inf(1)
		bestLabels  = make([]int, m.nSamples(points))
	)
	// subtract of mean of points for more accurate distance computations
	XMean := matrix.DenseMean(points, 0)
	matrix.DenseSubVector(points, XMean, 0)

	for i := 0; i < m.NInit; i++ {
		seed = time.Now().Unix()
		switch m.Algorithm {
		case Full:
			m.singleFull(points, seed)
		case Elkan:
			m.singleElkan(points, seed)
		}
		if m.cost < bestCost {
			copy(bestLabels, m.labels)
			bestCost = m.cost
			bestCenters = m.centers
		}
	}
	matrix.DenseAddVector(points, XMean, 0)
	matrix.DenseAddVector(bestCenters, XMean, 0)

	m.labels = bestLabels
	m.centers = bestCenters
	m.cost = bestCost
	m.done = true
	return nil
}

func (m *KMeans) HasFitted() bool { return m.done }

func (m *KMeans) Transform(X mat.Vector) int {
	var (
		dist             float64
		cluster, minDist = 0, math.Inf(1)
	)

	for k := 0; k < m.NClusters; k++ {
		dist = utils.EuclideanSquare(X, m.centers.RowView(k))
		if dist < minDist {
			cluster, minDist = k, dist
		}
	}
	return cluster
}

func (m *KMeans) singleFull(points *mat.Dense, seed int64) {
	m.initCenters(points, seed)
	for iter := 0; iter < m.MaxIter && !m.unchanged; iter++ {
		m.assign(points)
		m.update(points)
		if m.Verbose {
			log.Printf("[Epoch %d] Cost: %f\n", iter, m.cost)
		}
	}
}

// assign 将所有点分配到最近的聚类中心（类似于EM中的E步）
func (m *KMeans) assign(points *mat.Dense) {
	var (
		limit    = make(chan struct{}, m.NGoroutines)
		done     = make(chan struct{})
		converge = make(chan bool)
	)

	_assign := func(i int) bool {
		var dist float64
		cluster, minDist := 0, math.Inf(1)
		for k := 0; k < m.NClusters; k++ {
			dist = utils.EuclideanSquare(points.RowView(i), m.centers.RowView(k))
			if dist < minDist {
				cluster, minDist = k, dist
			}
		}
		if m.labels[i] == cluster {
			return true
		}
		m.labels[i] = cluster
		return false
	}

	go func() {
		m.unchanged = true
		for i := 0; i < m.nSamples(points); i++ {
			m.unchanged = <-converge && m.unchanged // note short circuit
		}
		done <- struct{}{}
	}()

	for i := 0; i < m.nSamples(points); i++ {
		limit <- struct{}{}
		go func(i int) {
			converge <- _assign(i)
			<-limit
		}(i)
	}
	<-done

	close(limit)
	close(converge)
	close(done)
}

// update 更新聚类中心（类似于EM中的M步）
func (m *KMeans) update(points *mat.Dense) {
	var (
		costs        = make(chan float64)
		_, nFeatures = points.Dims()
	)

	_update := func(k int) float64 {
		var (
			centroid = mat.NewVecDense(nFeatures, nil)
			cnt      int
			cost     float64
		)

		for x, class := range m.labels {
			if class != k {
				continue
			}
			centroid.AddVec(centroid, points.RowView(x))
			cnt++
		}
		centroid.ScaleVec(1/float64(cnt), centroid)

		m.centers.SetRow(k, centroid.RawVector().Data)
		for x, class := range m.labels {
			if class != k {
				continue
			}
			cost += utils.EuclideanSquare(centroid, points.RowView(x))
		}
		return cost
	}

	for k := 0; k < m.NClusters; k++ {
		p := k
		go func(k int) {
			costs <- _update(k)
		}(p)
	}

	m.cost = 0
	for k := 0; k < m.NClusters; k++ {
		m.cost += <-costs // todo: possible overflow
	}
	close(costs)
}

func (m *KMeans) singleElkan(points *mat.Dense, seed int64) {
	m.initCenters(points, seed)
	// todo: to be implemented
	panic("not implemented")
}

func (m *KMeans) Center(i int) mat.Vector { return m.centers.RowView(i) }

func (m *KMeans) Labels() []int { return m.labels }

func (m *KMeans) Cost() float64 { return m.cost }

func (m *KMeans) nSamples(points *mat.Dense) int {
	n, _ := points.Dims()
	return n
}

// initCenters 初始化聚类起点（使用kmeans++方式初始化起点，各个簇的起点相对分离）
func (m *KMeans) initCenters(points *mat.Dense, seed int64) {
	_, nFeatures := points.Dims()
	m.centers = mat.NewDense(m.NClusters, nFeatures, nil)
	m.labels = make([]int, m.nSamples(points))

	var (
		minDist float64
		centers = make([]int, m.NClusters)
		sampler = utils.NewSampler(m.nSamples(points))
		chosen  = make(map[int]struct{})
	)

	rand.Seed(seed)

	// 初始化第一个点
	centers[0] = rand.Intn(m.nSamples(points))
	chosen[centers[0]] = struct{}{}

	// 初始化其余点（计算每个数据点和已有簇中心之间的最短距离，生成该样本被选为聚类中心的概率）
	for k := 1; k < m.NClusters; k++ {
		for j := 0; j < m.nSamples(points); j++ {
			minDist = utils.EuclideanSquare(points.RowView(j), points.RowView(centers[0]))
			for q := 1; q < k; q++ {
				minDist = math.Min(minDist,
					utils.EuclideanSquare(points.RowView(j), points.RowView(centers[q])))
			}
			sampler.Assign(j, minDist) // note: 这里的距离应有平方含义
		}
		for {
			centers[k] = sampler.Sample()
			if _, exist := chosen[centers[k]]; !exist {
				chosen[centers[k]] = struct{}{}
				break
			}
		}
	}
	if m.Verbose {
		log.Printf("[Init] %v\n", centers)
	}

	for k := 0; k < m.NClusters; k++ {
		m.centers.SetRow(k, points.RawRowView(centers[k]))
	}
}
