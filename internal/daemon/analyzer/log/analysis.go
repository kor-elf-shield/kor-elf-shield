package log

import (
	analysisServices "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/log/analysis"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/notifications"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

type Analysis interface {
	Alert(entry *analysisServices.Entry)
}

type analysis struct {
	alertService analysisServices.Alert
}

func NewAnalysis(alertRuleIndex analysisServices.AlertRuleIndex, logger log.Logger, notify notifications.Notifications) Analysis {
	return &analysis{
		alertService: analysisServices.NewAlert(alertRuleIndex, logger, notify),
	}
}

func (a *analysis) Alert(entry *analysisServices.Entry) {
	a.alertService.Analyze(entry)
}
