package log

import (
	analysisServices "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/log/analysis"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/log/analysis/alert_group"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/log/analysis/brute_force_protection_group"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/notifications"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

type Analysis interface {
	Alert(entry *analysisServices.Entry)
	BruteForceProtection(entry *analysisServices.Entry)
}

type analysis struct {
	alertService                analysisServices.Alert
	bruteForceProtectionService analysisServices.BruteForceProtection
}

func NewAnalysis(rulesIndex *analysisServices.RulesIndex, blockIPFunc analysisServices.BlockIPFunc, repositories db.Repositories, logger log.Logger, notify notifications.Notifications) Analysis {
	alertGroupService := alert_group.NewGroup(repositories.AlertGroup(), logger)
	bruteForceProtectionGroupService := brute_force_protection_group.NewGroup(repositories.BruteForceProtectionGroup(), logger)

	return &analysis{
		alertService:                analysisServices.NewAlert(rulesIndex, alertGroupService, logger, notify),
		bruteForceProtectionService: analysisServices.NewBruteForceProtection(rulesIndex, bruteForceProtectionGroupService, blockIPFunc, logger, notify),
	}
}

func (a *analysis) Alert(entry *analysisServices.Entry) {
	a.alertService.Analyze(entry)
}

func (a *analysis) BruteForceProtection(entry *analysisServices.Entry) {
	a.bruteForceProtectionService.Analyze(entry)
}
