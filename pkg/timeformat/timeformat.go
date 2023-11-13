package timeformat

import (
	"time"

	"github.com/andanhm/go-prettytime"
)

func FormatTime(t time.Time, showExactTime bool) string {
	if !showExactTime {
		return prettytime.Format(t)
	}
	return t.String()
}
