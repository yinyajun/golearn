/*
* @Author: Yajun
* @Date:   2021/11/12 15:52
 */

package matrix

import "fmt"

func ErrInvalidArgument(name string, data interface{}) string {
	return fmt.Sprintf("Invalid argument: %s=%v", name, data)
}
