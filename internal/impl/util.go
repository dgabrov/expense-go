package impl

import (
	"fmt"
	"time"

	"github.com/dgabrov/expense-go/internal/data"
)

func getDateRange(startMonth, startYear, endMonth, endYear int) (*data.Date, *data.Date, error) {
	if !validMonthYear(startMonth, startYear) {
		return nil, nil, fmt.Errorf("invalid start month/year %d/%d", startMonth, startYear)
	}

	start := &data.Date{
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

	end := &data.Date{
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
