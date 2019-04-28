package u_test

import (
	"fmt"
	"github.com/ssgo/u"
	"testing"
)

func TestUniqueId(t *testing.T) {
	fmt.Println(u.UniqueId())
	uids := map[string]bool{}
	for i := 0; i < 100000; i++ {
		uid := u.UniqueId()
		if uids[uid] {
			t.Error("unique id repeated ", uids, uid)
		}
		uids[uid] = true
	}
}

func TestShortUniqueId(t *testing.T) {
	fmt.Println(u.ShortUniqueId(), len(u.ShortUniqueId()))
	uids := map[string]bool{}
	for i := 0; i < 100000; i++ {
		uid := u.UniqueId()
		if uids[uid] {
			t.Error("short unique id repeated ", uids, uid)
		}
		uids[uid] = true
	}
}
