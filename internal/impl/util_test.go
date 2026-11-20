package impl

import (
	"testing"
	"time"

	"github.com/dgabrov/expense-go/internal/data"
)

func TestGetDateRange(t *testing.T) {
	tests := []struct {
		name       string
		startMonth int
		startYear  int
		endMonth   int
		endYear    int
		wantStart  *data.Date
		wantEnd    *data.Date
		wantError  bool
	}{
		{
			name:       "start only",
			startMonth: 1,
			startYear:  2025,
			wantStart:  date(2025, time.January),
		},
		{
			name:       "minimum month",
			startMonth: 1,
			startYear:  1,
			endMonth:   1,
			endYear:    1,
			wantStart:  date(1, time.January),
			wantEnd:    date(1, time.February),
		},
		{
			name:       "same start and end month",
			startMonth: 5,
			startYear:  2024,
			endMonth:   5,
			endYear:    2024,
			wantStart:  date(2024, time.May),
			wantEnd:    date(2024, time.June),
		},
		{
			name:       "inclusive and exclusive dates",
			startMonth: 3,
			startYear:  2024,
			endMonth:   6,
			endYear:    2024,
			wantStart:  date(2024, time.March),
			wantEnd:    date(2024, time.July),
		},
		{
			name:       "end rolls into next year",
			startMonth: 11,
			startYear:  2024,
			endMonth:   12,
			endYear:    2024,
			wantStart:  date(2024, time.November),
			wantEnd:    date(2025, time.January),
		},
		{
			name:       "maximum end year rolls forward",
			startMonth: 1,
			startYear:  9999,
			endMonth:   12,
			endYear:    9999,
			wantStart:  date(9999, time.January),
			wantEnd:    date(10000, time.January),
		},
		{
			name:       "missing end month and year",
			startMonth: 12,
			startYear:  2025,
			endMonth:   0,
			endYear:    0,
			wantStart:  date(2025, time.December),
		},
		{
			name:       "start month below range",
			startMonth: 0,
			startYear:  2025,
			wantError:  true,
		},
		{
			name:       "negative start month",
			startMonth: -1,
			startYear:  2025,
			wantError:  true,
		},
		{
			name:       "start month above range",
			startMonth: 13,
			startYear:  2025,
			wantError:  true,
		},
		{
			name:       "missing start year",
			startMonth: 1,
			startYear:  0,
			wantError:  true,
		},
		{
			name:       "negative start year",
			startMonth: 1,
			startYear:  -1,
			wantError:  true,
		},
		{
			name:       "start year above range",
			startMonth: 1,
			startYear:  10000,
			wantError:  true,
		},
		{
			name:       "end month without year",
			startMonth: 1,
			startYear:  2025,
			endMonth:   2,
			wantError:  true,
		},
		{
			name:       "end year without month",
			startMonth: 1,
			startYear:  2025,
			endYear:    2025,
			wantError:  true,
		},
		{
			name:       "end month below range",
			startMonth: 1,
			startYear:  2025,
			endMonth:   0,
			endYear:    2025,
			wantError:  true,
		},
		{
			name:       "negative end month",
			startMonth: 1,
			startYear:  2025,
			endMonth:   -1,
			endYear:    2025,
			wantError:  true,
		},
		{
			name:       "end month above range",
			startMonth: 1,
			startYear:  2025,
			endMonth:   13,
			endYear:    2025,
			wantError:  true,
		},
		{
			name:       "missing end year",
			startMonth: 1,
			startYear:  2025,
			endMonth:   1,
			endYear:    0,
			wantError:  true,
		},
		{
			name:       "negative end year",
			startMonth: 1,
			startYear:  2025,
			endMonth:   1,
			endYear:    -1,
			wantError:  true,
		},
		{
			name:       "end year above range",
			startMonth: 1,
			startYear:  2025,
			endMonth:   1,
			endYear:    10000,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStart, gotEnd, err := getDateRange(tt.startMonth, tt.startYear, tt.endMonth, tt.endYear)
			if (err != nil) != tt.wantError {
				t.Fatalf("getDateRange() error = %v, wantError %t", err, tt.wantError)
			}
			if tt.wantError {
				if gotStart != nil || gotEnd != nil {
					t.Fatalf("getDateRange() returned dates with invalid input: start=%v, end=%v", gotStart, gotEnd)
				}
				return
			}

			if gotStart == nil || *gotStart != *tt.wantStart {
				t.Errorf("getDateRange() start = %v, want %v", gotStart, tt.wantStart)
			}
			if tt.wantEnd == nil {
				if gotEnd != nil {
					t.Errorf("getDateRange() end = %v, want nil", gotEnd)
				}
			} else if gotEnd == nil || *gotEnd != *tt.wantEnd {
				t.Errorf("getDateRange() end = %v, want %v", gotEnd, tt.wantEnd)
			}
		})
	}
}

func date(year int, month time.Month) *data.Date {
	return &data.Date{Year: year, Month: month, Day: 1}
}
