package utility

import (
	"fmt"
	"testing"
)

func TestUniqueId(t *testing.T) {
	fmt.Println(UniqueId())
	uids := map[string]bool{}
	for i := 0; i < 100000; i++ {
		uid := UniqueId()
		if uids[uid] {
			t.Error("unique id repeated ", uids, uid)
		}
		uids[uid] = true
	}
}
