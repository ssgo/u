package u_test

import (
	"fmt"
	"github.com/ssgo/u"
	"testing"
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

func TestAes(t *testing.T) {

	testString := "Hello Password!"

	key := []byte("vpL54DlR2KG{JSAaAX7Tu;*#&DnG`M0o")
	iv := []byte("@z]zv@10-K.5Al0Dm`@foq9k\"VRfJ^~j")
	encrypted := u.EncryptAes(testString, key, iv)
	decrypted := u.DecryptAes(encrypted, key, iv)

	if decrypted != testString {
		t.Error("Decrypt failed", encrypted, decrypted)
	}
}
