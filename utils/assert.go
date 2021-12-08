/*
* @Author: Yajun
* @Date:   2021/11/9 23:11
 */

package utils

import (
	"fmt"
)

func Assert(condition bool, args ...interface{}) {
	if !condition {
		if msg := fmt.Sprint(args...); msg != "" {
			fmt.Printf("Assert failed, %s\n", msg)
		} else {
			fmt.Println("Assert failed")
		}
	}
}

func Assertf(condition bool, format string, args ...interface{}) {
	if !condition {
		if msg := fmt.Sprintf(format, args...); msg != "" {
			fmt.Printf("Assert failed, %s\n", msg)
		} else {
			fmt.Println("Assert failed")
		}
	}
}

func AssertFunc(fn func() error) {
	if err := fn(); err != nil {
		fmt.Printf("AssertFunc failed: %v", err)
	}
}

func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}
