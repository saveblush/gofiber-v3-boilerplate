package utils

import (
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/config"
)

// Now now
func Now() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.Local)
}

// CurTime current time
func CurTime() string {
	return Now().Format(config.CF.Web.TimeFormat)
}

// DateFormat date format
func DateFormat(t time.Time) string {
	return t.Format(config.CF.Web.DateFormat)
}

// DateTimeFormat date time format
func DateTimeFormat(t time.Time) string {
	return t.Format(config.CF.Web.DateTimeFormat)
}

// TimeFormat time format
func TimeFormat(t time.Time) string {
	return t.Format(config.CF.Web.TimeFormat)
}

// DateAndTimeAsDateTimeFormat date and time as date time format
func DateAndTimeAsDateTimeFormat(d time.Time, t string) time.Time {
	dt := strings.Split(t, ":")
	hour, _ := strconv.Atoi(dt[0])
	minute, _ := strconv.Atoi(dt[1])
	second, _ := strconv.Atoi(dt[2])

	return time.Date(d.Year(), d.Month(), d.Day(), hour, minute, second, 0, time.Local)
}

// DateAdd date add
// DateAdd(Now(), 3)
func DateAdd(t time.Time, amt int) time.Time {
	return t.Add(time.Duration(amt) * 24 * time.Hour)
}

// DateSub date sub
// DateSub(Now(), 3)
func DateSub(t time.Time, amt int) time.Time {
	return t.Add(time.Duration(-amt) * 24 * time.Hour)
}

// DateAddDuration date add duration
// DateAddDuration(Now(), time.Duration)
func DateAddDuration(t time.Time, d time.Duration) time.Time {
	return t.Add(d)
}

// DateDiff date diff
func DateDiff(s, e time.Time) time.Duration {
	first := time.Date(e.Year(), e.Month(), e.Day(), e.Hour(), e.Minute(), e.Second(), 0, time.Local)
	second := time.Date(s.Year(), s.Month(), s.Day(), s.Hour(), s.Minute(), s.Second(), 0, time.Local)

	return first.Sub(second)
}

// DateDiffFormat date diff format
func DateDiffFormat(a, b time.Time) (year, month, day, hour, min, sec int) {
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}

// TimeParseDuration time parse duration
// ตย. t = "1m15s"
func TimeParseDuration(t string) time.Duration {
	tt, _ := time.ParseDuration(t)
	return tt
}

// TimeTtl time ttl
// แปลงเป็น Milliseconds
func TimeTtl(t string) int64 {
	tt := TimeParseDuration(t)
	return tt.Milliseconds()
}

// TimeParse time parse
func TimeParse(layout, value string) time.Time {
	tt, _ := time.Parse(layout, value)
	return tt
}

// TimeZone time zone
func TimeZone() string {
	zone, _ := Now().Zone()
	return zone
}
