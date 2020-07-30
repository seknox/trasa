package stats

import (
	"github.com/huandu/go-sqlbuilder"
	"time"
)

func addTimeFilter(timeFilter string, now time.Time, sb sqlbuilder.SelectBuilder) sqlbuilder.SelectBuilder {
	switch timeFilter {
	case "Today":
		todayStart := now.Add(-time.Hour*time.Duration(now.Hour()) - time.Minute*time.Duration(now.Minute()))
		sb.Where(sb.GreaterThan(`login_time`, todayStart.UnixNano()))
	case "Yesterday":
		todayStart := now.Add(-time.Hour*time.Duration(now.Hour()) - time.Minute*time.Duration(now.Minute()))
		yesterdayStart := now.Add(-time.Hour*time.Duration(now.Hour()) - time.Minute*time.Duration(now.Minute())).Add(-time.Hour * 24)
		sb.Where(sb.LessThan(`login_time`, todayStart.UnixNano()))
		sb.Where(sb.GreaterThan(`login_time`, yesterdayStart.UnixNano()))
	case "7 Days":
		start := now.Add(-time.Hour*time.Duration(now.Hour()) - time.Minute*time.Duration(now.Minute())).Add(-time.Hour * 24 * 7)
		sb.Where(sb.GreaterThan(`login_time`, start.UnixNano()))
	case "30 Days":
		start := now.Add(-time.Hour*time.Duration(now.Hour()) - time.Minute*time.Duration(now.Minute())).Add(-time.Hour * 24 * 30)
		sb.Where(sb.GreaterThan(`login_time`, start.UnixNano()))
	case "90 Days":
		start := now.Add(-time.Hour*time.Duration(now.Hour()) - time.Minute*time.Duration(now.Minute())).Add(-time.Hour * 24 * 90)
		sb.Where(sb.GreaterThan(`login_time`, start.UnixNano()))
	}

	return sb
}
