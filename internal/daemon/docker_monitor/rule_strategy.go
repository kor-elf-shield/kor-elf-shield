package docker_monitor

import (
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor/client"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor/rule_strategy"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

func newRuleStrategy(config *Config, dockerClient client.Docker, logger log.Logger) (rule_strategy.Strategy, error) {
	generate := rule_strategy.NewGenerator(dockerClient, logger)

	switch config.RuleStrategy {
	case RuleStrategyRebuild:
		return rule_strategy.NewRebuildStrategy(generate), nil
	}

	return nil, nil
}
