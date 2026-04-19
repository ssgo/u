package u

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"runtime"

	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/hkdf"
)

type IntEncoder struct {
	radix     uint8
	digits    string
	decodeMap [256]int
}

func (enc *IntEncoder) EncodeInt(u uint64) []byte {
	return enc.AppendInt(nil, u)
}

func (enc *IntEncoder) AppendInt(buf []byte, u uint64) []byte {
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

func (enc *IntEncoder) FillInt(buf []byte, length int) []byte {
	currLen := len(buf)
	if currLen >= length {
		return buf
	}
	if cap(buf) < length {
		newBuf := make([]byte, currLen, length)
		copy(newBuf, buf)
		buf = newBuf
	}
	buf = buf[:length]
	radix := uint(enc.radix)
	for i := currLen; i < length; i++ {
		idx := GlobalRand2.UintN(radix)
		buf[i] = enc.digits[idx]
	}
	return buf
}

func (enc *IntEncoder) DecodeInt(buf []byte) uint64 {
	radix := uint64(enc.radix)
	if buf == nil {
		return 0
	}
	var n uint64 = 0
	for i := len(buf) - 1; i >= 0; i-- {
		p := enc.decodeMap[buf[i]]
		if p >= 0 {
			n = n*radix + uint64(p)
		}
	}
	return n
}

func NewIntEncoder(digits string, radix uint8) (*IntEncoder, error) {
	if len(digits) < int(radix) {
		return nil, errors.New("int encoder digits is bad")
	}

	e := IntEncoder{digits: digits, radix: radix, decodeMap: [256]int{}}
	for i := 0; i < 256; i++ {
		e.decodeMap[i] = -1
	}
	m := map[int32]bool{}
	for i, d := range digits {
		e.decodeMap[digits[i]] = i
		if m[d] {
			return nil, errors.New("int encoder digits is repeated " + digits)
		}
		m[d] = true
	}

	return &e, nil
}

var DefaultIntEncoder, _ = NewIntEncoder("9ukH1grX75TQS6LzpFAjIivsdZoO0mc8NBwnyYDhtMWEC2V3KaGxfJRPqe4lbU", 62)
var OrderedIntEncoder, _ = NewIntEncoder("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", 62)

func EncodeInt(u uint64) []byte {
	return DefaultIntEncoder.AppendInt(nil, u)
}

func AppendInt(buf []byte, u uint64) []byte {
	return DefaultIntEncoder.AppendInt(buf, u)
}

func FillInt(buf []byte, length int) []byte {
	return DefaultIntEncoder.FillInt(buf, length)
}

func DecodeInt(buf []byte) uint64 {
	return DefaultIntEncoder.DecodeInt(buf)
}

func ExchangeInt(buf []byte) []byte {
	size := len(buf)
	buf2 := make([]byte, size)
	buf2_i := 0
	buf2_ai := 0
	buf2_ri := size - 1
	for i := range size {
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
	return DefaultIntEncoder.HashInt(buf)
}

func (enc *IntEncoder) HashInt(buf []byte) []byte {
	if len(buf) == 0 {
		return buf
	}
	prevP := (len(buf) * 17) % int(enc.radix)
	for i, c := range buf {
		p := enc.decodeMap[c]
		if p < 0 {
			p = 0
		}
		p = (prevP + p) % int(enc.radix)
		buf[i] = enc.digits[p]
		prevP = p
	}
	return buf
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

func MD5Base64(data string) string {
	return base64.StdEncoding.EncodeToString(MD5([]byte(data)))
}

func MD5String(data string) string {
	return hex.EncodeToString(MD5([]byte(data)))
}

func MD5(data ...[]byte) []byte {
	hash := md5.New()
	for _, v := range data {
		hash.Write(v)
	}
	return hash.Sum([]byte{})
}

func Sha1Base64(data string) string {
	return base64.StdEncoding.EncodeToString(Sha1([]byte(data)))
}

func Sha1String(data string) string {
	return hex.EncodeToString(Sha1([]byte(data)))
}

func Sha1(data ...[]byte) []byte {
	hash := sha1.New()
	for _, v := range data {
		hash.Write(v)
	}
	return hash.Sum([]byte{})
}

func Sha256Base64(data string) string {
	return base64.StdEncoding.EncodeToString(Sha256([]byte(data)))
}

func Sha256String(data string) string {
	return hex.EncodeToString(Sha256([]byte(data)))
}

func Sha256(data ...[]byte) []byte {
	hash := sha256.New()
	for _, v := range data {
		hash.Write(v)
	}
	return hash.Sum([]byte{})
}

func Sha512Base64(data string) string {
	return base64.StdEncoding.EncodeToString(Sha512([]byte(data)))
}

func Sha512String(data string) string {
	return hex.EncodeToString(Sha512([]byte(data)))
}

func Sha512(data ...[]byte) []byte {
	hash := sha512.New()
	for _, v := range data {
		hash.Write(v)
	}
	return hash.Sum([]byte{})
}

func HmacMD5(key []byte, data ...[]byte) []byte {
	hash := hmac.New(md5.New, key)
	for _, v := range data {
		hash.Write(v)
	}
	return hash.Sum([]byte{})
}

func HmacSha1(key []byte, data ...[]byte) []byte {
	hash := hmac.New(sha1.New, key)
	for _, v := range data {
		hash.Write(v)
	}
	return hash.Sum([]byte{})
}

func HmacSha256(key []byte, data ...[]byte) []byte {
	hash := hmac.New(sha256.New, key)
	for _, v := range data {
		hash.Write(v)
	}
	return hash.Sum([]byte{})
}

func HmacSha512(key []byte, data ...[]byte) []byte {
	hash := hmac.New(sha512.New, key)
	for _, v := range data {
		hash.Write(v)
	}
	return hash.Sum([]byte{})
}

func Hex(data []byte) []byte {
	dst := make([]byte, hex.EncodedLen(len(data)))
	hex.Encode(dst, data)
	return dst
}

func HexString(data []byte) string {
	return hex.EncodeToString(data)
}

func UnHex(data []byte) ([]byte, error) {
	return hex.DecodeString(string(data))
}

func UnHexString(data string) ([]byte, error) {
	return hex.DecodeString(data)
}

func Base64(data []byte) []byte {
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(buf, data)
	return buf
}

func UrlBase64(data []byte) []byte {
	buf := make([]byte, base64.URLEncoding.EncodedLen(len(data)))
	base64.URLEncoding.Encode(buf, data)
	return buf
}

func Base64String(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func UrlBase64String(data []byte) string {
	return base64.URLEncoding.EncodeToString(data)
}

func UnBase64(data []byte) ([]byte, error) {
	dbuf := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	n, err := base64.StdEncoding.Decode(dbuf, data)
	return dbuf[:n], err
}

func UnUrlBase64(data []byte) ([]byte, error) {
	dbuf := make([]byte, base64.URLEncoding.DecodedLen(len(data)))
	n, err := base64.URLEncoding.Decode(dbuf, data)
	return dbuf[:n], err
}

func UnBase64String(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}

func UnUrlBase64String(data string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(data)
}

func MakeToken(size int) []byte {
	token := make([]byte, size)
	rand.Read(token)
	return token
}

// v3.3
// ------ Symmetric ---- //

type SymmetricCipher interface {
	Encrypt(data []byte, key []byte, iv []byte) ([]byte, error)
	Decrypt(data []byte, key []byte, iv []byte) ([]byte, error)
}

type Symmetric struct {
	cipher SymmetricCipher
	key    *SafeBuf
	iv     *SafeBuf
}

func NewSymmetric(cipher SymmetricCipher, safeKeyBuf, safeIvBuf *SafeBuf) (*Symmetric, error) {
	keyBuf := safeKeyBuf.Open()
	defer keyBuf.Close()
	ivBuf := safeIvBuf.Open()
	defer ivBuf.Close()
	return NewSymmetricWithOutEraseKey(cipher, keyBuf.Data, ivBuf.Data)
}

func NewSymmetricAndEraseKey(cipher SymmetricCipher, key, iv []byte) (*Symmetric, error) {
	defer ZeroMemory(key)
	defer ZeroMemory(iv)
	return NewSymmetricWithOutEraseKey(cipher, key, iv)
}

func NewSymmetricWithOutEraseKey(cipher SymmetricCipher, key, iv []byte) (*Symmetric, error) {
	keySize := 16
	if len(key) >= 32 {
		keySize = 32
	} else if len(key) >= 24 {
		keySize = 24
	} else if len(key) < 16 {
		return nil, errors.New("key or iv size is not 16 24 or 32")
	}
	symmetric := &Symmetric{cipher: cipher, key: NewSafeBuf(key[:keySize]), iv: NewSafeBuf(iv)}
	runtime.SetFinalizer(symmetric, func(obj *Symmetric) {
		obj.Close()
	})
	return symmetric, nil
}

func (_this *Symmetric) Close() {
	_this.key.Close()
	_this.iv.Close()
}

func (_this *Symmetric) Encrypt(safeBuf *SafeBuf) ([]byte, error) {
	buf := safeBuf.Open()
	defer buf.Close()
	return _this.EncryptBytes(buf.Data)
}

func (_this *Symmetric) EncryptBytes(data []byte) ([]byte, error) {
	key := _this.key.Open()
	defer key.Close()
	iv := _this.iv.Open()
	defer iv.Close()
	return _this.cipher.Encrypt(data, key.Data, iv.Data)
}

func (_this *Symmetric) Decrypt(data []byte) (*SafeBuf, error) {
	buf, err := _this.DecryptBytes(data)
	if err != nil {
		return nil, err
	}
	defer ZeroMemory(buf)
	return NewSafeBuf(buf), nil
}

func (_this *Symmetric) DecryptBytes(data []byte) ([]byte, error) {
	key := _this.key.Open()
	defer key.Close()
	iv := _this.iv.Open()
	defer iv.Close()
	return _this.cipher.Decrypt(data, key.Data, iv.Data)
}

func (_this *Symmetric) DecryptBytesN(data []byte) []byte {
	r, err := _this.DecryptBytes(data)
	if err != nil {
		return data
	}
	return r
}

// ------ AES ------

type AESCipher struct{ useGCM bool }

var AESCBC = &AESCipher{useGCM: false}
var AESGCM = &AESCipher{useGCM: true}

func NewAESCBC(safeKeyBuf, safeIvBuf *SafeBuf) (*Symmetric, error) {
	return NewSymmetric(AESCBC, safeKeyBuf, safeIvBuf)
}
func NewAESCBCAndEraseKey(safeKeyBuf, safeIvBuf []byte) (*Symmetric, error) {
	return NewSymmetricAndEraseKey(AESCBC, safeKeyBuf, safeIvBuf)
}
func NewAESCBCWithOutEraseKey(safeKeyBuf, safeIvBuf []byte) (*Symmetric, error) {
	return NewSymmetricWithOutEraseKey(AESCBC, safeKeyBuf, safeIvBuf)
}

func NewAESGCM(safeKeyBuf, safeIvBuf *SafeBuf) (*Symmetric, error) {
	return NewSymmetric(AESGCM, safeKeyBuf, safeIvBuf)
}
func NewAESGCMAndEraseKey(safeKeyBuf, safeIvBuf []byte) (*Symmetric, error) {
	return NewSymmetricAndEraseKey(AESGCM, safeKeyBuf, safeIvBuf)
}
func NewAESGCMWithOutEraseKey(safeKeyBuf, safeIvBuf []byte) (*Symmetric, error) {
	return NewSymmetricWithOutEraseKey(AESGCM, safeKeyBuf, safeIvBuf)
}

func (_this *AESCipher) Encrypt(data []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	if _this.useGCM {
		aesgcm, err := cipher.NewGCM(block)
		if err != nil {
			return nil, err
		}
		crypted := aesgcm.Seal(nil, iv[:aesgcm.NonceSize()], data, nil)
		return crypted, nil
	} else {
		origDataBytes := pkcs5Padding(data, blockSize)
		blockMode := cipher.NewCBCEncrypter(block, iv[:blockSize])
		crypted := make([]byte, len(origDataBytes))
		blockMode.CryptBlocks(crypted, origDataBytes)
		return crypted, nil
	}
}

func (_this *AESCipher) Decrypt(data []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	if _this.useGCM {
		aesgcm, err := cipher.NewGCM(block)
		if err != nil {
			return nil, err
		}
		return aesgcm.Open(nil, iv[:aesgcm.NonceSize()], data, nil)
	} else {
		blockMode := cipher.NewCBCDecrypter(block, iv[:blockSize])
		origData := make([]byte, len(data))
		blockMode.CryptBlocks(origData, data)
		origData = pkcs5UnPadding(origData)
		return origData, nil
	}
}

// ------ Asymmetric ------

type AsymmetricAlgorithm interface {
	ParsePrivateKey(der []byte) (any, error)
	ParsePublicKey(der []byte) (any, error)
	Sign(privateKey any, data []byte, hash ...crypto.Hash) ([]byte, error)
	Verify(publicKey any, data []byte, signature []byte, hash ...crypto.Hash) (bool, error)
}

type AsymmetricCipherAlgorithm interface {
	Encrypt(publicKey any, data []byte) ([]byte, error)
	Decrypt(privateKey any, data []byte) ([]byte, error)
}

type Asymmetric struct {
	algorithm     AsymmetricAlgorithm
	privateKeyBuf *SafeBuf
	publicKeyBuf  *SafeBuf
	privateKey    any
	publicKey     any
}

func NewAsymmetric(algorithm AsymmetricAlgorithm, safePrivateKeyBuf, safePublicKeyBuf *SafeBuf) (*Asymmetric, error) {
	var privKey, pubKey []byte
	if safePrivateKeyBuf != nil {
		privateKeyBuf := safePrivateKeyBuf.Open()
		defer privateKeyBuf.Close()
		privKey = privateKeyBuf.Data
	}
	if safePublicKeyBuf != nil {
		publicKeyBuf := safePublicKeyBuf.Open()
		defer publicKeyBuf.Close()
		pubKey = publicKeyBuf.Data
	}
	return NewAsymmetricWithoutEraseKey(algorithm, privKey, pubKey, false)
}

func NewAsymmetricAndEraseKey(algorithm AsymmetricAlgorithm, privateKey, publicKey []byte) (*Asymmetric, error) {
	if privateKey != nil {
		defer ZeroMemory(privateKey)
	}
	if publicKey != nil {
		defer ZeroMemory(publicKey)
	}
	return NewAsymmetricWithoutEraseKey(algorithm, privateKey, publicKey, false)
}

func NewAsymmetricWithoutEraseKey(algorithm AsymmetricAlgorithm, privateKey, publicKey []byte, fastModeButIsNotSecure bool) (*Asymmetric, error) {
	_this := &Asymmetric{algorithm: algorithm}
	var err error
	if privateKey != nil {
		if fastModeButIsNotSecure {
			if _this.privateKey, err = algorithm.ParsePrivateKey(privateKey); err != nil {
				return nil, err
			}
		} else {
			_this.privateKeyBuf = NewSafeBuf(privateKey)
		}
	}
	if publicKey != nil {
		if fastModeButIsNotSecure {
			if _this.publicKey, err = algorithm.ParsePublicKey(publicKey); err != nil {
				return nil, err
			}
		} else {
			_this.publicKeyBuf = NewSafeBuf(publicKey)
		}
	}
	runtime.SetFinalizer(_this, func(obj *Asymmetric) {
		obj.Close()
	})
	return _this, nil
}

func (_this *Asymmetric) Close() {
	if _this.privateKeyBuf != nil {
		_this.privateKeyBuf.Close()
		_this.privateKeyBuf = nil
	}
	if _this.publicKeyBuf != nil {
		_this.publicKeyBuf.Close()
		_this.publicKeyBuf = nil
	}
}

func (_this *Asymmetric) Sign(data []byte) ([]byte, error) {
	privateKey := _this.privateKey
	if privateKey == nil && _this.privateKeyBuf != nil {
		privBuf := _this.privateKeyBuf.Open()
		defer privBuf.Close()
		var err error
		if privateKey, err = _this.algorithm.ParsePrivateKey(privBuf.Data); err != nil {
			return nil, err
		}
	}
	if privateKey == nil {
		return nil, errors.New("private key is not set")
	}
	return _this.algorithm.Sign(privateKey, data)
}

func (_this *Asymmetric) Verify(data []byte, signature []byte) (bool, error) {
	publicKey := _this.publicKey
	if publicKey == nil && _this.publicKeyBuf != nil {
		pubBuf := _this.publicKeyBuf.Open()
		defer pubBuf.Close()
		var err error
		if publicKey, err = _this.algorithm.ParsePublicKey(pubBuf.Data); err != nil {
			return false, err
		}
	}
	if publicKey == nil {
		return false, errors.New("public key is not set")
	}
	return _this.algorithm.Verify(publicKey, data, signature)
}

func (_this *Asymmetric) Encrypt(data []byte) ([]byte, error) {
	cipherAlgo, ok := _this.algorithm.(AsymmetricCipherAlgorithm)
	if !ok {
		return nil, errors.New("the current algorithm does not support encryption")
	}

	publicKey := _this.publicKey
	if publicKey == nil && _this.publicKeyBuf != nil {
		pubBuf := _this.publicKeyBuf.Open()
		defer pubBuf.Close()
		var err error
		if publicKey, err = _this.algorithm.ParsePublicKey(pubBuf.Data); err != nil {
			return nil, err
		}
	}
	if publicKey == nil {
		return nil, errors.New("public key is not set")
	}

	return cipherAlgo.Encrypt(publicKey, data)
}

func (_this *Asymmetric) Decrypt(data []byte) ([]byte, error) {
	cipherAlgo, ok := _this.algorithm.(AsymmetricCipherAlgorithm)
	if !ok {
		return nil, errors.New("the current algorithm does not support decryption")
	}

	privateKey := _this.privateKey
	if privateKey == nil && _this.privateKeyBuf != nil {
		privBuf := _this.privateKeyBuf.Open()
		defer privBuf.Close()
		var err error
		if privateKey, err = _this.algorithm.ParsePrivateKey(privBuf.Data); err != nil {
			return nil, err
		}
	}
	if privateKey == nil {
		return nil, errors.New("private key is not set")
	}

	return cipherAlgo.Decrypt(privateKey, data)
}

// ------ ECDSA ------ //

type ECDSAAlgorithm struct {
	useGCM  bool
	kdfInfo []byte
	kdfSalt []byte
	hash    crypto.Hash
}

var ECDSAGCM = &ECDSAAlgorithm{useGCM: true, hash: crypto.SHA256}
var ECDSACBC = &ECDSAAlgorithm{useGCM: false, hash: crypto.SHA256}

func NewECDSA(safePrivateKeyBuf, safePublicKeyBuf *SafeBuf) (*Asymmetric, error) {
	return NewAsymmetric(ECDSAGCM, safePrivateKeyBuf, safePublicKeyBuf)
}
func NewECDSAndEraseKey(safePrivateKeyBuf, safePublicKeyBuf []byte) (*Asymmetric, error) {
	return NewAsymmetricAndEraseKey(ECDSAGCM, safePrivateKeyBuf, safePublicKeyBuf)
}
func NewECDSAWithOutEraseKey(safePrivateKeyBuf, safePublicKeyBuf []byte) (*Asymmetric, error) {
	return NewAsymmetricWithoutEraseKey(ECDSAGCM, safePrivateKeyBuf, safePublicKeyBuf, false)
}

func NewECDSAAlgorithm(useGCM bool, hash crypto.Hash, kdfInfo, kdfSalt []byte) *ECDSAAlgorithm {
	return &ECDSAAlgorithm{useGCM: useGCM, kdfInfo: kdfInfo, kdfSalt: kdfSalt, hash: hash}
}

func GenerateECDSAKeyPair(bitSize int) ([]byte, []byte, error) {
	var curve elliptic.Curve
	switch bitSize {
	case 256:
		curve = elliptic.P256()
	case 384:
		curve = elliptic.P384()
	default:
		curve = elliptic.P521()
	}

	priKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	privateKey, err := x509.MarshalECPrivateKey(priKey)
	if err != nil {
		return nil, nil, err
	}
	publicKey, err := x509.MarshalPKIXPublicKey(&priKey.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, publicKey, nil
}

func (e *ECDSAAlgorithm) ParsePrivateKey(der []byte) (any, error) {
	return x509.ParseECPrivateKey(der)
}

func (e *ECDSAAlgorithm) ParsePublicKey(der []byte) (any, error) {
	pubKeyAny, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return nil, err
	}
	pubKey, ok := pubKeyAny.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("not an ECDSA public key")
	}
	return pubKey, nil
}

func (e *ECDSAAlgorithm) Sign(privateKeyObj any, data []byte, hash ...crypto.Hash) ([]byte, error) {
	privKey, ok := privateKeyObj.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("invalid private key")
	}
	if len(hash) == 0 {
		hash = append(hash, crypto.SHA256)
	}
	hasher := hash[0].New()
	hasher.Write(data)
	hashed := hasher.Sum(nil)
	return ecdsa.SignASN1(rand.Reader, privKey, hashed)
}

func (e *ECDSAAlgorithm) Verify(publicKeyObj any, data []byte, signature []byte, hash ...crypto.Hash) (bool, error) {
	pubKey, ok := publicKeyObj.(*ecdsa.PublicKey)
	if !ok {
		return false, errors.New("invalid public key")
	}
	if len(hash) == 0 {
		hash = append(hash, crypto.SHA256)
	}
	hasher := hash[0].New()
	hasher.Write(data)
	hashed := hasher.Sum(nil)
	return ecdsa.VerifyASN1(pubKey, hashed, signature), nil
}

func (e *ECDSAAlgorithm) Encrypt(publicKeyObj any, data []byte) ([]byte, error) {
	ecdsaPub, ok := publicKeyObj.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("invalid public key type for ECDSA")
	}
	ecdhPub, err := ecdsaPub.ECDH()
	if err != nil {
		return nil, err
	}

	ephemeralPriv, err := ecdhPub.Curve().GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	sharedSecret, err := ephemeralPriv.ECDH(ecdhPub)
	if err != nil {
		return nil, err
	}
	defer ZeroMemory(sharedSecret)

	hkdfReader := hkdf.New(e.hash.New, sharedSecret, e.kdfSalt, e.kdfInfo)
	aesKey := make([]byte, 32)
	if _, err := io.ReadFull(hkdfReader, aesKey); err != nil {
		return nil, err
	}
	defer ZeroMemory(aesKey)

	cipherAlgo := &AESCipher{useGCM: e.useGCM}
	ivLen := 16
	if e.useGCM {
		ivLen = 12
	}
	iv := MakeSafeToken(ivLen)

	cipherText, err := cipherAlgo.Encrypt(data, aesKey, iv)
	if err != nil {
		return nil, err
	}

	pubBytes := ephemeralPriv.PublicKey().Bytes()
	out := make([]byte, 0, len(pubBytes)+len(iv)+len(cipherText))
	out = append(out, pubBytes...)
	out = append(out, iv...)
	out = append(out, cipherText...)
	return out, nil
}

func (e *ECDSAAlgorithm) Decrypt(privateKeyObj any, data []byte) ([]byte, error) {
	ecdsaPriv, ok := privateKeyObj.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("invalid private key type for ECDSA")
	}

	ecdhPriv, err := ecdsaPriv.ECDH()
	if err != nil {
		return nil, err
	}

	pubKeyLen := len(ecdhPriv.PublicKey().Bytes())
	ivLen := 16
	if e.useGCM {
		ivLen = 12
	}

	if len(data) < pubKeyLen+ivLen {
		return nil, errors.New("invalid ciphertext package size")
	}

	ephemeralPubBytes := data[:pubKeyLen]
	iv := data[pubKeyLen : pubKeyLen+ivLen]
	cipherText := data[pubKeyLen+ivLen:]

	ephemeralPub, err := ecdhPriv.Curve().NewPublicKey(ephemeralPubBytes)
	if err != nil {
		return nil, err
	}

	sharedSecret, err := ecdhPriv.ECDH(ephemeralPub)
	if err != nil {
		return nil, err
	}
	defer ZeroMemory(sharedSecret)

	hkdfReader := hkdf.New(e.hash.New, sharedSecret, e.kdfSalt, e.kdfInfo)
	aesKey := make([]byte, 32)
	if _, err := io.ReadFull(hkdfReader, aesKey); err != nil {
		return nil, err
	}
	defer ZeroMemory(aesKey)

	cipherAlgo := &AESCipher{useGCM: e.useGCM}
	plainText, err := cipherAlgo.Decrypt(cipherText, aesKey, iv)
	return plainText, err
}

