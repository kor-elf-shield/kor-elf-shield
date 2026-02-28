package brute_force_protection_group

import (
	"fmt"
	"net"
	"time"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config/brute_force_protection"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/entity"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/repository"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/pkg/time_operation"
)

type Group interface {
	Analyze(group *brute_force_protection.Group, eventTime time.Time, ip net.IP, message string) (AnalysisResult, error)
	ClearDBData() error
}

type group struct {
	groupRepository repository.BruteForceProtectionGroupRepository
	logger          log.Logger
}

type AnalysisResult struct {
	Block    bool
	BlockSec uint32
	LastLogs []string
}

func NewGroup(groupRepository repository.BruteForceProtectionGroupRepository, logger log.Logger) Group {
	return &group{
		groupRepository: groupRepository,
		logger:          logger,
	}
}

func (g *group) Analyze(group *brute_force_protection.Group, eventTime time.Time, ip net.IP, message string) (AnalysisResult, error) {
	analysisResult := AnalysisResult{
		Block: false,
	}

	g.logger.Debug(fmt.Sprintf("Analyzing brute force protection group %s IP %s", group.Name, ip.String()))

	err := g.groupRepository.Update(group.Name, ip, func(entityGroup *entity.BruteForceProtectionGroup) (*entity.BruteForceProtectionGroup, error) {
		rateLimit, err := group.RateLimit(entityGroup.CurrentLevelTriggerCount)
		if err != nil {
			return entityGroup, err
		}

		if time_operation.IsRateLimited(entityGroup.LastTriggeredAtUnix, eventTime, int64(rateLimit.Period)) {
			g.logger.Debug(fmt.Sprintf("Brute force protection group %s is rate limited", group.Name))
			analysisResult, entityGroup = g.analysisResult(rateLimit, eventTime, message, entityGroup)
			return entityGroup, nil
		}

		entityGroup.TriggerCount = 0

		if time_operation.IsReset(entityGroup.LastTriggeredAtUnix, eventTime, int64(group.RateLimitResetPeriod)) {
			g.logger.Debug(fmt.Sprintf("Brute force protection group %s is reset", group.Name))
			entityGroup.Reset()
			rateLimit, err = group.RateLimit(0)
			if err != nil {
				return entityGroup, err
			}
		}

		g.logger.Debug(fmt.Sprintf("Brute force protection not rate limited"))
		analysisResult, entityGroup = g.analysisResult(rateLimit, eventTime, message, entityGroup)

		return entityGroup, nil
	})

	if err != nil {
		return AnalysisResult{
			Block: false,
		}, err
	}

	return analysisResult, nil
}

func (g *group) ClearDBData() error {
	return g.groupRepository.Clear()
}

func (g *group) analysisResult(rateLimit brute_force_protection.RateLimit, eventTime time.Time, message string, entityGroup *entity.BruteForceProtectionGroup) (AnalysisResult, *entity.BruteForceProtectionGroup) {
	analysisResult := AnalysisResult{
		Block: false,
	}

	entityGroup.LastTriggeredAtUnix = eventTime.Unix()
	entityGroup.TriggerCount++
	entityGroup.LastLogs = append(entityGroup.LastLogs, fmt.Sprintf("event time: %s, message: %s", eventTime.Format(time.RFC3339), message))
	g.logger.Debug(fmt.Sprintf("Brute force protection triggered. Count: %d", entityGroup.TriggerCount))

	if entityGroup.TriggerCount >= uint64(rateLimit.Count) {
		g.logger.Debug(fmt.Sprintf("Brute force protection reached rate limit"))
		analysisResult.LastLogs = entityGroup.LastLogs
		analysisResult.Block = true
		analysisResult.BlockSec = rateLimit.BlockingTimeSeconds

		entityGroup.CurrentLevelTriggerCount++
		entityGroup.TriggerCount = 0
		entityGroup.LastLogs = []string{}
	} else {
		g.logger.Debug(fmt.Sprintf("Brute force protection not reached rate limit"))
	}

	return analysisResult, entityGroup
}
