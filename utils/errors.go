/*
* @Author: Yajun
* @Date:   2021/11/10 22:16
 */

package utils

import "fmt"

func ErrInvalidArgument(name string, data interface{}) string {
	return fmt.Sprintf("Invalid argument: %s=%v", name, data)
}
