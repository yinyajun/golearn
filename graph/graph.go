/*
* @Author: Yajun
* @Date:   2021/12/7 23:15
 */

package graph

type Vertex interface{}

type BiGraph struct {
	Vertices []Vertex
	X, Y     []int
}

type Graph struct {
	Vertices []Vertex
}
