package models

// LogEventStat is a single slice for pie charts (action frequency).
type LogEventStat struct {
	Event string `json:"event" example:"user.login"`
	Count int64  `json:"count" example:"42"`
}

// LogDailyStat is a single bar for bar charts (activity per day).
type LogDailyStat struct {
	Date  string `json:"date" example:"2026-05-28"`
	Count int64  `json:"count" example:"12"`
}
