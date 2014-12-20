package gocostmodel

import (
	"math"
	"testing"
)

var (
	x         int
	val       = 1
	rng       = []int{1}
	blockedCh = make(chan bool)
	closedCh  = make(chan bool)
	fullBufCh = make(chan bool, 1)
)

func init() {
	close(closedCh)
	fullBufCh <- true
}

func BenchmarkSwitch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		switch val {
		case 0:
			x++
		case 1:
			x--
		default:
			x += 2
		}
	}
}

func BenchmarkIfEquals(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if val == 1 {
			x++
		}
	}
}

func BenchmarkIfNotEquals(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if val != 5 {
			x++
		}
	}
}

func BenchmarkFor1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for val < 0 {
			x++
		}
	}
}

func BenchmarkFor3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// loop only once, the goal is to measure the loop construct
		for j := 0; j < val; j++ {
			x++
		}
	}
}

func BenchmarkForRange(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// do not use empty range, would force go1.4
		for _ = range rng {
			x++
		}
	}
}

func BenchmarkForRangeClosedChan(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// do not use empty range, would force go1.4
		for _ = range closedCh {
			x++
		}
	}
}

func BenchmarkSelectBlockedDefault(b *testing.B) {
	for i := 0; i < b.N; i++ {
		select {
		case <-blockedCh:
			x--
		default:
			x++
		}
	}
}

func BenchmarkSelectBlockedClosed(b *testing.B) {
	for i := 0; i < b.N; i++ {
		select {
		case <-blockedCh:
			x--
		case <-closedCh:
			x++
		}
	}
}

func BenchmarkSelectTrySend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		select {
		case blockedCh <- false:
		default:
			x++
		}
	}
}

func BenchmarkPanicRecover(b *testing.B) {
	for i := 0; i < b.N; i++ {
		panicRecover()
	}
}

func panicRecover() {
	defer func() {
		// a simple call to recover would be enough, but the usual way
		// this pattern is used is to do something conditionnally if there
		// was a panic.
		if e := recover(); e != nil {
			x++
		}
	}()
	panic("zomg!")
}

func BenchmarkSelectTrySendBuf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		select {
		case fullBufCh <- false:
		default:
			x++
		}
	}
}

func BenchmarkFunc1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		x = sum1(i1)
	}
}

func BenchmarkFunc2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		x = sum2(i1, i2)
	}
}

func BenchmarkFunc3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		x = sum3(i1, i2, i3)
	}
}

func sum1(n int) int {
	if n > 0 {
		return n
	}

	// instructions just to avoid inlining
	math.Pow(float64(n), 1)
	panic("unreachable")
}

func sum2(n1, n2 int) int {
	if n1 > 0 {
		x := n1 + n2
		return x
	}

	// instructions just to avoid inlining
	math.Pow(float64(n1), 1)
	panic("unreachable")
}

func sum3(n1, n2, n3 int) int {
	if n1 > 0 {
		x := n1 + n2 + n3
		return x
	}

	// instructions just to avoid inlining
	math.Pow(float64(n1), 1)
	panic("unreachable")
}