// ------ RSA ------ //

type RSAAlgorithm struct {
	isPSS  bool
	isOAEP bool
}

var RSA = &RSAAlgorithm{isPSS: true, isOAEP: true}

// Deprecated: RSAPKCS1v15 is not recommended.
// Please use RSA instead.
var RSAPKCS1v15 = &RSAAlgorithm{isPSS: false, isOAEP: false}

func NewRSA(safePrivateKeyBuf, safePublicKeyBuf *SafeBuf) (*Asymmetric, error) {
	return NewAsymmetric(RSA, safePrivateKeyBuf, safePublicKeyBuf)
}
func NewRSAndEraseKey(safePrivateKeyBuf, safePublicKeyBuf []byte) (*Asymmetric, error) {
	return NewAsymmetricWithoutEraseKey(RSA, safePrivateKeyBuf, safePublicKeyBuf, false)
}
func NewRSAWithOutEraseKey(safePrivateKeyBuf, safePublicKeyBuf []byte) (*Asymmetric, error) {
	return NewAsymmetricWithoutEraseKey(RSA, safePrivateKeyBuf, safePublicKeyBuf, false)
}

func GenerateRSAKeyPair(bitSize int) ([]byte, []byte, error) {
	if bitSize < 2048 {
		bitSize = 2048
	}
	priKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, nil, err
	}
	privateKey, err := x509.MarshalPKCS8PrivateKey(priKey)
	if err != nil {
		return nil, nil, err
	}
	publicKey, err := x509.MarshalPKIXPublicKey(&priKey.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, publicKey, nil
}

