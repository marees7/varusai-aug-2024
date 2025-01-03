package helper

import (
	"shopping-site/pkg/loggers"
	"shopping-site/utils/constants"
	"strconv"

	"github.com/google/uuid"
)

// parse the string to uuid
func PasreUuid(id string) (uuid.UUID, error) {
	ParsedId, err := uuid.Parse(id)
	if err != nil {
		loggers.ErrorLog.Println(err)
		return uuid.Nil, err
	}

	return ParsedId, nil
}

// convert the limit and offset string to int and calculate offset using it
func Pagination(limitStr string, offsetStr string) (int, int, error) {
	offset, err := strconv.Atoi(offsetStr)
	if offsetStr == "" {
		offset = constants.OffsetDefault
	} else if err != nil {
		return 0, 0, err
	}

	limit, err := strconv.Atoi(limitStr)
	if limitStr == "" {
		limit = constants.LimitDefault
	} else if err != nil {
		return 0, 0, err
	}

	offset = (offset - 1) * limit

	return limit, offset, nil
}

func ToFloat(input string) (float64, error) {
	result, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0, err
	}

	return result, nil
}
