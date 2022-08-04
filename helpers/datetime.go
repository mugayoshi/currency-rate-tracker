package helpers

import (
	"fmt"
	"time"
)

func GetYyyyMmDd(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())
}