func (r *RSAAlgorithm) ParsePrivateKey(der []byte) (any, error) {
	keyAny, err := x509.ParsePKCS8PrivateKey(der)
	if err != nil {
		keyAny, err = x509.ParsePKCS1PrivateKey(der)
		if err != nil {
			return nil, err
		}
	}
	privKey, ok := keyAny.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("not an RSA private key")
	}
	return privKey, nil
}

func (r *RSAAlgorithm) ParsePublicKey(der []byte) (any, error) {
	pubKeyAny, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return nil, err
	}
	pubKey, ok := pubKeyAny.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}
	return pubKey, nil
}

func (r *RSAAlgorithm) Sign(privateKeyObj any, data []byte, hash ...crypto.Hash) ([]byte, error) {
	privKey, ok := privateKeyObj.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("invalid private key type for RSA")
	}
	if len(hash) == 0 {
		hash = append(hash, crypto.SHA256)
	}
	hasher := hash[0].New()
	hasher.Write(data)
	hashed := hasher.Sum(nil)
	if r.isPSS {
		return rsa.SignPSS(rand.Reader, privKey, hash[0], hashed, nil)
	}
	return rsa.SignPKCS1v15(rand.Reader, privKey, hash[0], hashed)
}

func (r *RSAAlgorithm) Verify(publicKeyObj any, data []byte, signature []byte, hash ...crypto.Hash) (bool, error) {
	pubKey, ok := publicKeyObj.(*rsa.PublicKey)
	if !ok {
		return false, errors.New("invalid public key type for RSA")
	}
	if len(hash) == 0 {
		hash = append(hash, crypto.SHA256)
	}
	hasher := hash[0].New()
	hasher.Write(data)
	hashed := hasher.Sum(nil)
	var err error
	if r.isPSS {
		err = rsa.VerifyPSS(pubKey, hash[0], hashed, signature, nil)
	} else {
		err = rsa.VerifyPKCS1v15(pubKey, hash[0], hashed, signature)
	}
	return err == nil, nil
}

