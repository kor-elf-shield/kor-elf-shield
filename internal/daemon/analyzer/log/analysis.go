package log

import (
	analysisServices "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/log/analysis"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/log/analysis/alert_group"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/notifications"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

type Analysis interface {
	Alert(entry *analysisServices.Entry)
}

type analysis struct {
	alertService analysisServices.Alert
}

func NewAnalysis(rulesIndex *analysisServices.RulesIndex, repositories db.Repositories, logger log.Logger, notify notifications.Notifications) Analysis {
	alertGroupService := alert_group.NewGroup(repositories.AlertGroup(), logger)

	return &analysis{
		alertService: analysisServices.NewAlert(rulesIndex, alertGroupService, logger, notify),
	}
}

func (a *analysis) Alert(entry *analysisServices.Entry) {
	a.alertService.Analyze(entry)
}
