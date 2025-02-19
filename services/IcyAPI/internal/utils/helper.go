package utils

import (
	"time"

	logger "itsjaylen/IcyLogger"
)

func Retry(attempts int, delay time.Duration, fn func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		if err = fn(); err == nil {
			return nil 
		}
		logger.Warn.Printf("Retry attempt %d failed: %v", i+1, err)
		time.Sleep(delay)
	}
	logger.Error.Println("All retry attempts failed")
	return err
}
