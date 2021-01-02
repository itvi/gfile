package util

import "fmt"

// ConvertByteToMB ...
func ConvertByteTo(n int64) string {
	switch {
	case n == 0:
		return ""
	case n < 1024*1024:
		return fmt.Sprintf("%.2f", float64(n)/1024) + "KB"
	case n < 1024*1024*1024:
		return fmt.Sprintf("%.2f", float64(n)/1024/1024) + "MB"
	default:
		return fmt.Sprintf("%.2f", float64(n)/1024/1024/1024) + "GB"
	}
}
