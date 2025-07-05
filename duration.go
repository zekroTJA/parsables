package parsables

import "time"

// Duration implements encoding.TextUnmarshaler for time.Duration.
//
// The implementation only wraps time.ParseDuration. For information
// about formatting, please refer to the documentation of time.ParseDuration.
type Duration struct {
	time.Duration
}

// DurationFromString parses Duration from the given string s.
// For format details, see documentation of Duration.
func DurationFromString(s string) (Duration, error) {
	d, err := time.ParseDuration(s)
	return Duration{d}, err
}

func (t *Duration) UnmarshalText(p []byte) (err error) {
	*t, err = DurationFromString(string(p))
	return err
}
