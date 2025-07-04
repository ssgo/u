package u

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"math/big"
	"strings"
)

type intEncoder struct {
	radix  uint8
	digits string
}

func (enc *intEncoder) EncodeInt(u uint64) []byte {
	return enc.AppendInt(nil, u)
}

func (enc *intEncoder) AppendInt(buf []byte, u uint64) []byte {
	if buf == nil {
		buf = make([]byte, 0)
	}
	radix := uint64(enc.radix)
	for u >= radix {
		q := u / radix
		buf = append(buf, enc.digits[uint(u-q*radix)])
		u = q
	}
	buf = append(buf, enc.digits[uint(u)])
	return buf
}

func (enc *intEncoder) FillInt(buf []byte, length int) []byte {
	radix := int(enc.radix)
	for i := len(buf); i < length; i++ {
		buf = enc.AppendInt(buf, uint64(GlobalRand1.Intn(radix)))
	}

	if len(buf) > length {
		buf = buf[0:length]
	}
	return buf
}

func (enc *intEncoder) DecodeInt(buf []byte) uint64 {
	radix := uint64(enc.radix)
	if buf == nil {
		return 0
	}
	var n uint64 = 0
	for i := len(buf) - 1; i >= 0; i-- {
		p := strings.IndexByte(enc.digits, buf[i])
		if p >= 0 {
			n = n*radix + uint64(p)
		}
	}
	return n
}

func NewIntEncoder(digits string, radix uint8) (*intEncoder, error) {
	if len(digits) < int(radix) {
		return nil, errors.New("int encoder digits is bad")
	}

	m := map[int32]bool{}
	for _, d := range digits {
		if m[d] {
			return nil, errors.New("int encoder digits is repeated " + digits)
		}
		m[d] = true
	}

	e := intEncoder{}
	e.digits = digits
	e.radix = radix
	return &e, nil
}

var defaultIntEncoder, _ = NewIntEncoder("9ukH1grX75TQS6LzpFAjIivsdZoO0mc8NBwnyYDhtMWEC2V3KaGxfJRPqe4lbU", 62)

func EncodeInt(u uint64) []byte {
	return defaultIntEncoder.AppendInt(nil, u)
}

func AppendInt(buf []byte, u uint64) []byte {
	return defaultIntEncoder.AppendInt(buf, u)
}

func FillInt(buf []byte, length int) []byte {
	return defaultIntEncoder.FillInt(buf, length)
}

func DecodeInt(buf []byte) uint64 {
	return defaultIntEncoder.DecodeInt(buf)
}

func ExchangeInt(buf []byte) []byte {
	size := len(buf)
	buf2 := make([]byte, size)
	buf2_i := 0
	buf2_ai := 0
	buf2_ri := size - 1
	for i := 0; i < size; i++ {
		if i%2 == 0 {
			// 从后往前取
			buf2[buf2_i] = buf[buf2_ri]
			buf2_i++
			buf2_ri--
		} else {
			// 从前往后取
			buf2[buf2_i] = buf[buf2_ai]
			buf2_i++
			buf2_ai++
		}
	}
	return buf2
}

func HashInt(buf []byte) []byte {
	prevP := 0
	for i, c := range []byte(buf) {
		p := strings.IndexByte(defaultIntEncoder.digits, c)
		if p < 0 {
			p = 0
		}
		if i > 0 {
			p = (prevP + p) % 62
			buf[i] = defaultIntEncoder.digits[p]
		}
		prevP = p
	}
	return buf
}

func EncryptAes(origData string, key []byte, iv []byte) (out string) {
	buf, _ := EncryptAesBytes([]byte(origData), key, iv)
	return base64.URLEncoding.EncodeToString(buf)
}

func EncryptAesBytes(origData []byte, key []byte, iv []byte) (out []byte, outErr error) {
	defer func() {
		if r := recover(); r != nil {
			out = nil
			outErr = errors.New("failed to encrypt aes")
		}
	}()

	key, iv = makeKeyIv(key, iv)
	if iv == nil {
		iv = key
	}
	if block, err := aes.NewCipher(key); err != nil {
		return origData, err
	} else {
		origDataBytes := origData
		blockSize := block.BlockSize()
		origDataBytes = pkcs5Padding(origDataBytes, blockSize)
		blockMode := cipher.NewCBCEncrypter(block, iv[:blockSize])
		crypted := make([]byte, len(origDataBytes))
		blockMode.CryptBlocks(crypted, origDataBytes)
		return crypted, nil
	}
}

