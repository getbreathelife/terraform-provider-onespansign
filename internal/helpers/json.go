package helpers

import (
	"encoding/json"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func GetJsonNumber(i int64, base int) json.Number {
	return json.Number(strconv.FormatInt(i, 10))
}

func RandJsonNumber(min int, max int) json.Number {
	return json.Number(strconv.Itoa(acctest.RandIntRange(min, max)))
}
