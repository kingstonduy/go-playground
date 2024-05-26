package utils_class

import (
	"fmt"
	"strconv"
	"time"
)

type TimeUtils struct {
	now time.Time
}

func NewTimeUtils() TimeUtils {
	return TimeUtils{
		now: time.Now(),
	}
}

func (t *TimeUtils) GetCurrentTime() {
	t.now = time.Now()
}

func (t *TimeUtils) Format(now time.Time, format string) string {
	return now.Format(format)
}

func (t *TimeUtils) GetTimeYYYYMMDDHHMMSS() string {
	return t.now.Format("20060102150405") // YYYYMMDDHHMMSS
}

func (t *TimeUtils) GetTimeGmtMMDDHHMMSS() string {
	return t.now.In(time.FixedZone("GMT", 0)).Format("0102150405") // MMDDHHMMSS
}

func (t *TimeUtils) GetTimeHHMMSS() string {
	return t.now.Format("150405") // HHMMSS
}

func (t *TimeUtils) GetTimeMMDD() string {
	return t.now.Format("0102") // MMDD
}

func (t *TimeUtils) GetTimeYDDDHHNNNNNN(f11 string) string {
	lastLetterYear := strconv.Itoa(t.now.Year() % 10)
	dayOfYear := strconv.Itoa(t.now.YearDay())

	f37 := fmt.Sprintf("%s%s%s%s", lastLetterYear, dayOfYear, t.now.In(time.FixedZone("GMT", 0)).Format("15"), f11)
	return f37
}
