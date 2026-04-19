package u_test

import (
	"fmt"
	"testing"

	"github.com/ssgo/u"
)

//func TestId(t *testing.T) {
//	for _, nx := range []string{"99u", "UUU", "999u", "UUUU", "9999u", "UUUUU", "99999u", "UUUUUU"} {
//		fmt.Println(nx, u.DecodeInt([]byte(nx)))
//	}
//	for i:=0; i<10; i++ {
//		dd := u.ShortUniqueId()
//		fmt.Println(dd, len(dd))
//	}
//}

func TestUniqueId(t *testing.T) {
	fmt.Println(u.MakeId(12))
	uids := map[string]bool{}
	for range 100000 {
		uid := u.MakeId(12)
		if uids[uid] || len(uid) != 12 {
			t.Error("unique id repeated ", uid, len(uid))
		}
		uids[uid] = true
	}
}

func TestShortUniqueId(t *testing.T) {
	fmt.Println(u.MakeId(20), len(u.MakeId(20)))
	uids := map[string]bool{}
	for range 100000 {
		uid := u.MakeId(20)
		if uids[uid] || len(uid) != 20 {
			t.Error("short unique id repeated ", uid, len(uid))
		}
		uids[uid] = true
	}
}
