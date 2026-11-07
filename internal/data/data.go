package data

import (
	"errors"
)

type ConfigData struct {
	Version         string         `json:"version"`
	Context         string         `json:"context"`
	Auth            string         `json:"auth"`
	TokenTTLSeconds int            `json:"tokenTTLSeconds"`
	Database        DatabaseConfig `json:"database"`
	Server          ServerConfig   `json:"server"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type ServerConfig struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	LogRequests bool   `json:"logRequests"`
}

type IdsHolder struct {
	Ids []string `json:"ids"`
}

type FilterData struct {
	Month int
	Year  int
}

func (v FilterData) validate() error {
	if v.Month < 1 || v.Month > 12 {
		return errors.New("invalid month")
	}

	if v.Year <= 0 {
		return errors.New("invalid year")
	}

	return nil
}

type ReportRequest struct {
	Start       *FilterData
	End         *FilterData
	ExcludeNone bool
}

// ValidateStart checks only the start date; End and ExcludeNone are not considered.
func (r ReportRequest) ValidateStart() error {
	if r.Start == nil {
		return errors.New("please fill out the start date")
	}

	return r.Start.validate()
}

// Validate checks Start, and End when it is given, and rejects an End that is before Start.
func (r ReportRequest) Validate() error {
	// start must be filled out and valid
	if r.Start == nil {
		return errors.New("please fill out the start date")
	}

	err := r.Start.validate()
	if err != nil {
		return err
	}

	// end if filled out, must be valid
	if r.End != nil {
		err = r.End.validate()

		if err != nil {
			return err
		}
		// then if both of them are filled out and valid, get the date range and check they are valid as well
		date1, date2, err := GetDateRange(r.Start.Month, r.Start.Year, r.End.Month, r.End.Year)
		if err != nil {
			return err
		}

		if date2 != nil && !date2.After(*date1) {
			return errors.New("start date is after end date")
		}
	}

	return nil
}

type ReportLastResponse struct {
	AccountID string `json:"accountId"`
	LastDate  Date   `json:"lastDate"`
}

// ReportGroupResponseItem is one month's total for one group. GroupID is empty for
// transactions that have no category or group; it is still sent so they are visible.
type ReportGroupResponseItem struct {
	Month   int     `json:"month"`
	Year    int     `json:"year"`
	GroupID string  `json:"groupId"`
	Amount  float64 `json:"amount"`
}

type ReportGroupResponse struct {
	Items []ReportGroupResponseItem `json:"items"`
}

type ReportCategoryResponseItem struct {
	Month      int     `json:"month"`
	Year       int     `json:"year"`
	CategoryID string  `json:"categoryId"`
	Amount     float64 `json:"amount"`
}

type ReportCategoryResponse struct {
	Items []ReportCategoryResponseItem `json:"items"`
}

type ReportListResponse struct {
	Items []Transaction `json:"items"`
}
