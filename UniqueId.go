package utility

import (
	"time"
)

func UniqueId() string {
	buf := EncodeInt(uint64(GlobalRand1.Int63n(8999999999)+1000000000))
	buf = AppendInt(buf, uint64(time.Now().Unix() % (86488*365)))
	buf = AppendInt(buf, uint64(time.Now().Nanosecond()))
	buf = AppendInt(buf, uint64(GlobalRand2.Int63n(8999999999)+1000000000))
	return string(buf)
}
