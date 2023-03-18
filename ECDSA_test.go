package u_test

import (
	"fmt"
	"github.com/ssgo/u"
	"testing"
)

func TestSign(t *testing.T) {

	priS, pubS, err := u.GenECDSA256Key()
	if err != nil {
		t.Fatal(err.Error())
	}
	priK, err := u.MakeECDSA256PrivateKey(priS)
	if err != nil {
		t.Fatal(err.Error())
	}
	pubK, err := u.MakeECDSA256PublicKey(pubS)
	if err != nil {
		t.Fatal(err.Error())
	}

	okCount := 0
	testCount := 1000
	for i:=0; i<testCount; i++ {
		data := u.MakeToken(10)
		sign, err := u.SignECDSA(data, priK)
		if err != nil {
			t.Fatal(err.Error())
		}
		if u.VerifyECDSA(data, sign, pubK) {
			okCount ++
		}
	}
	if okCount != testCount {
		t.Fatal("VerifyECDSA Error "+u.String(okCount))
	}
	fmt.Println(okCount, testCount)
}
