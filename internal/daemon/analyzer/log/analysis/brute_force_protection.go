package analysis

import (
	"fmt"
	"net"
	"strings"
	"time"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config/brute_force_protection"
	analysisBruteForceProtection "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/log/analysis/brute_force_protection"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/log/analysis/brute_force_protection_group"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/blocking"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/types"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/geoip"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/notifications"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/i18n"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/pkg/format"
)

type BruteForceProtection interface {
	Analyze(entry *Entry)
	ClearDBData() error
}

type bruteForceProtection struct {
	rulesIndex         *RulesIndex
	groupService       brute_force_protection_group.Group
	blockService       brute_force_protection_group.BlockService
	notificationPolicy analysisBruteForceProtection.NotificationPolicy
	logger             log.Logger
	notify             notifications.Notifications
	ipInfo             geoip.Info
}

type bruteForceProtectionAnalyzeRuleReturn struct {
	found     bool
	fields    []*regexField
	ip        net.IP
	partition *string
}

type bruteForceProtectionNotify struct {
	rule         *brute_force_protection.Rule
	messages     []string
	blockIPCount uint64
	ip           net.IP
	ports        []types.L4Port
	time         time.Time
	fields       []*regexField
	blockSec     uint32
	err          error
}

func NewBruteForceProtection(
	rulesIndex *RulesIndex,
	groupService brute_force_protection_group.Group,
	blockService brute_force_protection_group.BlockService,
	notificationPolicy analysisBruteForceProtection.NotificationPolicy,
	logger log.Logger,
	notify notifications.Notifications,
	ipInfo geoip.Info,
) BruteForceProtection {
	return &bruteForceProtection{
		rulesIndex:         rulesIndex,
		groupService:       groupService,
		blockService:       blockService,
		notificationPolicy: notificationPolicy,
		logger:             logger,
		notify:             notify,
		ipInfo:             ipInfo,
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

		groupResult, err := p.groupService.Analyze(rule.Group, entry.Time, result.ip, entry.Message, result.partition)
		if err != nil {
			p.logger.Error(fmt.Sprintf("Failed to analyze brute force protection group: %s", err))
			continue
		}

		if !groupResult.Block {
			continue
		}

		ipWithPorts, l4Ports := groupResult.BlockConfig.PortsBlocked()
		if !ipWithPorts {
			p.handleBlockIP(entry, rule, &result, &groupResult)
			continue
		}

		p.handleBlockIPWithPorts(entry, rule, &result, &groupResult, l4Ports)
	}
}

func (p *bruteForceProtection) ClearDBData() error {
	return p.groupService.ClearDBData()
}

func (p *bruteForceProtection) handleBlockIP(
	entry *Entry,
	rule *brute_force_protection.Rule,
	result *bruteForceProtectionAnalyzeRuleReturn,
	groupResult *brute_force_protection_group.AnalysisResult,
) {
	blockIP := blocking.BlockIP{
		IP:          result.ip,
		TimeSeconds: groupResult.BlockSec,
		Reason:      rule.Message,
	}
	isBanned, err := p.blockService.BlockIP(blockIP)
	if isBanned == false {
		p.logger.Info(fmt.Sprintf("IP %s are not blocked (%s) (group:%s): %s. Err: %s", result.ip, rule.Name, rule.Group.Name, entry.Message, err.Error()))
		p.sendNotifyError(&bruteForceProtectionNotify{
			rule:         rule,
			ip:           result.ip,
			messages:     groupResult.LastLogs,
			blockIPCount: groupResult.BlockIPCount,
			time:         entry.Time,
			fields:       result.fields,
			blockSec:     groupResult.BlockSec,
			err:          err,
		})
		return
	}

	p.logger.Info(fmt.Sprintf("Block IP %s detected (%s) (group:%s): %s", result.ip, rule.Name, rule.Group.Name, entry.Message))
	p.sendNotifySuccess(&bruteForceProtectionNotify{
		rule:         rule,
		ip:           result.ip,
		messages:     groupResult.LastLogs,
		blockIPCount: groupResult.BlockIPCount,
		time:         entry.Time,
		fields:       result.fields,
		blockSec:     groupResult.BlockSec,
		err:          err,
	})
}

func (p *bruteForceProtection) handleBlockIPWithPorts(
	entry *Entry,
	rule *brute_force_protection.Rule,
	result *bruteForceProtectionAnalyzeRuleReturn,
	groupResult *brute_force_protection_group.AnalysisResult,
	l4Ports []types.L4Port,
) {
	blockIPWithPorts := blocking.BlockIPWithPorts{
		IP:          result.ip,
		TimeSeconds: groupResult.BlockSec,
		Reason:      rule.Message,
		Ports:       l4Ports,
	}
	isBanned, err := p.blockService.BlockIPWithPorts(blockIPWithPorts)
	if isBanned == false {
		p.logger.Info(fmt.Sprintf("IP %s are not blocked (%s) (group:%s): %s. Err: %s", result.ip, rule.Name, rule.Group.Name, entry.Message, err.Error()))
		p.sendNotifyError(&bruteForceProtectionNotify{
			rule:         rule,
			ip:           result.ip,
			ports:        l4Ports,
			messages:     groupResult.LastLogs,
			blockIPCount: groupResult.BlockIPCount,
			time:         entry.Time,
			fields:       result.fields,
			blockSec:     groupResult.BlockSec,
			err:          err,
		})
		return
	}

	p.logger.Info(fmt.Sprintf("Block IP %s detected (%s) (group:%s): %s", result.ip, rule.Name, rule.Group.Name, entry.Message))
	p.sendNotifySuccess(&bruteForceProtectionNotify{
		rule:         rule,
		ip:           result.ip,
		ports:        l4Ports,
		messages:     groupResult.LastLogs,
		blockIPCount: groupResult.BlockIPCount,
		time:         entry.Time,
		fields:       result.fields,
		blockSec:     groupResult.BlockSec,
		err:          err,
	})
}

