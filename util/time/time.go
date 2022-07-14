package time

import "time"

func TSO() int64 {
	return (time.Now().UnixNano() / int64(time.Millisecond)) << 18
}

func CheckPointTSO() uint64 {
	return uint64(TSO() - 1)
}
