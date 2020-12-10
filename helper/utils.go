// @Author: abbeymart | Abi Akindele | @Created: 2020-12-09 | @Updated: 2020-12-09
// @Company: mConnect.biz | @License: MIT
// @Description: crud utility functions

package helper

func ArrayStringContains(arr []string, val string) bool {
	for _, a := range arr {
		if a == val {
			return true
		}
	}
	return false
}

func ArrayIntContains(arr []int, val int) bool {
	for _, a := range arr {
		if a == val {
			return true
		}
	}
	return false
}