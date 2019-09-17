package rnd

import (
	"math/rand"
	"time"

	"github.com/mediacoin-pro/core/common/sys"
)

type TimeDuration struct {
	MinValue time.Duration // minimum time duration
	KDisp    float64       // k-dispersion
}

func (r TimeDuration) Duration() time.Duration {
	return Duration(r.MinValue, r.KDisp)
}

func (r TimeDuration) MaxValue() time.Duration {
	return time.Duration(float64(r.MinValue) * (1 + r.KDisp))
}

func Duration(minValue time.Duration, kDispersion float64, seed ...interface{}) time.Duration {
	return time.Duration(float64(minValue) * (1 + kDispersion*New(seed).Float64()))
}

func Sleep(minValue time.Duration, kDispersion float64) {
	sys.Sleep(Duration(minValue, kDispersion))
}

func String(ss ...string) string {
	return ss[rand.Intn(len(ss))]
}
