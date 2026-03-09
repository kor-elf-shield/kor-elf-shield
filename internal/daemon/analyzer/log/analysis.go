package log

import (
	"fmt"

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
	ClearDBData() ([]error, error)
}

type analysis struct {
	alertService                analysisServices.Alert
	bruteForceProtectionService analysisServices.BruteForceProtection
}

func NewAnalysis(rulesIndex *analysisServices.RulesIndex, blockService brute_force_protection_group.BlockService, repositories db.Repositories, logger log.Logger, notify notifications.Notifications) Analysis {
	alertGroupService := alert_group.NewGroup(repositories.AlertGroup(), logger)
	bruteForceProtectionGroupService := brute_force_protection_group.NewGroup(repositories.BruteForceProtectionGroup(), logger)

	return &analysis{
		alertService:                analysisServices.NewAlert(rulesIndex, alertGroupService, logger, notify),
		bruteForceProtectionService: analysisServices.NewBruteForceProtection(rulesIndex, bruteForceProtectionGroupService, blockService, logger, notify),
	}
}

func (a *analysis) Alert(entry *analysisServices.Entry) {
	a.alertService.Analyze(entry)
}

func (a *analysis) BruteForceProtection(entry *analysisServices.Entry) {
	a.bruteForceProtectionService.Analyze(entry)
}

func (a *analysis) ClearDBData() ([]error, error) {
	var errClearDB []error
	if err := a.alertService.ClearDBData(); err != nil {
		errClearDB = append(errClearDB, err)
	}
	if err := a.bruteForceProtectionService.ClearDBData(); err != nil {
		errClearDB = append(errClearDB, err)
	}

	if len(errClearDB) > 0 {
		return nil, fmt.Errorf("failed to clear database data: %v", errClearDB)
	}

	return errClearDB, nil
}
