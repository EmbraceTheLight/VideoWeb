package Utilities

import "fmt"

// RoundOff 四舍五入函数
func RoundOff(f float64) int64 {
	var decimal = f - float64(int(f))
	if decimal >= 0.5 {
		return int64(f) + 1
	}
	return int64(f)
}

func SecondToTime(second int64) string {
	hour := second / 3600
	minute := (second % 3600) / 60
	second = second % 60
	return fmt.Sprintf("%02d:%02d:%02d", hour, minute, second)
}
