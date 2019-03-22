package utility

import (
	"errors"
	"strings"
)

type intEncoder struct {
	Digits string
}

func (enc *intEncoder) EncodeInt(u uint64) []byte {
	return enc.AppendInt(nil, u)
}

func (enc *intEncoder) AppendInt(buf []byte, u uint64) []byte {
	if buf == nil {
		buf = make([]byte, 0)
	}
	for u >= 62 {
		q := u / 62
		buf = append(buf, enc.Digits[uint(u-q*62)])
		u = q
	}
	buf = append(buf, enc.Digits[uint(u)])
	return buf
}

func (enc *intEncoder) DecodeInt(buf []byte) uint64 {
	if buf == nil {
		return 0
	}
	var n uint64 = 0
	for i := len(buf) - 1; i >= 0; i-- {
		p := strings.IndexByte(enc.Digits, buf[i])
		if p >= 0 {
			n = n*62 + uint64(p)
		}
	}
	return n
}

func NewIntEncoder(digits string) (*intEncoder, error) {
	if len(digits) != 62 {
		return nil, errors.New("int encoder digits is bad " + digits)
	}

	m := map[int32]bool{}
	for _, d := range digits {
		if m[d] {
			return nil, errors.New("int encoder digits is repeated " + digits)
		}
		m[d] = true
	}

	e := intEncoder{}
	e.Digits = digits
	return &e, nil
}

var defaultIntEncoder, _ = NewIntEncoder("9ukH1grX75TQS6LzpFAjIivsdZoO0mc8NBwnyYDhtMWEC2V3KaGxfJRPqe4lbU")

func EncodeInt(u uint64) []byte {
	return defaultIntEncoder.AppendInt(nil, u)
}

func AppendInt(buf []byte, u uint64) []byte {
	return defaultIntEncoder.AppendInt(buf, u)
}

func DecodeInt(buf []byte) uint64 {
	return defaultIntEncoder.DecodeInt(buf)
}
