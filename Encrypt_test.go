package utility

import (
	"fmt"
	"testing"
)

func TestEncodeInt(t *testing.T) {
	fmt.Println(string(EncodeInt(GlobalRand2.Uint64())))
	for i := 0; i < 100000; i++ {
		n := GlobalRand2.Uint64()
		s := EncodeInt(uint64(n))
		n2 := DecodeInt(s)
		if n2 != n {
			t.Error("decode not match ", s, n, n2)
		}
	}
}


func TestAes(t *testing.T) {

	testString := "Hello Password!"

	key := []byte("vpL54DlR2KG{JSAaAX7Tu;*#&DnG`M0o")
	iv := []byte("@z]zv@10-K.5Al0Dm`@foq9k\"VRfJ^~j")
	encrypted := EncryptAes(testString, key, iv)
	decrypted := DecryptAes(encrypted, key, iv)

	if decrypted != testString {
		t.Error("Decrypt failed", encrypted, decrypted)
	}
}
