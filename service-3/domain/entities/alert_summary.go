package entities

import (
	"time"
)

type AlertSummary struct {
	TotalAlerts    int64     `json:"total_alerts"`
	ResolvedAlerts int64     `json:"resolved_alerts"`
	ActiveAlerts   int64     `json:"active_alerts"`
	LastAlertDate  time.Time `json:"last_alert_date"`
}
