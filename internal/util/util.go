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

func Find[T comparable](elems []T, p func(item T) bool) (T, bool) {
	var result T
	for _, v := range elems {
		if p(v) {
			return v, true
		}
	}
	return result, false
}

func Schedule(interval time.Duration, handler func()) func() {
	ticker := time.NewTicker(interval)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				handler()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	return func() {
		close(quit)
	}
}

func ScheduleControlled(interval time.Duration, handler func() bool) {
	ticker := time.NewTicker(interval)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				stop := handler()
				if stop {
					close(quit)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
