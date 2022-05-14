package time2

import (
	"time"

	"github.com/pingcap/errors"
)

//Duration  持续时间，是time.Duration的封装
type Duration struct {
	time.Duration
}

//NewDuration 通过time.Duration的d获得持续时间
func NewDuration(d time.Duration) Duration {
	return Duration{
		Duration: d,
	}
}

//MarshalJSON 序列化JSON,json的序列化器使用
func (d *Duration) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.String() + `"`), nil
}

//UnmarshalJSON 反序列化JSON,json的序列化器使用
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
