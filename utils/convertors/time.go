package convertors

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ConvertDateStringToTime(input string) (time.Time, error) {
	parsedTime, err := time.Parse("2006-01-02", input)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format (%s): %v", input, err)
	}

	return parsedTime, nil
}

func ParseDurationToPGInterval(durationStr string) (pgtype.Interval, error) {
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return pgtype.Interval{}, fmt.Errorf("invalid duration format: %v", err)
	}

	return pgtype.Interval{
		Microseconds: duration.Microseconds(),
		Valid:        true,
	}, nil
}
