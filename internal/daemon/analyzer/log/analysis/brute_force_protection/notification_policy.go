package brute_force_protection

import (
	"time"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config/brute_force_protection"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/entity"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/repository"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

type NotificationPolicy interface {
	IsNotify(rule *brute_force_protection.Rule) bool
}

type notificationPolicy struct {
	notifyPolicyRepository repository.BruteForceProtectionNotifyPolicyRepository
	logger                 log.Logger
}

func NewNotificationPolicy(notifyPolicyRepository repository.BruteForceProtectionNotifyPolicyRepository, logger log.Logger) NotificationPolicy {
	return &notificationPolicy{
		notifyPolicyRepository: notifyPolicyRepository,
		logger:                 logger,
	}
}

func (s *notificationPolicy) IsNotify(rule *brute_force_protection.Rule) bool {
	if !rule.IsNotification {
		return false
	}

	if rule.NotificationCooldown == 0 && rule.NotificationEvery == 0 {
		return true
	}

	isNotify := false

	err := s.notifyPolicyRepository.Update(rule.Name, func(notify *entity.BruteForceProtectionNotifyPolicy) (*entity.BruteForceProtectionNotifyPolicy, error) {
		if isEvery(rule, notify) {
			isNotify = true
			return resetNotifyPolicy(rule, notify), nil
		}

		if time.Now().Unix() >= notify.CooldownTime {
			isNotify = true
			return resetNotifyPolicy(rule, notify), nil
		}

		return notify, nil
	})

	if err != nil {
		s.logger.Error(err.Error())
		return true
	}

	return isNotify
}

func isEvery(rule *brute_force_protection.Rule, entity *entity.BruteForceProtectionNotifyPolicy) bool {
	if rule.NotificationEvery == 0 {
		return false
	}

	if entity.Every > 0 {
		entity.Every--
	}

	return entity.Every == 0
}

func resetNotifyPolicy(rule *brute_force_protection.Rule, entity *entity.BruteForceProtectionNotifyPolicy) *entity.BruteForceProtectionNotifyPolicy {
	entity.Every = rule.NotificationEvery
	cooldownTime := time.Now().Add(time.Duration(int(rule.NotificationCooldown)) * time.Second)
	entity.CooldownTime = cooldownTime.Unix()
	return entity
}
