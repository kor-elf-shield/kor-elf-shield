package i18nmock

import (
	"fmt"
	"sync"
	"testing"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/i18n"
)

var langMu sync.Mutex

type fakeLang struct{}

func (f *fakeLang) T(messageID string, data ...map[string]interface{}) string {
	if len(data) > 0 {
		if p, ok := data[0]["Parameter"]; ok {
			return fmt.Sprintf("%s:%v", messageID, p)
		}
	}
	return messageID
}

func (f *fakeLang) ChangeLang(_ string) error { return nil }

// Use replaces global i18n.Lang for a single test and restores it automatically.
func Use(t testing.TB) {
	t.Helper()

	langMu.Lock()
	prev := i18n.Lang
	i18n.Lang = &fakeLang{}

	t.Cleanup(func() {
		i18n.Lang = prev
		langMu.Unlock()
	})
}
