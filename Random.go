package u

import (
	"math/rand"
	"time"
)

var GlobalRand1 = rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
var GlobalRand2 = rand.New(rand.NewSource(int64(time.Now().Unix())))
