package analysis

import config2 "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"

type RulesBucket interface {
	Alerts() []*config2.AlertRule

	addRule(rule *config2.AlertRule)
}

type rulesBucket struct {
	alerts []*config2.AlertRule
}

func (rb *rulesBucket) Alerts() []*config2.AlertRule {
	return rb.alerts
}

func (rb *rulesBucket) addRule(rule *config2.AlertRule) {
	rb.alerts = append(rb.alerts, rule)
}

func newRulesBucket() RulesBucket {
	return &rulesBucket{
		alerts: make([]*config2.AlertRule, 0),
	}
}