func (r *RSAAlgorithm) Encrypt(publicKeyObj any, data []byte) ([]byte, error) {
	pubKey, ok := publicKeyObj.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("invalid public key type for RSA")
	}

	if r.isOAEP {
		return rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey, data, nil)
	}
	return rsa.EncryptPKCS1v15(rand.Reader, pubKey, data)
}

func (r *RSAAlgorithm) Decrypt(privateKeyObj any, data []byte) ([]byte, error) {
	privKey, ok := privateKeyObj.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("invalid private key type for RSA")
	}

	if r.isOAEP {
		return rsa.DecryptOAEP(sha256.New(), rand.Reader, privKey, data, nil)
	}
	return rsa.DecryptPKCS1v15(rand.Reader, privKey, data)
}

// ------ Ed25519 ------ //

type Ed25519Algorithm struct{}

var ED25519 = &Ed25519Algorithm{}

func NewED25519(safePrivateKeyBuf, safePublicKeyBuf *SafeBuf) (*Asymmetric, error) {
	return NewAsymmetric(ED25519, safePrivateKeyBuf, safePublicKeyBuf)
}
func NewED25519AndEraseKey(safePrivateKeyBuf, safePublicKeyBuf []byte) (*Asymmetric, error) {
	return NewAsymmetricWithoutEraseKey(ED25519, safePrivateKeyBuf, safePublicKeyBuf, false)
}
func NewED25519WithOutEraseKey(safePrivateKeyBuf, safePublicKeyBuf []byte) (*Asymmetric, error) {
	return NewAsymmetricWithoutEraseKey(ED25519, safePrivateKeyBuf, safePublicKeyBuf, false)
}

