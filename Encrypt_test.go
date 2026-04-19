package u_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/ssgo/u"
)

func TestEncodeInt(t *testing.T) {
	fmt.Println(string(u.EncodeInt(u.GlobalRand2.Uint64())))
	for i := 0; i < 100000; i++ {
		n := u.GlobalRand2.Uint64()
		s := u.EncodeInt(uint64(n))
		n2 := u.DecodeInt(s)
		if n2 != n {
			t.Error("decode not match ", s, n, n2)
		}
	}
}

func TestEncodeSha(t *testing.T) {
	fmt.Println(u.MD5String("Hello"))
	fmt.Println(u.MD5Base64("Hello"))
	fmt.Println(u.Sha1String("Hello"))
	fmt.Println(u.Sha1Base64("Hello"))
	fmt.Println(u.Sha256String("Hello"))
	fmt.Println(u.Sha256Base64("Hello"))
	fmt.Println(u.Sha512String("Hello"))
	fmt.Println(u.Sha512Base64("Hello"))
	fmt.Println(u.Hex(u.HmacSha256([]byte("key"), []byte("data"), []byte("data2"))))
}

func TestAes(t *testing.T) {
	testData := []byte("Hello Password!")
	defer u.ZeroMemory(testData)
	testBuf := u.NewSafeBuf(testData)

	aes, _ := u.NewAESGCMAndEraseKey([]byte("vpL54DlR2KG{JSAaAX7Tu;*#&DnG`M0o"), []byte("@z]zv@10-K.5Al0Dm`@foq9k\"VRfJ^~j"))
	encrypted, _ := aes.Encrypt(testBuf)
	decrypted, _ := aes.Decrypt(encrypted)

	decryptedBuf := decrypted.Open()
	defer decryptedBuf.Close()
	if !bytes.Equal(decryptedBuf.Data, testData) {
		t.Error("Decrypt failed", encrypted, decrypted)
	}
}
