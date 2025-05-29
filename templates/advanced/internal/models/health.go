package models

import "time"

// HealthReport represents the overall health status of the service
type HealthReport struct {
	Status      string        `json:"status"`
	Timestamp   time.Time     `json:"timestamp"`
	Version     string        `json:"version"`
	Uptime      time.Duration `json:"uptime"`
	UptimeHuman string        `json:"uptime_human"`
}

// ServerTime represents server time information with multiple formats
type ServerTime struct {
	Timestamp time.Time `json:"timestamp"`
	Timezone  string    `json:"timezone"`
	Unix      int64     `json:"unix"`
	UnixMilli int64     `json:"unix_milli"`
	ISO8601   string    `json:"iso8601"`
	Formatted string    `json:"formatted"`
}
