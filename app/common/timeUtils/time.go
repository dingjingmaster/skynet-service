package timeUtils

import (
	"fmt"
	"strconv"
	"time"
)

func GetTimeStamp () int {
	tim := time.Now()

	ret, _ := strconv.Atoi(fmt.Sprintf("%04d%02d%02d%02d", tim.Year(), tim.Month(), tim.Day(), tim.Hour()))

	return ret
}