package u

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
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


func EncryptAes(origData string, key []byte, iv []byte) string {
	key, iv = makeKeyIv(key, iv)
	block, err := aes.NewCipher(key)
	if err != nil {
		return ""
	}
	origDataBytes := []byte(origData)
	blockSize := block.BlockSize()
	origDataBytes = pkcs5Padding(origDataBytes, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv[:blockSize])
	crypted := make([]byte, len(origDataBytes))
	blockMode.CryptBlocks(crypted, origDataBytes)
	return base64.StdEncoding.EncodeToString(crypted)
}

func DecryptAes(crypted string, key []byte, iv []byte) string {
	key, iv = makeKeyIv(key, iv)
	cryptedBytes, err := base64.StdEncoding.DecodeString(crypted)
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, iv[:blockSize])
	origData := make([]byte, len(cryptedBytes))
	blockMode.CryptBlocks(origData, cryptedBytes)
	origData = pkcs5UnPadding(origData)
	return string(origData)
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func makeKeyIv(key []byte, iv []byte) ([]byte, []byte) {
	if len(key) >= 32 {
		key = key[0:32]
	} else if len(key) >= 16 {
		key = key[0:16]
	} else {
		for i := len(key); i < 16; i++ {
			key = append(key, 0)
		}
	}
	if len(iv) < len(key) {
		for i := len(iv); i < len(key); i++ {
			iv = append(iv, 0)
		}
	}
	return key, iv
}