func DecryptAes(crypted string, key []byte, iv []byte) (out string) {
	var base64Encoding *base64.Encoding
	if strings.ContainsRune(crypted, '_') || strings.ContainsRune(crypted, '-') {
		base64Encoding = base64.URLEncoding
	} else {
		base64Encoding = base64.StdEncoding
	}
	if cryptedBytes, err := base64Encoding.DecodeString(crypted); err != nil {
		return crypted
	} else {
		buf, err := DecryptAesBytes(cryptedBytes, key, iv)
		if err == nil && len(buf) > 0 {
			return String(buf)
		} else {
			return crypted
		}
	}
}

func DecryptAesBytes(cryptedBytes []byte, key []byte, iv []byte) (out []byte, outErr error) {
	defer func() {
		if r := recover(); r != nil {
			out = nil
			outErr = errors.New("failed to decrypt aes")
		}
	}()

	key, iv = makeKeyIv(key, iv)
	if iv == nil {
		iv = key
	}
	if block, err := aes.NewCipher(key); err != nil {
		return cryptedBytes, err
	} else {
		blockSize := block.BlockSize()
		blockMode := cipher.NewCBCDecrypter(block, iv[:blockSize])
		origData := make([]byte, len(cryptedBytes))
		blockMode.CryptBlocks(origData, cryptedBytes)
		origData = pkcs5UnPadding(origData)
		return origData, nil
	}
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	if length > 0 {
		unpadding := int(origData[length-1])
		pos := length - unpadding
		if pos < 0 || pos >= length {
			return nil
		}
		return origData[:pos]
	}
	return origData
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
	if iv != nil {
		if len(iv) < len(key) {
			for i := len(iv); i < len(key); i++ {
				iv = append(iv, 0)
			}
		}
	}
	return key, iv
}

func MD5Base64(data string) string {
	return base64.StdEncoding.EncodeToString(MD5([]byte(data)))
}

func MD5String(data string) string {
	return hex.EncodeToString(MD5([]byte(data)))
}

func MD5(data []byte) []byte {
	hash := md5.New()
	hash.Write(data)
	return hash.Sum([]byte{})
}

func Sha1Base64(data string) string {
	return base64.StdEncoding.EncodeToString(Sha1([]byte(data)))
}

func Sha1String(data string) string {
	return hex.EncodeToString(Sha1([]byte(data)))
}

func Sha1(data []byte) []byte {
	hash := sha1.New()
	hash.Write(data)
	return hash.Sum([]byte{})
}

func Sha256Base64(data string) string {
	return base64.StdEncoding.EncodeToString(Sha256([]byte(data)))
}

func Sha256String(data string) string {
	return hex.EncodeToString(Sha256([]byte(data)))
}

func Sha256(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum([]byte{})
}

func Sha512Base64(data string) string {
	return base64.StdEncoding.EncodeToString(Sha512([]byte(data)))
}

func Sha512String(data string) string {
	return hex.EncodeToString(Sha512([]byte(data)))
}

func Sha512(data []byte) []byte {
	hash := sha512.New()
	hash.Write(data)
	return hash.Sum([]byte{})
}

func Base64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func UrlBase64(data []byte) string {
	return base64.URLEncoding.EncodeToString(data)
}

func Base64String(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func UrlBase64String(data string) string {
	return base64.URLEncoding.EncodeToString([]byte(data))
}

func UnBase64(data string) []byte {
	buf, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return []byte{}
	}
	return buf
}

func UnUrlBase64(data string) []byte {
	buf, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		return []byte{}
	}
	return buf
}

func UnBase64String(data string) string {
	return string(UnBase64(data))
}

func UnUrlBase64String(data string) string {
	return string(UnUrlBase64(data))
}

type Aes struct {
	key []byte
	iv  []byte
}

func NewAes(key, iv []byte) *Aes {
	return &Aes{key: key, iv: iv}
}

func (_this *Aes) EncryptBytes(data []byte) (out []byte, outErr error) {
	return EncryptAesBytes(data, _this.key, _this.iv)
	//defer func() {
	//	if r := recover(); r != nil {
	//		out = data
	//	}
	//}()
	//
	//key, iv := makeKeyIv(_this.key, _this.iv)
	//if iv == nil {
	//	iv = key
	//}
	//
	//block, err := aes.NewCipher(key)
	//if err != nil {
	//	return data
	//}
	//origDataBytes := data
	//blockSize := block.BlockSize()
	//origDataBytes = pkcs5Padding(origDataBytes, blockSize)
	//blockMode := cipher.NewCBCEncrypter(block, iv[:blockSize])
	//crypted := make([]byte, len(origDataBytes))
	//blockMode.CryptBlocks(crypted, origDataBytes)
	//return crypted
}

