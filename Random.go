package u

import (
	"math/rand/v2"
	"sync"
	"sync/atomic"
	"time"
)

type Rand struct {
	isPCG bool
}

var GlobalRand1 = &Rand{}
var GlobalRand2 = &Rand{isPCG: true}

var seedCounter atomic.Uint64
var pcgPool = sync.Pool{
	New: func() any {
		return rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), seedCounter.Add(1)))
	},
}

func (r *Rand) Int64() int64 {
	if !r.isPCG {
		return rand.Int64()
	}
	instance := pcgPool.Get().(*rand.Rand)
	v := instance.Int64()
	pcgPool.Put(instance)
	return v
}

func (r *Rand) Uint32() uint32 {
	if !r.isPCG {
		return rand.Uint32()
	}
	instance := pcgPool.Get().(*rand.Rand)
	v := instance.Uint32()
	pcgPool.Put(instance)
	return v
}

func (r *Rand) Uint64() uint64 {
	if !r.isPCG {
		return rand.Uint64()
	}
	instance := pcgPool.Get().(*rand.Rand)
	v := instance.Uint64()
	pcgPool.Put(instance)
	return v
}

func (r *Rand) Int32() int32 {
	if !r.isPCG {
		return rand.Int32()
	}
	instance := pcgPool.Get().(*rand.Rand)
	v := instance.Int32()
	pcgPool.Put(instance)
	return v
}

func (r *Rand) Int() int {
	if !r.isPCG {
		return rand.Int()
	}
	instance := pcgPool.Get().(*rand.Rand)
	v := instance.Int()
	pcgPool.Put(instance)
	return v
}

// Deprecated: use Int64N instead
func (r *Rand) Int63n(n int64) int64 {
	return r.Int64N(n)
}

func (r *Rand) Int64N(n int64) int64 {
	if !r.isPCG {
		return rand.Int64N(n)
	}
	instance := pcgPool.Get().(*rand.Rand)
	v := instance.Int64N(n)
	pcgPool.Put(instance)
	return v
}

// Deprecated: use IntN instead
func (r *Rand) Int31n(n int32) int32 {
	return r.Int32N(n)
}

func (r *Rand) Int32N(n int32) int32 {
	if !r.isPCG {
		return rand.Int32N(n)
	}
	instance := pcgPool.Get().(*rand.Rand)
	v := instance.Int32N(n)
	pcgPool.Put(instance)
	return v
}

// Deprecated: use IntN instead
func (r *Rand) Intn(n int) int {
	return r.IntN(n)
}

func (r *Rand) IntN(n int) int {
	if !r.isPCG {
		return rand.IntN(n)
	}
	instance := pcgPool.Get().(*rand.Rand)
	v := instance.IntN(n)
	pcgPool.Put(instance)
	return v
}

func (r *Rand) UintN(n uint) uint {
	if !r.isPCG {
		return rand.UintN(n)
	}
	instance := pcgPool.Get().(*rand.Rand)
	v := instance.UintN(n)
	pcgPool.Put(instance)
	return v
}

func (r *Rand) Uint64N(n uint64) uint64 {
	if !r.isPCG {
		return rand.Uint64N(n)
	}
	instance := pcgPool.Get().(*rand.Rand)
	v := instance.Uint64N(n)
	pcgPool.Put(instance)
	return v
}

func (r *Rand) Uint32N(n uint32) uint32 {
	if !r.isPCG {
		return rand.Uint32N(n)
	}
	instance := pcgPool.Get().(*rand.Rand)
	v := instance.Uint32N(n)
	pcgPool.Put(instance)
	return v
}

func (r *Rand) Float64() float64 {
	if !r.isPCG {
		return rand.Float64()
	}
	instance := pcgPool.Get().(*rand.Rand)
	v := instance.Float64()
	pcgPool.Put(instance)
	return v
}

func (r *Rand) Float32() float32 {
	if !r.isPCG {
		return rand.Float32()
	}
	instance := pcgPool.Get().(*rand.Rand)
	v := instance.Float32()
	pcgPool.Put(instance)
	return v
}

func (r *Rand) Perm(n int) []int {
	if !r.isPCG {
		return rand.Perm(n)
	}
	instance := pcgPool.Get().(*rand.Rand)
	v := instance.Perm(n)
	pcgPool.Put(instance)
	return v
}

func (r *Rand) Shuffle(n int, swap func(i, j int)) {
	if !r.isPCG {
		rand.Shuffle(n, swap)
		return
	}
	instance := pcgPool.Get().(*rand.Rand)
	instance.Shuffle(n, swap)
	pcgPool.Put(instance)
}
