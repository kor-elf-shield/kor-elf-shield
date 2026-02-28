package analysis

import (
	config2 "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config/brute_force_protection"
)

type RulesBucket interface {
	Alerts() []*config2.AlertRule
	BruteForceProtectionRules() []*brute_force_protection.Rule

	addAlertRule(rule *config2.AlertRule)
	addBruteForceProtectionRule(rule *brute_force_protection.Rule)
}

type rulesBucket struct {
	alerts                    []*config2.AlertRule
	bruteForceProtectionRules []*brute_force_protection.Rule
}

func (rb *rulesBucket) Alerts() []*config2.AlertRule {
	return rb.alerts
}

func (rb *rulesBucket) BruteForceProtectionRules() []*brute_force_protection.Rule {
	return rb.bruteForceProtectionRules
}

func (rb *rulesBucket) addAlertRule(rule *config2.AlertRule) {
	rb.alerts = append(rb.alerts, rule)
}

func (rb *rulesBucket) addBruteForceProtectionRule(rule *brute_force_protection.Rule) {
	rb.bruteForceProtectionRules = append(rb.bruteForceProtectionRules, rule)
}

func newRulesBucket() RulesBucket {
	return &rulesBucket{
		alerts: make([]*config2.AlertRule, 0),
	}
}
