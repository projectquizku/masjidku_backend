package helper

import (
	"fmt"
	"strconv"
)

func ParseCoordinates(latStr, lngStr string) (float64, float64, error) {
	lat, err1 := strconv.ParseFloat(latStr, 64)
	lng, err2 := strconv.ParseFloat(lngStr, 64)

	if err1 != nil || err2 != nil {
		return 0, 0, fmt.Errorf("koordinat tidak valid")
	}
	return lat, lng, nil
}
