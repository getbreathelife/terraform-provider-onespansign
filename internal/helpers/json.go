package helpers

import (
	"encoding/json"
	"strconv"
)

func GetJsonNumber(i int64, base int) json.Number {
	return json.Number(strconv.FormatInt(i, 10))
}
