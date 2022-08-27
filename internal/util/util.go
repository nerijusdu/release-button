package util

import "time"

func Contains[T comparable](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func Schedule(interval time.Duration, handler func()) func() {
	ticker := time.NewTicker(interval)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <- ticker.C:
				handler()
			case <- quit:
				ticker.Stop()
				return
			}
		}
	}()

	return func() {
		close(quit)
	}
}