package oracle

import "time"

func TSOracle() uint64 {
	return uint64((time.Now().UnixNano() / int64(time.Millisecond)) << 18)
}

func CheckPointTs() uint64 {
	return TSOracle() - 1
}
