package u

import (
	"math/rand"
	"sync"
	"time"
)

var GlobalRand1 = NewRand(rand.NewSource(int64(time.Now().Nanosecond())))
var GlobalRand2 = NewRand(rand.NewSource(int64(time.Now().Nanosecond())))

func NewRand(source rand.Source) *Rand {
	return &Rand{
		goRand: rand.New(source),
	}
}

type Rand struct {
	goRand *rand.Rand
	lock   sync.Mutex
}

func (r *Rand) Seed(seed int64) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.goRand.Seed(seed)
}

func (r *Rand) Int63() int64 {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.goRand.Int63()
}

func (r *Rand) Uint32() uint32 {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.goRand.Uint32()
}

func (r *Rand) Uint64() uint64 {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.goRand.Uint64()
}

func (r *Rand) Int31() int32 {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.goRand.Int31()
}

func (r *Rand) Int() int {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.goRand.Int()
}

func (r *Rand) Int63n(n int64) int64 {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.goRand.Int63n(n)
}

func (r *Rand) Int31n(n int32) int32 {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.goRand.Int31n(n)
}

func (r *Rand) Intn(n int) int {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.goRand.Intn(n)
}

func (r *Rand) Float64() float64 {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.goRand.Float64()
}

func (r *Rand) Float32() float32 {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.goRand.Float32()
}

func (r *Rand) Perm(n int) []int {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.goRand.Perm(n)
}

func (r *Rand) Shuffle(n int, swap func(i, j int)) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.goRand.Shuffle(n, swap)
}

func (r *Rand) Read(p []byte) (n int, err error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.goRand.Read(p)
}
