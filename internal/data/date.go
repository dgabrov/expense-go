package data

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

const dateLayout = "2006-01-02"

// Date is a calendar date with no time and no zone. It is sent as "YYYY-MM-DD"
// in JSON and stored in MySQL DATE columns, so no timezone conversion happens.
type Date struct {
	Year  int
	Month time.Month
	Day   int
}

// After reports whether d is later than other.
func (d *Date) After(other Date) bool {
	if d.Year != other.Year {
		return d.Year > other.Year
	}
	if d.Month != other.Month {
		return d.Month > other.Month
	}
	return d.Day > other.Day
}

func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Date) UnmarshalJSON(raw []byte) error {
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return fmt.Errorf("date must be a JSON string in YYYY-MM-DD format: %w", err)
	}

	return d.parse(s)
}

// Value implements driver.Valuer; the string form is accepted by MySQL DATE columns.
func (d Date) Value() (driver.Value, error) {
	return d.String(), nil
}

// Scan implements sql.Scanner for DATE columns read with parseTime enabled
// (time.Time at midnight UTC) or without it (bytes or string).
func (d *Date) Scan(src any) error {
	switch v := src.(type) {
	case time.Time:
		d.Year, d.Month, d.Day = v.Date()
		return nil
	case []byte:
		return d.parse(string(v))
	case string:
		return d.parse(v)
	default:
		return fmt.Errorf("cannot scan %T into Date", src)
	}
}

func (d *Date) parse(s string) error {
	parsed, err := time.Parse(dateLayout, s)
	if err != nil {
		return fmt.Errorf("invalid date %q, expected YYYY-MM-DD: %w", s, err)
	}

	d.Year, d.Month, d.Day = parsed.Year(), parsed.Month(), parsed.Day()
	return nil
}

func GetDateRange(startMonth, startYear, endMonth, endYear int) (*Date, *Date, error) {
	if !validMonthYear(startMonth, startYear) {
		return nil, nil, fmt.Errorf("invalid start month/year %d/%d", startMonth, startYear)
	}

	start := &Date{
		Year:  startYear,
		Month: time.Month(startMonth),
		Day:   1,
	}

	if endMonth == 0 && endYear == 0 {
		return start, nil, nil
	}
	if !validMonthYear(endMonth, endYear) {
		return nil, nil, fmt.Errorf("invalid end month/year %d/%d", endMonth, endYear)
	}

	end := &Date{
		Year:  endYear,
		Month: time.Month(endMonth) + 1,
		Day:   1,
	}
	if end.Month > time.December {
		end.Month = time.January
		end.Year++
	}

	return start, end, nil
}

func validMonthYear(month, year int) bool {
	return month >= int(time.January) && month <= int(time.December) && year >= 1 && year <= 9999
}