func GenerateEd25519KeyPair() ([]byte, []byte, error) {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return privKey, pubKey, nil
}

func (e *Ed25519Algorithm) ParsePrivateKey(der []byte) (any, error) {
	if len(der) != ed25519.PrivateKeySize {
		return nil, errors.New("invalid Ed25519 private key size")
	}
	return ed25519.PrivateKey(der), nil
}

func (e *Ed25519Algorithm) ParsePublicKey(der []byte) (any, error) {
	if len(der) != ed25519.PublicKeySize {
		return nil, errors.New("invalid Ed25519 public key size")
	}
	return ed25519.PublicKey(der), nil
}

func (e *Ed25519Algorithm) Sign(privateKeyObj any, data []byte, hash ...crypto.Hash) ([]byte, error) {
	privKey, ok := privateKeyObj.(ed25519.PrivateKey)
	if !ok {
		return nil, errors.New("invalid private key type for Ed25519")
	}
	return ed25519.Sign(privKey, data), nil
}

func (e *Ed25519Algorithm) Verify(publicKeyObj any, data []byte, signature []byte, hash ...crypto.Hash) (bool, error) {
	pubKey, ok := publicKeyObj.(ed25519.PublicKey)
	if !ok {
		return false, errors.New("invalid public key type for Ed25519")
	}
	return ed25519.Verify(pubKey, data, signature), nil
}