func (_this *Aes) DecryptBytes(data []byte) (out []byte, outErr error) {
	return DecryptAesBytes(data, _this.key, _this.iv)
	//defer func() {
	//	if r := recover(); r != nil {
	//		out = data
	//	}
	//}()
	//
	//key, iv := makeKeyIv(_this.key, _this.iv)
	//if iv == nil {
	//	iv = key
	//}
	//block, err := aes.NewCipher(key)
	//if err != nil {
	//	fmt.Println(err)
	//	return nil
	//}
	//blockSize := block.BlockSize()
	//blockMode := cipher.NewCBCDecrypter(block, iv[:blockSize])
	//origData := make([]byte, len(data))
	//blockMode.CryptBlocks(origData, data)
	//origData = pkcs5UnPadding(origData)
	//return origData
}

func (_this *Aes) EncryptBytesToHex(data []byte) string {
	if buf, err := _this.EncryptBytes(data); err == nil && len(buf) > 0 {
		return hex.EncodeToString(buf)
	} else {
		return string(data)
	}
}

func (_this *Aes) EncryptBytesToBase64(data []byte) string {
	if buf, err := _this.EncryptBytes(data); err == nil && len(buf) > 0 {
		return base64.StdEncoding.EncodeToString(buf)
	} else {
		return string(data)
	}
}

func (_this *Aes) EncryptBytesToUrlBase64(data []byte) string {
	if buf, err := _this.EncryptBytes(data); err == nil && len(buf) > 0 {
		return base64.URLEncoding.EncodeToString(buf)
	} else {
		return string(data)
	}
}

func (_this *Aes) EncryptStringToHex(data string) string {
	if buf, err := _this.EncryptBytes([]byte(data)); err == nil && len(buf) > 0 {
		return hex.EncodeToString(buf)
	} else {
		return data
	}
}

func (_this *Aes) EncryptStringToBase64(data string) string {
	if buf, err := _this.EncryptBytes([]byte(data)); err == nil && len(buf) > 0 {
		return base64.StdEncoding.EncodeToString(buf)
	} else {
		return data
	}
}

func (_this *Aes) EncryptStringToUrlBase64(data string) string {
	if buf, err := _this.EncryptBytes([]byte(data)); err == nil && len(buf) > 0 {
		return base64.URLEncoding.EncodeToString(buf)
	} else {
		return data
	}
}

func (_this *Aes) DecryptHexToBytes(data string) []byte {
	buf, err := hex.DecodeString(data)
	if err != nil {
		return []byte(data)
	}
	if buf1, err1 := _this.DecryptBytes(buf); err1 == nil && len(buf1) > 0 {
		return buf1
	} else {
		return []byte(data)
	}
}

func (_this *Aes) DecryptBase64ToBytes(data string) []byte {
	buf, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return []byte(data)
	}
	if buf1, err1 := _this.DecryptBytes(buf); err1 == nil && len(buf1) > 0 {
		return buf1
	} else {
		return []byte(data)
	}
}

func (_this *Aes) DecryptUrlBase64ToBytes(data string) []byte {
	buf, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		return []byte(data)
	}
	if buf1, err1 := _this.DecryptBytes(buf); err1 == nil && len(buf1) > 0 {
		return buf1
	} else {
		return []byte(data)
	}
}

func (_this *Aes) DecryptHexToString(data string) string {
	buf, err := hex.DecodeString(data)
	if err != nil {
		return data
	}
	if buf1, err1 := _this.DecryptBytes(buf); err1 == nil && len(buf1) > 0 {
		return string(buf1)
	} else {
		return data
	}
}

func (_this *Aes) DecryptBase64ToString(data string) string {
	buf, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return data
	}
	if buf1, err1 := _this.DecryptBytes(buf); err1 == nil && len(buf1) > 0 {
		return string(buf1)
	} else {
		return data
	}
}

func (_this *Aes) DecryptUrlBase64ToString(data string) string {
	buf, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		return data
	}
	if buf1, err1 := _this.DecryptBytes(buf); err1 == nil && len(buf1) > 0 {
		return string(buf1)
	} else {
		return data
	}
}

func GenECDSA521Key() (privateKey string, publicKey string, err error) {
	return GenECDSAKey(elliptic.P521())
}

func GenECDSA256Key() (privateKey string, publicKey string, err error) {
	return GenECDSAKey(elliptic.P256())
}

