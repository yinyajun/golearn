/*
* @Author: Yajun
* @Date:   2021/12/7 23:29
 */

package utils

import "gonum.org/v1/gonum/mat"

type ID string

type VectorSpace interface {
	Size() int
	Query(ID) mat.Vector
	AddWithID(ID, mat.Vector)
	Deserialize(map[ID]string)
}
