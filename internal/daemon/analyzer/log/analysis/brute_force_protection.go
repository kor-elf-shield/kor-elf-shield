package analysis

import (
	"fmt"
	"net"
	"time"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config/brute_force_protection"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/log/analysis/brute_force_protection_group"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/blocking"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/notifications"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/i18n"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

type BruteForceProtection interface {
	Analyze(entry *Entry)
	ClearDBData() error
}

type BlockIPFunc func(blockIP blocking.BlockIP) error

type bruteForceProtection struct {
	rulesIndex   *RulesIndex
	groupService brute_force_protection_group.Group
	blockIP      BlockIPFunc
	logger       log.Logger
	notify       notifications.Notifications
}

type bruteForceProtectionAnalyzeRuleReturn struct {
	found  bool
	fields []*regexField
	ip     net.IP
}

type bruteForceProtectionNotify struct {
	rule     *brute_force_protection.Rule
	messages []string
	ip       net.IP
	time     time.Time
	fields   []*regexField
	blockSec uint32
}

func NewBruteForceProtection(rulesIndex *RulesIndex, groupService brute_force_protection_group.Group, blockIP BlockIPFunc, logger log.Logger, notify notifications.Notifications) BruteForceProtection {
	return &bruteForceProtection{
		rulesIndex:   rulesIndex,
		groupService: groupService,
		blockIP:      blockIP,
		logger:       logger,
		notify:       notify,
	}
}

func (p *bruteForceProtection) Analyze(entry *Entry) {
	rules, err := p.rulesIndex.BruteForceProtections(entry)
	if err != nil {
		p.logger.Error(fmt.Sprintf("Failed to get brute force protection rules for entry: %v", err))
		return
	}
	for _, rule := range rules {
		if rule.Group == nil {
			p.logger.Error("Brute force protection rule without group")
			continue
		}

		result := p.analyzeRule(rule, entry.Message)
		if !result.found {
			continue
		}

		groupResult, err := p.groupService.Analyze(rule.Group, entry.Time, result.ip, entry.Message)
		if err != nil {
			p.logger.Error(fmt.Sprintf("Failed to analyze brute force protection group: %s", err))
			continue
		}

		if !groupResult.Block {
			continue
		}

		blockIP := blocking.BlockIP{
			IP:          result.ip,
			TimeSeconds: groupResult.BlockSec,
			Reason:      rule.Message,
		}
		if err := p.blockIP(blockIP); err != nil {
			p.logger.Info(fmt.Sprintf("IP %s are not blocked (%s) (group:%s): %s. Err: %s", result.ip, rule.Name, rule.Group.Name, entry.Message, err.Error()))
			p.sendNotifyError(&bruteForceProtectionNotify{
				rule:     rule,
				ip:       result.ip,
				messages: groupResult.LastLogs,
				time:     entry.Time,
				fields:   result.fields,
				blockSec: groupResult.BlockSec,
			}, err)
			continue
		}

		p.logger.Info(fmt.Sprintf("Block IP %s detected (%s) (group:%s): %s", result.ip, rule.Name, rule.Group.Name, entry.Message))
		p.sendNotify(&bruteForceProtectionNotify{
			rule:     rule,
			ip:       result.ip,
			messages: groupResult.LastLogs,
			time:     entry.Time,
			fields:   result.fields,
			blockSec: groupResult.BlockSec,
		})
	}
}

func (p *bruteForceProtection) ClearDBData() error {
	return p.groupService.ClearDBData()
}

func (p *bruteForceProtection) analyzeRule(rule *brute_force_protection.Rule, message string) bruteForceProtectionAnalyzeRuleReturn {
	result := bruteForceProtectionAnalyzeRuleReturn{
		found:  false,
		fields: []*regexField{},
		ip:     nil,
	}

	for _, pattern := range rule.Patterns {
		re, err := pattern.Regexp.Get()
		if err != nil {
			p.logger.Error(fmt.Sprintf("Failed to compile regexp: %s", err))
			continue
		}

		idx := re.FindStringSubmatchIndex(message)

		if idx != nil {
			start, end, err := getValueStartEndByRegexIndex(int(pattern.IP), idx)
			if err != nil {
				p.logger.Error(fmt.Sprintf("Failed to get ip value: %s", err))
				return result
			}
			ipText := message[start:end]
			result.ip = net.ParseIP(ipText)
			if result.ip == nil {
				p.logger.Error(fmt.Sprintf("Failed to parse ip: %s", ipText))
				return bruteForceProtectionAnalyzeRuleReturn{
					found: false,
				}
			}

			for _, value := range pattern.Values {
				start, end, err := getValueStartEndByRegexIndex(int(value.Value), idx)
				if err != nil {
					result.fields = append(result.fields, &regexField{name: value.Name, value: i18n.Lang.T("unknown")})
					continue
				}
				result.fields = append(result.fields, &regexField{name: value.Name, value: message[start:end]})
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

func (p *bruteForceProtection) sendNotify(notify *bruteForceProtectionNotify) {
	if !notify.rule.IsNotification {
		return
	}

	groupName := notify.rule.Group.Name
	groupMessage := notify.rule.Group.Message + "\n\n"

	subject := i18n.Lang.T("alert.bruteForceProtection.subject", map[string]any{
		"Name":      notify.rule.Name,
		"GroupName": groupName,
		"IP":        notify.ip,
	})
	text := subject + "\n\n" + groupMessage + notify.rule.Message + "\n\n"
	text += "IP: " + notify.ip.String() + "\n"
	text += i18n.Lang.T("time", map[string]any{
		"Time": notify.time,
	}) + "\n"
	for _, field := range notify.fields {
		text += fmt.Sprintf("%s: %s\n", field.name, field.value)
	}
	text += "\n" + i18n.Lang.T("log") + "\n"
	for _, message := range notify.messages {
		text += message + "\n"
	}
	p.notify.SendAsync(notifications.Message{Subject: subject, Body: text})
}

func (p *bruteForceProtection) sendNotifyError(notify *bruteForceProtectionNotify, err error) {
	if !notify.rule.IsNotification {
		return
	}

	groupName := notify.rule.Group.Name
	groupMessage := notify.rule.Group.Message + "\n\n"

	subject := i18n.Lang.T("alert.bruteForceProtection.subject-error", map[string]any{
		"Name":      notify.rule.Name,
		"GroupName": groupName,
		"IP":        notify.ip,
	})
	text := subject + "\n\n" + groupMessage + notify.rule.Message + "\n\n"
	text += i18n.Lang.T("alert.bruteForceProtection.error", map[string]any{
		"Error": err.Error(),
	}) + "\n"
	text += "IP: " + notify.ip.String() + "\n"
	text += i18n.Lang.T("blockSec", map[string]any{
		"BlockSec": notify.blockSec,
	}) + "\n"
	text += i18n.Lang.T("time", map[string]any{
		"Time": notify.time,
	}) + "\n"
	for _, field := range notify.fields {
		text += fmt.Sprintf("%s: %s\n", field.name, field.value)
	}
	text += "\n" + i18n.Lang.T("log") + "\n"
	for _, message := range notify.messages {
		text += message + "\n"
	}
	p.notify.SendAsync(notifications.Message{Subject: subject, Body: text})
}
