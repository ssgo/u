package u

import (
	"time"
)

const randMin6BitNumber = 916132832                          // 99999u (100000)
const randMax6BitNumber = 56800235583                        // UUUUUU (999999)
const rand6BitNumber = randMax6BitNumber - randMin6BitNumber // 6位

const randMin5BitNumber = 14776336                           // 9999u (10000)
const randMax5BitNumber = 916132831                          // UUUUU (99999)
const rand5BitNumber = randMax5BitNumber - randMin5BitNumber // 5位

const randMin4BitNumber = 238328                             // 999u (1000)
const randMax4BitNumber = 14776335                           // UUUU (9999)
const rand4BitNumber = randMax4BitNumber - randMin4BitNumber // 4位

const randMin3BitNumber = 3844                               // 99u (100)
const randMax3BitNumber = 238327                             // UUU (999)
const rand3BitNumber = randMax3BitNumber - randMin3BitNumber // 3位

func UniqueId() string {
	buf := EncodeInt(uint64(GlobalRand1.Int63n(rand6BitNumber) + randMin6BitNumber))
	buf = AppendInt(buf, uint64(time.Now().UnixNano()/1000))
	buf = AppendInt(buf, uint64(GlobalRand2.Int63n(rand6BitNumber)+randMin6BitNumber))
	return string(cutId(buf, 20))
}

func ShortUniqueId() string {
	buf := EncodeInt(uint64(GlobalRand1.Int63n(rand3BitNumber) + randMin3BitNumber))
	buf = AppendInt(buf, uint64(time.Now().UnixNano()/1000%(86400000000*100)))
	buf = AppendInt(buf, uint64(GlobalRand2.Int63n(rand4BitNumber)+randMin4BitNumber))
	return string(cutId(buf, 14))
}

// 约312303亿亿
func Id12() string {
	buf := EncodeInt(uint64(GlobalRand1.Int63n(rand6BitNumber) + randMin6BitNumber))
	buf = AppendInt(buf, uint64(GlobalRand2.Int63n(rand6BitNumber)+randMin6BitNumber))
	return string(cutId(buf, 12))
}

// 约81亿亿
func Id10() string {
	buf := EncodeInt(uint64(GlobalRand1.Int63n(rand5BitNumber) + randMin5BitNumber))
	buf = AppendInt(buf, uint64(GlobalRand2.Int63n(rand5BitNumber)+randMin5BitNumber))
	return string(cutId(buf, 10))
}

// 约 2113536亿
func Id8() string {
	buf := EncodeInt(uint64(GlobalRand1.Int63n(rand4BitNumber) + randMin4BitNumber))
	buf = AppendInt(buf, uint64(GlobalRand2.Int63n(rand4BitNumber)+randMin4BitNumber))
	return string(cutId(buf, 8))
}

// 约 550亿
func Id6() string {
	buf := EncodeInt(uint64(GlobalRand1.Int63n(rand3BitNumber) + randMin3BitNumber))
	buf = AppendInt(buf, uint64(GlobalRand2.Int63n(rand3BitNumber)+randMin3BitNumber))
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
