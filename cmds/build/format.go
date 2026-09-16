package build_cmd

import (
	"app/data"
	"fmt"
	"math"
	"time"
)

func (BuildCommand) formatDurationTimestamp(duration float32) string {
	d := int(math.Round(float64(duration)))
	return fmt.Sprintf("%02d:%02d:%02d", d/(60*60), (d/60)%60, d%60)
}

func (BuildCommand) formatDurationHMS(duration float32) string {
	d := int(math.Round(float64(duration)))
	return fmt.Sprintf("%dh %dm %ds", d/(60*60), (d/60)%60, d%60)
}

func (BuildCommand) formatDateRange(startDate, endDate time.Time) string {
	if startDate.IsZero() && endDate.IsZero() {
		return ""
	}
	return fmt.Sprintf("%s - %s", startDate.Format(data.DATE_FORMAT), endDate.Format(data.DATE_FORMAT))
}
