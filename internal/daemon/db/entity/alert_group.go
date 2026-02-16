package entity

type AlertGroup struct {
	LastTriggeredAtUnix      int64    `json:"LastTriggeredAtUnix"`
	TriggerCount             uint64   `json:"TriggerCount"`
	CurrentLevelTriggerCount uint64   `json:"CurrentLevelTriggerCount"`
	LastLogs                 []string `json:"LastLogs"`
}

func (g *AlertGroup) Reset() {
	g.LastTriggeredAtUnix = 0
	g.TriggerCount = 0
	g.CurrentLevelTriggerCount = 0
	g.LastLogs = []string{}
}
