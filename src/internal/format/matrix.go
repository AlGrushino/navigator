package format

import (
	"fmt"
	"strings"
)

// Matrix formats an int matrix into aligned columns.
func Matrix(m [][]int) string {
	if len(m) == 0 {
		return ""
	}
	width := 1
	for i := range m {
		for j := range m[i] {
			w := len(fmt.Sprintf("%d", m[i][j]))
			if w > width {
				width = w
			}
		}
	}
	var b strings.Builder
	for i := range m {
		for j := range m[i] {
			if j > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(fmt.Sprintf("%*d", width, m[i][j]))
		}
		if i < len(m)-1 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}
