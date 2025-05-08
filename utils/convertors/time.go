package convertors

import (
	"fmt"
	"time"
)

func GetReleaseDateAsTime(strDate string) (time.Time, error) {
	parsedTime, err := time.Parse("2006-01-02", strDate)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format: %v", err)
	}

	return parsedTime, nil
}
