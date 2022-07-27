package helpers

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"time"
)

func GetJsonNumber(i int64, base int) json.Number {
	return json.Number(strconv.FormatInt(i, 10))
}

func RandJsonNumber(min int, max int) json.Number {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(max-min+1) + min
	return json.Number(strconv.Itoa(i))
}

func GetInt(n json.Number) (int, error) {
	i, err := n.Int64()

	if err != nil {
		return 0, err
	}

	return int(i), nil
}
