package uid

import (
	"sync"
	"time"
)

// NewV7 constructs a new v7 UUID. Enforces method 3 of monotonicity.
func NewV7() UUID { return make7(tick) }

// Time returns the embedded timestamp of UUID. For non-V7 zero(time.Time) is returned. If you don't pre-check version
// use `.IsZero()` to ensure time is "real".
//
//nolint:mnd // locality of behavior
func (u UUID) Time() time.Time {
	if u.Version() != Version7 {
		return time.Time{}
	}
	// rebuild unix_ts_ms
	ms := int64(u.b[0])<<40 | int64(u.b[1])<<32 | int64(u.b[2])<<24 | int64(u.b[3])<<16 | int64(u.b[4])<<8 | int64(u.b[5])
	ra := int64(u.b[6]&0x0f)<<8 | // top 4 of rand_a
		int64(u.b[7]) // bottom 8 of rand_a
	ns := ra * 1000 * 1000 / 4096
	return time.UnixMilli(ms).Add(time.Nanosecond * time.Duration(ns))
}

//nolint:mnd // locality of behavior
func make7(tickFn func() (int64, int64)) UUID {
	var b [16]byte
	ms, ra := tickFn()
	if ms < 0 {
		panic("v7 UUID does not support time before epoch")
	}
	// set unix_ts_ms
	b[0], b[1], b[2], b[3], b[4], b[5] = byte(ms>>40), byte(ms>>32), byte(ms>>24), byte(ms>>16), byte(ms>>8), byte(ms)
	// set rand_a
	b[6] = byte((ra >> 8)) & 0x0f // set top 4 bytes of rand_a
	b[7] = byte(ra)
	// fill rand_b
	_, _ = rng.Read(b[8:]) //nolint:errcheck // never returns errors
	// version, variant
	b[6], b[8] = (b[6]&0x0f)|0x70, (b[8]&0x3f)|0x80
	return UUID{b}
}

//nolint:mnd // locality of behavior
func msranda(t time.Time) (int64, int64) {
	ms := t.UnixMilli()
	nsr := t.UnixNano() - ms*1000*1000 // ns remainder
	return ms, int64(float64(nsr) * 4096.0 / 1000.0 / 1000.0)
}

func tick() (int64, int64) { return msranda(now()) }

/*
NewV7Strict returns a v7 UUID with guaranteed (beyond RFC method 3) local monotonicity.
If you think you need this, you don't. If you know you need this, your design is bad.
*/
func NewV7Strict() UUID { return make7(tickBatch) }

//nolint:gochecknoglobals // unexported
var (
	mux              = new(sync.Mutex)
	lastms, lasttick = msranda(now())
)

//nolint:nonamedreturns // golf
func tickBatch() (ms, tick int64) {
	mux.Lock()
	defer mux.Unlock()
	for {
		ms, tick = msranda(now())
		if ms > lastms {
			// now ms newer than last ms
			lastms, lasttick = ms, tick
			return
		}
		if ms == lastms && tick > lasttick {
			// in same ms, but tick is newer than last tick
			lasttick = tick
			return
		}
	}
}
