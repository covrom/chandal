package verr

import (
	"errors"
	"time"
)

var TimeoutErr = errors.New("timeout")

func SendTimeout[T any](result chan T, v T, timeout time.Duration) error {
	tout := time.NewTimer(timeout)
	select {
	case result <- v:
	case <-tout.C:
		return TimeoutErr
	}
	if !tout.Stop() {
		<-tout.C
	}
	return nil
}
