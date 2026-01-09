package docker_monitor

type Config struct {
	Path         string
	RuleStrategy RuleStrategy
}

type RuleStrategy int8

const (
	RuleStrategyRebuild RuleStrategy = iota + 1
)
