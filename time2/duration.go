package time2

import (
	"errors"
	"time"
)

type Duration struct {
	time.Duration
}

func NewDuration(d time.Duration) Duration {
	return Duration{
		Duration: d,
	}
}

func (d *Duration) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.String() + `"`), nil
}

func (d *Duration) UnmarshalJSON(b []byte) (err error) {
	if len(b) < 2 {
		return errors.New("duration is not string")
	}
	if string(b) == `""` {
		return nil
	}
	d.Duration, err = time.ParseDuration(string(b[1 : len(b)-1]))
	return
}