// ------ X25519 ------ //

type X25519Algorithm struct {
	useGCM  bool
	kdfInfo []byte
	kdfSalt []byte
	hash    crypto.Hash
}

var X25519GCM = &X25519Algorithm{useGCM: true, hash: crypto.SHA256}
var X25519CBC = &X25519Algorithm{useGCM: false, hash: crypto.SHA256}

func NewX25519Algorithm(useGCM bool, hash crypto.Hash, kdfInfo, kdfSalt []byte) *X25519Algorithm {
	return &X25519Algorithm{useGCM: useGCM, kdfInfo: kdfInfo, kdfSalt: kdfSalt, hash: hash}
}

func NewX25519(safePrivateKeyBuf, safePublicKeyBuf *SafeBuf) (*Asymmetric, error) {
	return NewAsymmetric(X25519GCM, safePrivateKeyBuf, safePublicKeyBuf)
}
func NewX25519AndEraseKey(safePrivateKeyBuf, safePublicKeyBuf []byte) (*Asymmetric, error) {
	return NewAsymmetricWithoutEraseKey(X25519GCM, safePrivateKeyBuf, safePublicKeyBuf, false)
}
func NewX25519WithOutEraseKey(safePrivateKeyBuf, safePublicKeyBuf []byte) (*Asymmetric, error) {
	return NewAsymmetricWithoutEraseKey(X25519GCM, safePrivateKeyBuf, safePublicKeyBuf, false)
}

