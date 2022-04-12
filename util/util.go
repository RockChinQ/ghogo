package util

import (
	"fmt"
	"strconv"
	"time"
)

func GetTimeStr(t time.Time) string {
	return strconv.Itoa(t.Year()) + "-" +
		format(int(t.Month())) + "-" +
		format(t.Day()) + "," +
		format(t.Hour()) + ":" +
		format(t.Minute()) + ":" +
		format(t.Second())
}

func format(i int) string {
	return fmt.Sprintf("%02d", i)
}
