package alert_group

import (
	"fmt"
	"time"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/entity"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/repository"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/pkg/time_operation"
)

type Group interface {
	Analyze(alertGroup *config.AlertGroup, eventTime time.Time, message string, partition *string) (AnalysisResult, error)
	ClearDBData() error
}

type group struct {
	alertGroupRepository repository.AlertGroupRepository
	logger               log.Logger
}

type AnalysisResult struct {
	Alerted     bool
	AlertNumber uint64
	LastLogs    []string
}

func NewGroup(alertGroupRepository repository.AlertGroupRepository, logger log.Logger) Group {
	return &group{
		alertGroupRepository: alertGroupRepository,
		logger:               logger,
	}
}

func (g *group) Analyze(alertGroup *config.AlertGroup, eventTime time.Time, message string, partition *string) (AnalysisResult, error) {
	analysisResult := AnalysisResult{
		Alerted: false,
	}

	g.logger.Debug(fmt.Sprintf("Analyzing alert group %s", alertGroup.Name))

	err := g.alertGroupRepository.Update(alertGroup.Name, partition, func(entityAlertGroup *entity.AlertGroup) (*entity.AlertGroup, error) {
		rateLimit, err := alertGroup.RateLimit(entityAlertGroup.CurrentLevelTriggerCount)
		if err != nil {
			return entityAlertGroup, err
		}

		if time_operation.IsRateLimited(entityAlertGroup.LastTriggeredAtUnix, eventTime, int64(rateLimit.Period)) {
			g.logger.Debug(fmt.Sprintf("Alert group %s is rate limited", alertGroup.Name))
			analysisResult, entityAlertGroup = g.analysisResult(rateLimit, eventTime, message, entityAlertGroup)
			return entityAlertGroup, nil
		}

		entityAlertGroup.TriggerCount = 0

		if time_operation.IsReset(entityAlertGroup.LastTriggeredAtUnix, eventTime, int64(alertGroup.RateLimitResetPeriod)) {
			g.logger.Debug(fmt.Sprintf("Alert group %s is reset", alertGroup.Name))
			entityAlertGroup.Reset()
			rateLimit, err = alertGroup.RateLimit(0)
			if err != nil {
				return entityAlertGroup, err
			}
		}

		g.logger.Debug(fmt.Sprintf("Alert not rate limited"))
		entityAlertGroup.LastLogs = []string{}
		analysisResult, entityAlertGroup = g.analysisResult(rateLimit, eventTime, message, entityAlertGroup)

		return entityAlertGroup, nil
	})

	if err != nil {
		return AnalysisResult{
			Alerted: false,
		}, err
	}

	return analysisResult, nil
}

func (g *group) ClearDBData() error {
	return g.alertGroupRepository.Clear()
}

func (g *group) analysisResult(rateLimit config.RateLimit, eventTime time.Time, message string, entityAlertGroup *entity.AlertGroup) (AnalysisResult, *entity.AlertGroup) {
	analysisResult := AnalysisResult{
		Alerted: false,
	}

	entityAlertGroup.LastTriggeredAtUnix = eventTime.Unix()
	entityAlertGroup.TriggerCount++
	entityAlertGroup.LastLogs = append(entityAlertGroup.LastLogs, fmt.Sprintf("event time: %s, message: %s", eventTime.Format(time.RFC3339), message))
	g.logger.Debug(fmt.Sprintf("Alert triggered. Count: %d", entityAlertGroup.TriggerCount))

	if entityAlertGroup.TriggerCount >= uint64(rateLimit.Count) {
		g.logger.Debug(fmt.Sprintf("Alert reached rate limit"))
		analysisResult.LastLogs = entityAlertGroup.LastLogs
		analysisResult.Alerted = true

		entityAlertGroup.CurrentLevelTriggerCount++
		entityAlertGroup.TriggerCount = 0
		entityAlertGroup.LastLogs = []string{}

		analysisResult.AlertNumber = entityAlertGroup.CurrentLevelTriggerCount
	} else {
		g.logger.Debug(fmt.Sprintf("Alert not reached rate limit"))
	}

	return analysisResult, entityAlertGroup
}