func GenerateX25519KeyPair() ([]byte, []byte, error) {
	privKey := make([]byte, curve25519.ScalarSize)
	if _, err := rand.Read(privKey); err != nil {
		return nil, nil, err
	}
	pubKey, err := curve25519.X25519(privKey, curve25519.Basepoint)
	if err != nil {
		return nil, nil, err
	}
	return privKey, pubKey, nil
}

func (x *X25519Algorithm) ParsePrivateKey(der []byte) (any, error) {
	if len(der) != curve25519.ScalarSize {
		return nil, errors.New("invalid X25519 private key size")
	}
	key := make([]byte, curve25519.ScalarSize)
	copy(key, der)
	return key, nil
}

func (x *X25519Algorithm) ParsePublicKey(der []byte) (any, error) {
	if len(der) != curve25519.PointSize {
		return nil, errors.New("invalid X25519 public key size")
	}
	key := make([]byte, curve25519.PointSize)
	copy(key, der)
	return key, nil
}

func (x *X25519Algorithm) Sign(privateKeyObj any, data []byte, hash ...crypto.Hash) ([]byte, error) {
	return nil, errors.New("X25519 does not support signing, use Ed25519 instead")
}

func (x *X25519Algorithm) Verify(publicKeyObj any, data []byte, signature []byte, hash ...crypto.Hash) (bool, error) {
	return false, errors.New("X25519 does not support verification, use Ed25519 instead")
}

