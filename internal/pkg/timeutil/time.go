package timeutil

import "time"

func NowRFC3339() string {
	return time.Now().Format(time.RFC3339)
}

func MustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		return time.Local
	}
	return loc
}

func ToLocation(t time.Time, location string) time.Time {
	return t.In(MustLoadLocation(location))
}
