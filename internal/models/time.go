package models

import (
	"strings"
	"time"
)

// TickTickTime
type TickTickTime struct {
	time.Time
}

const (
	TickTickTimeFormat  = "2006-01-02T15:04:05.000-0700"
	TickTickTimeFormat2 = "2006-01-02T15:04:05-0700"
	TickTickTimeFormat3 = "2006-01-02T15:04:05.000+0000"
	TickTickTimeFormat4 = "2006-01-02T15:04:05+0000"
)

func (t *TickTickTime) UnmarshalJSON(data []byte) error {
	if len(data) <= 2 {
		return nil
	}

	// trim quotes and check for null or empty string
	str := strings.Trim(string(data), `"`)
	if str == "" || str == "null" {
		return nil
	}

	formats := []string{
		TickTickTimeFormat,
		TickTickTimeFormat2,
		TickTickTimeFormat3,
		TickTickTimeFormat4,
		time.RFC3339,
		time.RFC3339Nano,
	}

	for _, format := range formats {
		if parsed, err := time.Parse(format, str); err == nil {
			t.Time = parsed
			return nil
		}
	}

	if strings.Contains(str, "+0000") {
		str = strings.Replace(str, "+0000", "Z", 1)
		if parsed, err := time.Parse(time.RFC3339Nano, str); err == nil {
			t.Time = parsed
			return nil
		}
		if parsed, err := time.Parse(time.RFC3339, str); err == nil {
			t.Time = parsed
			return nil
		}
	}

	return &time.ParseError{
		Layout: "TickTick time format",
		Value:  str,
	}
}

func (t TickTickTime) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + t.Time.Format(TickTickTimeFormat3) + `"`), nil
}

func (t TickTickTime) String() string {
	if t.Time.IsZero() {
		return ""
	}
	return t.Time.Format("2006-01-02 15:04")
}
