// Package utils provides utility functions.
package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	logger "itsjaylen/IcyLogger"
)

// ErrInvalidDurationFormat is returned when a duration string format is invalid.
var ErrInvalidDurationFormat = errors.New("invalid duration format")

// Retry executes a function multiple times with exponential backoff.
func Retry(attempts int, delay time.Duration, fn func() error) error {
	var err error
	for i := range attempts {
		if err = fn(); err == nil {
			return nil
		}
		logger.Warn.Printf("Retry attempt %d failed: %v", i+1, err)
		time.Sleep(delay)
	}
	logger.Error.Println("All retry attempts failed")

	return err
}

// ParseDuration converts duration strings like "1d", "2h30m", "15m" into time.Duration.
func ParseDuration(durationStr string) (time.Duration, error) {
	re := regexp.MustCompile(`(\d+)([smhd])`)
	matches := re.FindAllStringSubmatch(durationStr, -1)

	if len(matches) == 0 {
		return 0, fmt.Errorf("%w: %s", ErrInvalidDurationFormat, durationStr)
	}

	var totalDuration time.Duration
	unitMap := map[string]time.Duration{
		"s": time.Second,
		"m": time.Minute,
		"h": time.Hour,
		"d": time.Hour * 24,
	}

	for _, match := range matches {
		value, _ := strconv.Atoi(match[1])
		unit := match[2]
		totalDuration += time.Duration(value) * unitMap[unit]
	}

	return totalDuration, nil
}
