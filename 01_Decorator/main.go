package main

import (
	"fmt"
	"math"
	"reflect"
	"sync"
	"time"
)

func Pi(n int) float64 {
	ch := make(chan float64)
	defer close(ch)

	for k := 0; k < n; k++ {
		go func(ch chan float64, k float64) {
			ch <- 4 * math.Pow(-1, k) / (2*k + 1)
		}(ch, float64(k))
	}

	pi := 0.0
	for k := 0; k < n; k++ {
		pi += <-ch
	}
	return pi
}

func wrapperLogger(f func(int) float64) func(int) float64 {
	return func(n int) float64 {
		fu := func(ns int) (res float64) {
			defer func(t time.Time) {
				fmt.Printf("toke=%v, n=%d, pi=%v\n", time.Since(t), n, res)
			}(time.Now())

			return f(n)
		}
		return fu(n)
	}
}

func wrapperCache(f func(int) float64, cache *sync.Map) func(int) float64 {
	return func(n int) float64 {
		name := fmt.Sprintf("%s_%d", reflect.ValueOf(&f).Type().Name(), n)
		if v, ok := cache.Load(name); ok {
			fmt.Printf("find in cache(name=%s) n = %d, pi = %v\n", name, n, v)
			return v.(float64)
		}

		res := f(n)
		cache.Store(name, res)
		return res
	}
}

func divide(n int) float64 {
	return float64(n) / 2.0
}

func main() {
	cache := sync.Map{}
	wrapperCache(wrapperLogger(Pi), &cache)(100000)
	wrapperCache(wrapperLogger(Pi), &cache)(200000)
	wrapperCache(wrapperLogger(Pi), &cache)(800000)
	wrapperCache(wrapperLogger(Pi), &cache)(200000)

	divideCache := sync.Map{}
	wrapperCache(wrapperLogger(divide), &divideCache)(200000)
	wrapperCache(wrapperLogger(divide), &divideCache)(100000)
	wrapperCache(wrapperLogger(divide), &divideCache)(500000)
	wrapperCache(wrapperLogger(divide), &divideCache)(200000)
}