func (x *X25519Algorithm) Encrypt(publicKeyObj any, data []byte) ([]byte, error) {
	targetPub, ok := publicKeyObj.([]byte)
	if !ok || len(targetPub) != curve25519.PointSize {
		return nil, errors.New("invalid public key type/size for X25519")
	}

	ephemeralPriv, ephemeralPub, err := GenerateX25519KeyPair()
	if err != nil {
		return nil, err
	}
	defer ZeroMemory(ephemeralPriv)

	sharedSecret, err := curve25519.X25519(ephemeralPriv, targetPub)
	if err != nil {
		return nil, err
	}
	defer ZeroMemory(sharedSecret)

	hashFunc := x.hash
	if hashFunc == 0 {
		hashFunc = crypto.SHA256
	}
	hkdfReader := hkdf.New(hashFunc.New, sharedSecret, x.kdfSalt, x.kdfInfo)

	aesKey := make([]byte, 32)
	defer ZeroMemory(aesKey)

	if _, err := io.ReadFull(hkdfReader, aesKey); err != nil {
		return nil, err
	}

	cipherAlgo := &AESCipher{useGCM: x.useGCM}
	ivLen := 16
	if x.useGCM {
		ivLen = 12
	}
	iv := MakeSafeToken(ivLen)

	cipherText, err := cipherAlgo.Encrypt(data, aesKey, iv)
	if err != nil {
		return nil, err
	}

	out := make([]byte, 0, len(ephemeralPub)+len(iv)+len(cipherText))
	out = append(out, ephemeralPub...)
	out = append(out, iv...)
	out = append(out, cipherText...)

	return out, nil
}

func (x *X25519Algorithm) Decrypt(privateKeyObj any, data []byte) ([]byte, error) {
	myPriv, ok := privateKeyObj.([]byte)
	if !ok || len(myPriv) != curve25519.ScalarSize {
		return nil, errors.New("invalid private key type/size for X25519")
	}

	pubKeyLen := curve25519.PointSize
	ivLen := 16
	if x.useGCM {
		ivLen = 12
	}

	if len(data) < pubKeyLen+ivLen {
		return nil, errors.New("invalid ciphertext package size")
	}

	ephemeralPub := data[:pubKeyLen]
	iv := data[pubKeyLen : pubKeyLen+ivLen]
	cipherText := data[pubKeyLen+ivLen:]

	sharedSecret, err := curve25519.X25519(myPriv, ephemeralPub)
	if err != nil {
		return nil, err
	}
	defer ZeroMemory(sharedSecret)

	hashFunc := x.hash
	if hashFunc == 0 {
		hashFunc = crypto.SHA256
	}
	hkdfReader := hkdf.New(hashFunc.New, sharedSecret, x.kdfSalt, x.kdfInfo)

	aesKey := make([]byte, 32)
	defer ZeroMemory(aesKey)

	if _, err := io.ReadFull(hkdfReader, aesKey); err != nil {
		return nil, err
	}

	cipherAlgo := &AESCipher{useGCM: x.useGCM}
	plainText, err := cipherAlgo.Decrypt(cipherText, aesKey, iv)

	return plainText, err
}
