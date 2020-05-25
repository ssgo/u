package u

import (
	"time"
)

const randMinNumber = 916132832
const randMaxNumber = 56800235583
const randNumber = randMaxNumber - randMinNumber

const randMinShortNumber = 3844
const randMaxShortNumber = 238327
const randShortNumber = randMaxShortNumber - randMinShortNumber

func UniqueId() string {
	buf := EncodeInt(uint64(GlobalRand1.Int63n(randNumber) + randMinNumber))
	buf = AppendInt(buf, uint64(time.Now().UnixNano()/1000))
	buf = AppendInt(buf, uint64(GlobalRand2.Int63n(randNumber)+randMinNumber))
	return string(cutId(buf, 20))
}

func ShortUniqueId() string {
	buf := EncodeInt(uint64(GlobalRand1.Int63n(randShortNumber) + randMinShortNumber))
	buf = AppendInt(buf, uint64(time.Now().UnixNano()/1000%(86400000000*100)))
	buf = AppendInt(buf, uint64(GlobalRand2.Int63n(randShortNumber)+randMinShortNumber))
	return string(buf)
}

func Id() string {
	buf := EncodeInt(uint64(GlobalRand1.Int63n(randNumber) + randMinNumber))
	buf = AppendInt(buf, uint64(GlobalRand2.Int63n(randNumber)+randMinNumber))
	return string(cutId(buf, 12))
}

func ShortId() string {
	buf := EncodeInt(uint64(GlobalRand1.Int63n(randShortNumber) + randMinShortNumber))
	buf = AppendInt(buf, uint64(GlobalRand2.Int63n(randShortNumber)+randMinShortNumber))
	return string(cutId(buf, 6))
}

func cutId(buf []byte, size int) []byte {
	if len(buf) > size {
		buf = buf[0:size]
	} else if len(buf) < size {
		for i := len(buf); i < size; i++ {
			buf = append(buf, '0')
		}
	}
	return buf
}
