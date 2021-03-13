// 数值处理与数值转换
package numeral

import (
	"fmt"
	"strconv"
)

// 取小数点多少位
func FloatRound(f float64, n int) float64 {
	format := "%." + strconv.Itoa(n) + "f"
	res, _ := strconv.ParseFloat(fmt.Sprintf(format, f), 64)
	return res
}

// 浮点数转整数
// precision 浮点数转整数时是否进行精度运算(>=0.5时加1)
func Float2Int(f float64, precision bool) int64 {
	var i = int64(f)
	if !precision {
		return i
	}
	if f-float64(i) >= 0.5 {
		i += 1
	}
	return i
}

func Str2Float(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func Float2Str(f float64) string {
	return fmt.Sprintf("%f", f)
}

func Str2Int(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func Int2Str(s int64) string {
	return fmt.Sprintf("%d", s)
}

func Interface2Str(i interface{}) string {
	return fmt.Sprintf("%v", i)
}
