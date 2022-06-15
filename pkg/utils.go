package pkg

import "strings"

func SplitContainer(s, sep, substr string) bool {
	ss := strings.Split(s, sep)
	if len(ss) > 0 {
		m := make(map[string]bool)
		for _, v := range ss {
			m[v] = true
		}

		_, has := m[substr]
		return has
	}
	return false
}
