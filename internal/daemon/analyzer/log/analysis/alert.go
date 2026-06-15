package analysis

import (
	"fmt"
	"time"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/log/analysis/alert_group"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/geoip"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/notifications"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/i18n"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

type Alert interface {
	Analyze(entry *Entry)
	ClearDBData() error
}

type alert struct {
	rulesIndex        *RulesIndex
	alertGroupService alert_group.Group
	logger            log.Logger
	notify            notifications.Notifications
	ipInfo            geoip.Info
}

type alertAnalyzeRuleReturn struct {
	found  bool
	fields []*regexField
}

type alertNotify struct {
	rule        *config.AlertRule
	messages    []string
	alertNumber uint64
	time        time.Time
	fields      []*regexField
}

func NewAlert(
	rulesIndex *RulesIndex,
	alertGroupService alert_group.Group,
	logger log.Logger,
	notify notifications.Notifications,
	ipInfo geoip.Info,
) Alert {
	return &alert{
		rulesIndex:        rulesIndex,
		alertGroupService: alertGroupService,
		logger:            logger,
		notify:            notify,
		ipInfo:            ipInfo,
	}
}

func (a *alert) Analyze(entry *Entry) {
	rules, err := a.rulesIndex.Alerts(entry)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Failed to get alert rules: %s", err))
	}
	for _, rule := range rules {
		result := a.analyzeRule(rule, entry.Message)
		if !result.found {
			continue
		}
		groupName := ""
		alertNumber := uint64(0)
		messages := []string{}
		if rule.Group != nil {
			alertGroup, err := a.alertGroupService.Analyze(rule.Group, entry.Time, entry.Message)
			if err != nil {
				a.logger.Error(fmt.Sprintf("Failed to analyze alert group: %s", err))
				continue
			}
			if !alertGroup.Alerted {
				continue
			}

			groupName = rule.Group.Name
			for _, lastLog := range alertGroup.LastLogs {
				messages = append(messages, lastLog)
			}
			alertNumber = alertGroup.AlertNumber
		} else {
			messages = append(messages, entry.Message)
		}
		a.logger.Info(fmt.Sprintf("Alert detected (%s) (group:%s): %s", rule.Name, groupName, entry.Message))
		a.sendNotify(&alertNotify{
			rule:        rule,
			messages:    messages,
			alertNumber: alertNumber,
			time:        entry.Time,
			fields:      result.fields,
		})
	}
}

func (a *alert) ClearDBData() error {
	return a.alertGroupService.ClearDBData()
}

func (a *alert) analyzeRule(rule *config.AlertRule, message string) alertAnalyzeRuleReturn {
	result := alertAnalyzeRuleReturn{
		found:  false,
		fields: []*regexField{},
	}

	for _, pattern := range rule.Patterns {
		re, err := pattern.Regexp.Get()
		if err != nil {
			a.logger.Error(fmt.Sprintf("Failed to compile regexp: %s", err))
			continue
		}

		idx := re.FindStringSubmatchIndex(message)

		if idx != nil {
			for _, value := range pattern.Values {
				start, end, err := getValueStartEndByRegexIndex(int(value.Value), idx)
				if err != nil {
					result.fields = append(result.fields, &regexField{name: value.Name, value: i18n.Lang.T("unknown")})
					continue
				}
				result.fields = append(result.fields, &regexField{name: value.Name, value: message[start:end], typeValue: value.Type})
			}

			if len(pattern.Values) != len(result.fields) {
				continue
			}

			result.found = true
			return result
		}
	}

	return result
}

func (a *alert) sendNotify(notify *alertNotify) {
	if !notify.rule.IsNotification {
		return
	}

	groupName := ""
	groupMessage := ""
	if notify.rule.Group != nil {
		groupName = notify.rule.Group.Name
		groupMessage = notify.rule.Group.Message + "\n\n"
	}

	subject := i18n.Lang.T("alert.subject", map[string]any{
		"Name":      notify.rule.Name,
		"GroupName": groupName,
	})
	text := subject + "\n\n" + groupMessage + notify.rule.Message + "\n\n"
	text += i18n.Lang.T("time", map[string]any{
		"Time": notify.time,
	}) + "\n"

	for _, field := range notify.fields {
		v := field.value
		if field.typeValue == config.PatternValueIP {
			if ipInfo, err := a.ipInfo(field.value); err != nil {
				a.logger.Error(fmt.Sprintf("Failed to get geoip info for ip %s: %s", v, err))
			} else {
				v = ipInfo
			}
		}
		text += fmt.Sprintf("%s: %s\n", field.name, v)
	}
	if notify.alertNumber > 0 {
		text += i18n.Lang.T("alertNumber", map[string]any{
			"Count": notify.alertNumber,
		}) + "\n"
	}
	text += "\n" + i18n.Lang.T("log", map[string]any{
		"Count": len(notify.messages),
	}) + "\n"
	for _, message := range notify.messages {
		text += message + "\n\n"
	}
	a.notify.SendAsync(notifications.Message{Subject: subject, Body: text})
}