func (p *bruteForceProtection) analyzeRule(rule *brute_force_protection.Rule, message string) bruteForceProtectionAnalyzeRuleReturn {
	result := bruteForceProtectionAnalyzeRuleReturn{
		found:     false,
		fields:    []*regexField{},
		ip:        nil,
		partition: nil,
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
				continue
			}
			ipText := message[start:end]
			ip := net.ParseIP(ipText)
			if ip == nil {
				p.logger.Error(fmt.Sprintf("Failed to parse ip: %s", ipText))
				continue
			}
			result.ip = ip

			if pattern.Partition != nil {
				start, end, err := getValueStartEndByRegexIndex(int(pattern.Partition.Value), idx)
				if err != nil {
					p.logger.Error(fmt.Sprintf("brute_force_protection. Failed to get partition value: %s", err))
					result.ip = nil
					continue
				}
				partition := pattern.Partition.Normalize(message[start:end])
				if pattern.Partition.Type != nil {
					if accepts, partitionCode := pattern.Partition.Type.Accepts(partition); !accepts {
						continue
					} else {
						partition = partitionCode
					}
					if pattern.Partition.Type.IsPartitioned() {
						result.partition = &partition
					}
				} else {
					result.partition = &partition
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
				p.logger.Error(fmt.Sprintf("brute_force_protection. analyzeRule len(pattern.Values) != len(result.fields): %d != %d", len(pattern.Values), len(result.fields)))
				result.fields = []*regexField{}
				result.partition = nil
				result.ip = nil
				continue
			}

			result.found = true
			return result
		}
	}

	return result
}

func (p *bruteForceProtection) sendNotifySuccess(notify *bruteForceProtectionNotify) {
	if !p.notificationPolicy.IsNotify(notify.rule) {
		return
	}

	groupName := notify.rule.Group.Name

	subject := i18n.Lang.T("alert.bruteForceProtection.subject", map[string]any{
		"Name":      notify.rule.Name,
		"GroupName": groupName,
		"IP":        notify.ip,
	})

	p.sendNotify(subject, notify)
}

func (p *bruteForceProtection) sendNotifyError(notify *bruteForceProtectionNotify) {
	if !notify.rule.IsNotification {
		return
	}

	groupName := notify.rule.Group.Name

	subject := i18n.Lang.T("alert.bruteForceProtection.subject-error", map[string]any{
		"Name":      notify.rule.Name,
		"GroupName": groupName,
		"IP":        notify.ip,
	})

	p.sendNotify(subject, notify)
}

func (p *bruteForceProtection) sendNotify(subject string, notify *bruteForceProtectionNotify) {
	if !notify.rule.IsNotification {
		return
	}

	groupMessage := notify.rule.Group.Message + "\n\n"

	text := subject + "\n\n" + groupMessage + notify.rule.Message + "\n\n"
	if notify.err != nil {
		text += i18n.Lang.T("alert.bruteForceProtection.error", map[string]any{
			"Error": notify.err.Error(),
		}) + "\n"
	}

	ipInfo, err := p.ipInfo(notify.ip.String())
	if err != nil {
		ipInfo = notify.ip.String()
		p.logger.Error(fmt.Sprintf("Failed to get geoip info for ip %s: %s", notify.ip, err))
	}

	text += "IP: " + ipInfo + "\n"
	if len(notify.ports) > 0 {
		var ports []string
		for _, port := range notify.ports {
			ports = append(ports, port.ToString())
		}
		text += i18n.Lang.T("ports", map[string]any{
			"Ports": strings.Join(ports, ", "),
		}) + "\n"
	}
	text += i18n.Lang.T("blockSec", map[string]any{
		"BlockSec": format.HumanDuration(time.Duration(notify.blockSec) * time.Second),
	}) + "\n"
	text += i18n.Lang.T("time", map[string]any{
		"Time": notify.time,
	}) + "\n"
	for _, field := range notify.fields {
		text += fmt.Sprintf("%s: %s\n", field.name, field.value)
	}
	text += i18n.Lang.T("blockIPCount", map[string]any{
		"Count": notify.blockIPCount,
	}) + "\n"
	text += "\n" + i18n.Lang.T("log", map[string]any{
		"Count": len(notify.messages),
	}) + "\n"
	for _, message := range notify.messages {
		text += message + "\n\n"
	}
	p.notify.SendAsync(notifications.Message{Subject: subject, Body: text})
}
