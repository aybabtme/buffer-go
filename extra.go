package buffer

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const d = false

// UnixTime is a time.Time that marshals and unmarshals to Unix time.
type UnixTime struct {
	time.Time
}

func (us *UnixTime) UnmarshalJSON(p []byte) error {
	var i int64
	if err := json.Unmarshal(p, &i); err != nil {
		return err
	}
	us.Time = time.Unix(i, 0)
	return nil
}

func (us UnixTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(us.Time.Unix())
}

// Weekday as known to Buffer.
type Weekday struct {
	time.Weekday
}

var (
	prefixSunday    = strings.ToLower(time.Sunday.String()[:3])
	prefixMonday    = strings.ToLower(time.Monday.String()[:3])
	prefixTuesday   = strings.ToLower(time.Tuesday.String()[:3])
	prefixWednesday = strings.ToLower(time.Wednesday.String()[:3])
	prefixThursday  = strings.ToLower(time.Thursday.String()[:3])
	prefixFriday    = strings.ToLower(time.Friday.String()[:3])
	prefixSaturday  = strings.ToLower(time.Saturday.String()[:3])
)

func (wd *Weekday) UnmarshalJSON(p []byte) error {
	var s string
	if err := json.Unmarshal(p, &s); err != nil {
		return err
	}
	switch s {
	case prefixSunday:
		wd.Weekday = time.Sunday
	case prefixMonday:
		wd.Weekday = time.Monday
	case prefixTuesday:
		wd.Weekday = time.Tuesday
	case prefixWednesday:
		wd.Weekday = time.Wednesday
	case prefixThursday:
		wd.Weekday = time.Thursday
	case prefixFriday:
		wd.Weekday = time.Friday
	case prefixSaturday:
		wd.Weekday = time.Saturday
	default:
		return errors.Errorf("unknown weekday: %q", s)
	}
	return nil
}

func (wd Weekday) MarshalJSON() ([]byte, error) {
	var prefix string
	switch wd.Weekday {
	case time.Sunday:
		prefix = prefixSunday
	case time.Monday:
		prefix = prefixMonday
	case time.Tuesday:
		prefix = prefixTuesday
	case time.Wednesday:
		prefix = prefixWednesday
	case time.Thursday:
		prefix = prefixThursday
	case time.Friday:
		prefix = prefixFriday
	case time.Saturday:
		prefix = prefixSaturday
	default:
		return nil, errors.Errorf("unknown weekday: %q", wd)
	}
	return json.Marshal(prefix)
}

// Daytime is a time during a day. The date does not matter, only the hours
// and minutes.
type Daytime struct {
	time.Time
}

func (dt *Daytime) UnmarshalJSON(p []byte) error {
	var str string
	err := json.Unmarshal(p, &str)
	if err != nil {
		return err
	}
	dt.Time, err = time.Parse("15:04", str)
	return err
}

func (dt Daytime) MarshalJSON() ([]byte, error) {
	t := dt.Time
	return json.Marshal(fmt.Sprintf("%02d:%02d", t.Hour(), t.Minute()))
}

// Timezone is a timezone in the world.
type Timezone struct {
	*time.Location
}

func (tz *Timezone) UnmarshalJSON(p []byte) error {
	var str string
	err := json.Unmarshal(p, &str)
	if err != nil {
		return err
	}
	tz.Location, err = time.LoadLocation(str)
	return err
}

func (tz Timezone) MarshalJSON() ([]byte, error) {
	return json.Marshal(tz.Location.String())
}
