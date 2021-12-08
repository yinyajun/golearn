/*
* @Author: Yajun
* @Date:   2021/11/10 20:32
 */

package cluster

import "errors"

const (
	ErrIndexOutOfRange = "index out of range"
	ErrInvalidArgument = "invalid argument"
	ErrEmptyInput      = "empty input"
	//ErrEigenFactorization = "eigen factorization fails"

)

var (
	ErrEigenFactorization = errors.New("eigen factorization fails")
	ErrFitHasNotDone      = errors.New("fit has not done")
)
