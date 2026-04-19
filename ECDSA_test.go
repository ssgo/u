package u_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/ssgo/u"
)

func TestSign(t *testing.T) {

	priKeyBuf, pubKeyBuf, err := u.GenerateECDSAKeyPair(521)
	if err != nil {
		t.Fatal(err.Error())
	}
	ecdsa, _ := u.NewECDSAndEraseKey(priKeyBuf, pubKeyBuf)

	okCount := 0
	testCount := 100
	t1 := time.Now()

	for i := 0; i < testCount; i++ {
		data := u.MakeToken(10)
		sign, err := ecdsa.Sign(data)
		if err != nil {
			t.Fatal(err.Error())
		}
		if ok, err := ecdsa.Verify(data, sign); ok && err == nil {
			okCount++
		}
	}
	if okCount != testCount {
		t.Fatal("VerifyECDSA Error " + u.String(okCount))
	}
	t2 := time.Now()
	fmt.Println(okCount, testCount, t2.Sub(t1))

	data := u.MakeToken(10)
	crypted, err := ecdsa.Encrypt(data)
	if err != nil {
		t.Fatal(err)
	}
	decryptedData, err := ecdsa.Decrypt(crypted)
	t3 := time.Now()
	fmt.Println("DecryptECDSA Time:", t3.Sub(t2))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(decryptedData, data) {
		t.Fatal("DecryptECDSA Error", decryptedData, data)
	}

}
