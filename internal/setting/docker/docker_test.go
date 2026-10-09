package docker

import (
	"path/filepath"
	"testing"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/fakers"
)

func TestDocker_InitSetting_Success(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "docker.toml")

	if err := testutil.SaveTomlEmpty(path); err != nil {
		t.Fatalf("failed to create toml file: %v", err)
	}

	_, err := InitSetting(path)
	if err != nil {
		t.Errorf("InitSetting() should not return an error, got %v", err)
		return
	}
}

func TestDocker_InitSetting_Error(t *testing.T) {
	t.Run("invalid path", func(t *testing.T) {
		if _, err := InitSetting("docker.toml"); err == nil {
			t.Errorf("InitSetting() should return an error for non-absolute path")
		}
	})

	t.Run("file does not exist", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "not_exists.toml")

		if _, err := InitSetting(path); err == nil {
			t.Errorf("InitSetting() should return an error for missing file")
		}
	})

	t.Run("file not TOML", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "not_toml.txt")
		if err := testutil.SaveTomlEmpty(path); err != nil {
			t.Errorf("testutil.SaveTomlEmpty() failed: %v", err)
		}

		if _, err := InitSetting(path); err == nil {
			t.Errorf("InitSetting() should return an error for non-TOML file")
		}
	})
}

func TestDocker_Validate_Success(t *testing.T) {
	s := defaultDockerForTest()
	if err := s.Validate(); err != nil {
		t.Errorf("Validate() should not return an error, got %v", err)
	}
}

func TestDocker_ToRuleStrategy_Success(t *testing.T) {
	s := defaultDockerForTest()
	if ruleStrategy, err := s.ToRuleStrategy(); err != nil {
		t.Errorf("ToRuleStrategy() should not return an error, got %v", err)
	} else if !testRuleStrategy(ruleStrategy, s.RuleStrategy) {
		t.Errorf("ToRuleStrategy() should return the correct rule strategy, got %v", ruleStrategy)
	}
}

func TestDocker_ToRuleStrategy_Error(t *testing.T) {
	t.Run("empty rule strategy", func(t *testing.T) {
		s := defaultDockerForTest()
		s.RuleStrategy = ""
		if _, err := s.ToRuleStrategy(); err == nil {
			t.Errorf("ToRuleStrategy() should return an error for empty rule strategy")
		}
	})

	t.Run("invalid rule strategy", func(t *testing.T) {
		s := defaultDockerForTest()
		s.RuleStrategy = "invalid"
		if _, err := s.ToRuleStrategy(); err == nil {
			t.Errorf("ToRuleStrategy() should return an error for invalid rule strategy")
		}
	})
}

func defaultDockerForTest() *Setting {
	return &Setting{
		Enabled:      true,
		RuleStrategy: fakeRandRuleStrategy(),
	}
}

func fakeRandRuleStrategy() string {
	items := []string{"rebuild", "incremental"}

	return fakers.RandItem(items)
}

func testRuleStrategy(ruleStrategy docker_monitor.RuleStrategy, configRuleStrategy string) bool {
	switch ruleStrategy {
	case docker_monitor.RuleStrategyIncremental:
		return configRuleStrategy == "incremental"
	case docker_monitor.RuleStrategyRebuild:
		return configRuleStrategy == "rebuild"
	default:
		return false
	}
}
