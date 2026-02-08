package analysis

import (
	"errors"
	"fmt"

	config2 "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
)

type AlertRuleIndex struct {
	byKey map[string][]*config2.AlertRule
}

func (idx *AlertRuleIndex) Add(source *config2.Source) error {
	key, err := generateIndexKeyBySource(source)
	if err != nil {
		return err
	}

	idx.byKey[key] = append(idx.byKey[key], source.AlertRule)
	return nil
}

func (idx *AlertRuleIndex) Rules(entry *Entry) ([]*config2.AlertRule, error) {
	var rules []*config2.AlertRule

	keys, err := generateIndexKeysByEntry(entry)
	if err != nil {
		return rules, err
	}

	for _, key := range keys {
		rules = append(rules, idx.byKey[key]...)
	}

	return rules, nil
}

func NewAlertRuleIndex() AlertRuleIndex {
	return AlertRuleIndex{byKey: make(map[string][]*config2.AlertRule)}
}

func generateIndexKeyBySource(source *config2.Source) (string, error) {
	switch source.Type {
	case config2.SourceTypeJournal:
		match := source.Journal.JournalctlMatch()
		if source.Journal.Field == "" || source.Journal.Match == "" {
			return "", errors.New("journalctl match is empty")
		}
		return string(source.Type) + ":" + match, nil
	}

	return "", errors.New(fmt.Sprintf("unknown source type: %s", source.Type))
}

func generateIndexKeysByEntry(entry *Entry) ([]string, error) {
	var keys []string

	switch entry.Source {
	case config2.SourceTypeJournal:
		source := string(entry.Source) + ":"

		keys = append(keys, source+string(config2.JournalFieldSystemdUnit)+"="+entry.Unit)
		keys = append(keys, source+string(config2.JournalFieldSyslogIdentifier)+"="+entry.SyslogIdentifier)

		return keys, nil
	}

	return []string{}, errors.New(fmt.Sprintf("unknown source type: %s", entry.Source))
}