func GenECDSA384Key() (privateKey string, publicKey string, err error) {
	return GenECDSAKey(elliptic.P384())
}

func GenECDSAKey(curve elliptic.Curve) (privateKey string, publicKey string, err error) {
	priKey, err := ecdsa.GenerateKey(curve, GlobalRand2)
	if err != nil {
		return "", "", err
	}
	//ecPrivateKey, err := x509.MarshalECPrivateKey(priKey)
	//if err != nil {
	//	return "", "", err
	//}
	privateKey = base64.URLEncoding.EncodeToString(priKey.D.Bytes())
	var buf bytes.Buffer
	buf.Write(priKey.X.Bytes())
	buf.Write(priKey.Y.Bytes())
	publicKey = base64.URLEncoding.EncodeToString(buf.Bytes())
	return
}

func MakeECDSA256PrivateKey(privateKeyStr string) (priKey *ecdsa.PrivateKey, err error) {
	return MakeECDSAPrivateKey(privateKeyStr, elliptic.P256())
}

func MakeECDSA384PrivateKey(privateKeyStr string) (priKey *ecdsa.PrivateKey, err error) {
	return MakeECDSAPrivateKey(privateKeyStr, elliptic.P384())
}

func MakeECDSA521PrivateKey(privateKeyStr string) (priKey *ecdsa.PrivateKey, err error) {
	return MakeECDSAPrivateKey(privateKeyStr, elliptic.P521())
}

func MakeECDSAPrivateKey(privateKeyStr string, curve elliptic.Curve) (priKey *ecdsa.PrivateKey, err error) {
	bytes, err := base64.URLEncoding.DecodeString(privateKeyStr)
	if err != nil {
		return nil, err
	}
	//priKey, err = x509.ParseECPrivateKey(bytes)
	//if err != nil {
	//	return nil, err
	//}
	x, y := curve.ScalarBaseMult(bytes)
	return &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		},
		D: new(big.Int).SetBytes(bytes),
	}, nil
}

func MakeECDSA256PublicKey(publicKeyStr string) (pubKey *ecdsa.PublicKey, err error) {
	return MakeECDSAPublicKey(publicKeyStr, elliptic.P256())
}

func MakeECDSA384PublicKey(publicKeyStr string) (pubKey *ecdsa.PublicKey, err error) {
	return MakeECDSAPublicKey(publicKeyStr, elliptic.P384())
}

func MakeECDSA521PublicKey(publicKeyStr string) (pubKey *ecdsa.PublicKey, err error) {
	return MakeECDSAPublicKey(publicKeyStr, elliptic.P521())
}

func MakeECDSAPublicKey(publicKeyStr string, curve elliptic.Curve) (pubKey *ecdsa.PublicKey, err error) {
	bytes, err := base64.URLEncoding.DecodeString(publicKeyStr)
	if err != nil {
		return nil, err
	}
	x := new(big.Int)
	y := new(big.Int)
	byteLen := len(bytes) / 2
	x.SetBytes(bytes[0:byteLen])
	y.SetBytes(bytes[byteLen:])
	pub := ecdsa.PublicKey{Curve: curve, X: x, Y: y}
	pubKey = &pub
	return
}

func SignECDSA(content []byte, priKey *ecdsa.PrivateKey) (signature string, err error) {
	r, s, err := ecdsa.Sign(GlobalRand1, priKey, Sha256(content))
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	buf.WriteByte(byte(len(r.Bytes())))
	buf.Write(r.Bytes())
	buf.Write(s.Bytes())
	signature = base64.URLEncoding.EncodeToString(buf.Bytes())
	return
}

func VerifyECDSA(content []byte, signature string, pubKey *ecdsa.PublicKey) bool {
	bytes, e := base64.URLEncoding.DecodeString(signature)
	if e != nil {
		return false
	}
	r := new(big.Int)
	s := new(big.Int)
	byteLen := bytes[0]
	r.SetBytes(bytes[1 : byteLen+1])
	s.SetBytes(bytes[byteLen+1:])
	//fmt.Println(len(bytes[0:byteLen]), len(bytes[byteLen:]), "/", len(bytes))
	return ecdsa.Verify(pubKey, Sha256(content), r, s)
}

func MakeToken(size int) []byte {
	token := make([]byte, size)
	for i := 0; i < size; i++ {
		var r int
		if i%2 == 1 {
			r = GlobalRand1.Intn(255)
		} else {
			r = GlobalRand2.Intn(255)
		}
		token[i] = byte(r)
	}
	return token
}
