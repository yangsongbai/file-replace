//File stringuitl.go
//@Author Songbai Yang
//@Date 2019/7/31
package tools

import (
	"fmt"
	"strconv"
	"strings"
)

func IntToString(num int) string {
	return strconv.Itoa(num)
}

//TrimSpaces will trim space and line break
func TrimSpaces(str string) string {
	return strings.TrimSpace(str)
}

func Contains(source, sub string) bool {
	return strings.Contains(source, sub)
}
func Index(source, sub string) int {
	return strings.Index(source, sub)
}

func ToString(str []string) string {
	return fmt.Sprint(str)
}

func StartWith(str, start string) bool {
	if str[:len(start)] == start {
		return true
	}
	return false
}

func Match(s, p string) bool {
	m := len(s)
	n := len(p)
	//var f [m+1][n+1]bool
	//f := make([][]bool,100,100)
	f := make([][]bool, m+1)
	for i := 0; i < m+1; i++ {
		f[i] = make([]bool, n+1)
	}

	f[0][0] = true
	for i := 1; i <= n; i++ {
		f[0][i] = f[0][i-1] && p[i-1] == '*'
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if s[i-1] == p[j-1] || p[j-1] == '?' {
				f[i][j] = f[i-1][j-1]
			}
			if p[j-1] == '*' {
				f[i][j] = f[i][j-1] || f[i-1][j]
			}
		}
	}

	return f[m][n]
}
