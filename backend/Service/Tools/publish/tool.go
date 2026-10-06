package publish

import (
	"strconv"
	"time"
)

// Go 的无符号整型运算按 2^32 取模自动环绕，不会溢出 panic。
func hash32(s string) uint32 {
	var h uint32 = 5381
	for i := 0; i < len(s); i++ {
		h = h*33 + uint32(s[i])
	}
	return h
}

// pad36 把 base36 字符串左侧补 '0' 到 width 位；超过 width 则保留末尾 width 位。
func pad36(s string, width int) string {
	for len(s) < width {
		s = "0" + s
	}
	if len(s) > width {
		s = s[len(s)-width:]
	}
	return s
}

func MakeID(account string, ts ...int64) string {
	t := time.Now().UnixMilli()
	if len(ts) > 0 {
		t = ts[0]
	}
	timePart := pad36(strconv.FormatInt(t, 36), 8)
	accPart := pad36(strconv.FormatUint(uint64(hash32(account)), 36), 3)
	return timePart + accPart
}
