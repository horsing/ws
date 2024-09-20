// Package utils contains all kinds of utilities
package utils

func Or(args ...string) string {
	for _, v := range args {
		if len(v) > 0 {
			return v
		}
	}
	return ""
}
