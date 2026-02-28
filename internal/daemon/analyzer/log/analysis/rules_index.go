package analysis

import (
	"errors"
	"fmt"

	config2 "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config/brute_force_protection"
)

type RulesIndex struct {
	byKey map[indexKey]RulesBucket
}

type indexKey struct {
	source config2.SourceType
	val    string
}

func (idx *RulesIndex) Add(source *config2.Source) error {
	if source.AlertRule == nil && source.BruteForceProtectionRule == nil {
		return fmt.Errorf("no alert rule or brute force protection rule")
	}

	key, err := generateIndexKeyBySource(source)
	if err != nil {
		return err
	}

	if _, ok := idx.byKey[key]; !ok {
		idx.byKey[key] = newRulesBucket()
	}

	if source.AlertRule != nil {
		idx.byKey[key].addAlertRule(source.AlertRule)
	}

	if source.BruteForceProtectionRule != nil {
		idx.byKey[key].addBruteForceProtectionRule(source.BruteForceProtectionRule)
	}

	return nil
}

func (idx *RulesIndex) Alerts(entry *Entry) ([]*config2.AlertRule, error) {
	rules := make([]*config2.AlertRule, 0)

	keys, err := generateIndexKeysByEntry(entry)
	if err != nil {
		return rules, err
	}

	for _, key := range keys {
		b, ok := idx.byKey[key]
		if !ok {
			continue
		}

		rules = append(rules, b.Alerts()...)
	}

	return rules, nil
}

func (idx *RulesIndex) BruteForceProtections(entry *Entry) ([]*brute_force_protection.Rule, error) {
	rules := make([]*brute_force_protection.Rule, 0)

	keys, err := generateIndexKeysByEntry(entry)
	if err != nil {
		return rules, err
	}

	for _, key := range keys {
		b, ok := idx.byKey[key]
		if !ok {
			continue
		}

		rules = append(rules, b.BruteForceProtectionRules()...)
	}

	return rules, nil
}

func NewRulesIndex() *RulesIndex {
	return &RulesIndex{byKey: make(map[indexKey]RulesBucket)}
}

func generateIndexKeyBySource(source *config2.Source) (indexKey, error) {
	switch source.Type {
	case config2.SourceTypeJournal:
		match := source.Journal.JournalctlMatch()
		if source.Journal.Field == "" || source.Journal.Match == "" {
			return indexKey{}, errors.New("journalctl match is empty")
		}
		return indexKey{source: source.Type, val: match}, nil
	case config2.SourceTypeFile:
		return indexKey{source: source.Type, val: source.File.Path}, nil
	}

	return indexKey{}, errors.New(fmt.Sprintf("unknown source type: %s", source.Type))
}

func generateIndexKeysByEntry(entry *Entry) ([]indexKey, error) {
	var keys []indexKey

	switch entry.Source {
	case config2.SourceTypeJournal:
		keys = append(keys, indexKey{source: entry.Source, val: string(config2.JournalFieldSystemdUnit) + "=" + entry.Unit})
		keys = append(keys, indexKey{source: entry.Source, val: string(config2.JournalFieldSyslogIdentifier) + "=" + entry.SyslogIdentifier})
		return keys, nil
	case config2.SourceTypeFile:
		keys = append(keys, indexKey{source: entry.Source, val: entry.File})
		return keys, nil
	}

	return []indexKey{}, errors.New(fmt.Sprintf("unknown source type: %s", entry.Source))
}
