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
