package u

import (
	"crypto/rand"
	"encoding/binary"
	"runtime"
	"sync"
	"time"
	"unsafe"

	"golang.org/x/crypto/chacha20"
)

func MakeSafeToken(size int) []byte {
	fixsize := GlobalRand1.IntN(size/10+2) + 1
	key := make([]byte, size+fixsize*2)
	rand.Read(key)
	return key[fixsize : fixsize+size]
}

func EncryptChaCha20(raw []byte, key []byte, salt []byte) []byte {
	cipherObj, _ := chacha20.NewUnauthenticatedCipher(key[:chacha20.KeySize], salt[:chacha20.NonceSizeX])
	cipher := make([]byte, len(raw))
	cipherObj.XORKeyStream(cipher, raw)
	return cipher
}

func DecryptChaCha20(cipher []byte, key []byte, salt []byte) []byte {
	cipherObj, _ := chacha20.NewUnauthenticatedCipher(key[:chacha20.KeySize], salt[:chacha20.NonceSizeX])
	plaintext := make([]byte, len(cipher))
	cipherObj.XORKeyStream(plaintext, cipher)
	return plaintext
}

func ZeroMemory(buf []byte) {
	seed := uint64(time.Now().UnixNano())
	i := 0
	n := len(buf)
	for i <= n-8 {
		seed ^= seed << 13
		seed ^= seed >> 7
		seed ^= seed << 17
		binary.LittleEndian.PutUint64(buf[i:], seed)
		i += 8
	}
	if i < n {
		seed ^= seed << 13
		seed ^= seed >> 7
		seed ^= seed << 17
		var tmp [8]byte
		binary.LittleEndian.PutUint64(tmp[:], seed)
		copy(buf[i:], tmp[:n-i])
	}
	runtime.KeepAlive(buf)
}

var chachaGlobalKey = make([]byte, chacha20.KeySize)

func init() {
	rand.Read(chachaGlobalKey)
}

var safeBufEncrypt = func(raw []byte) ([]byte, []byte) {
	salt := MakeSafeToken(chacha20.NonceSizeX)
	return EncryptChaCha20(raw, chachaGlobalKey, salt), salt
}

var safeBufDecrypt = func(cipher []byte, salt []byte) []byte {
	return DecryptChaCha20(cipher, chachaGlobalKey, salt)
}

var setObfOnce sync.Once

func SetSafeBufObfuscator(encrypt func([]byte) ([]byte, []byte), decrypt func([]byte, []byte) []byte) {
	setObfOnce.Do(func() {
		safeBufEncrypt = encrypt
		safeBufDecrypt = decrypt
	})
}

type SafeBuf struct {
	buf  []byte
	salt []byte
}

type SecretPlaintext struct {
	Data []byte
}

func (sp *SecretPlaintext) String() string {
	if len(sp.Data) == 0 {
		return ""
	}
	return unsafe.String(&sp.Data[0], len(sp.Data))
}

func (sp *SecretPlaintext) Close() {
	if sp.Data != nil {
		ZeroMemory(sp.Data)
		sp.Data = nil
	}
}

func NewSafeBuf(raw []byte) *SafeBuf {
	cipher, salt := safeBufEncrypt(raw)
	return &SafeBuf{buf: cipher, salt: salt}
}

func MakeSafeBuf(cipher, salt []byte) *SafeBuf {
	return &SafeBuf{buf: cipher, salt: salt}
}

func (sb *SafeBuf) Open() *SecretPlaintext {
	sp := &SecretPlaintext{Data: safeBufDecrypt(sb.buf, sb.salt)}
	runtime.SetFinalizer(sp, func(obj *SecretPlaintext) {
		obj.Close()
	})
	return sp
}

func (sb *SafeBuf) Close() {
	ZeroMemory(sb.buf)
	ZeroMemory(sb.salt)
}

func NewSafeString(raw []byte) (*SecretPlaintext, string) {
	sp := &SecretPlaintext{Data: raw}
	runtime.SetFinalizer(sp, func(obj *SecretPlaintext) {
		obj.Close()
	})
	return sp, sp.String()
}
