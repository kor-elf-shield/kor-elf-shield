package guard

import (
	"context"
	"fmt"
	"sync"
	"time"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/guard/config"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/info"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/notifications"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/i18n"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

type firewallGuardTarget interface {
	HasRules() (bool, error)
	Reload(daemonInfo info.Info) error
}

type RulesGuard interface {
	Run(daemonInfo info.Info, ctx context.Context)
	Close() error
}

type rulesGuard struct {
	config   *config.GuardConfig
	firewall firewallGuardTarget
	notify   notifications.Notifications
	logger   log.Logger

	mu     sync.Mutex
	cancel context.CancelFunc
}

func NewRulesGuard(config *config.GuardConfig, firewall firewallGuardTarget, notify notifications.Notifications, logger log.Logger) RulesGuard {
	return &rulesGuard{
		config:   config,
		firewall: firewall,
		notify:   notify,
		logger:   logger,
	}
}

func (g *rulesGuard) Run(daemonInfo info.Info, ctx context.Context) {
	if !g.config.Enable {
		g.logger.Debug("firewall rules guard is disabled")
		return
	}

	g.logger.Debug("firewall rules guard is enabled")
	guardCtx, cancel := context.WithCancel(ctx)

	g.mu.Lock()
	g.cancel = cancel
	g.mu.Unlock()

	go g.run(daemonInfo, guardCtx)
}

func (g *rulesGuard) Close() error {
	g.mu.Lock()
	cancel := g.cancel
	g.cancel = nil
	g.mu.Unlock()

	if cancel != nil {
		cancel()
	}

	return nil
}

func (g *rulesGuard) run(daemonInfo info.Info, ctx context.Context) {
	interval := time.Duration(g.config.Interval) * time.Second

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			g.checkRules(daemonInfo)

		case <-ctx.Done():
			g.logger.Debug("firewall rules guard stopped")
			return
		}
	}
}

func (g *rulesGuard) checkRules(daemonInfo info.Info) {
	hasRules, err := g.firewall.HasRules()
	if err != nil {
		g.logger.Error(fmt.Sprintf("failed to check firewall rules: %s", err))
		if g.config.Notifications {
			message := notifications.Message{
				Subject: i18n.Lang.T("firewall_rules_not_found"),
				Body: i18n.Lang.T("firewall_rules_not_found_body_check_error", map[string]interface{}{
					"Error": err.Error(),
				}),
			}
			g.notify.SendAsync(message)
		}
		return
	}

	if hasRules {
		g.logger.Debug("firewall rules exists")
		return
	}

	g.logger.Warn("firewall rules not found")
	if g.config.Recovery {
		if err := g.firewall.Reload(daemonInfo); err != nil {
			g.logger.Error(fmt.Sprintf("failed to recover firewall rules: %s", err))
			if g.config.Notifications {
				message := notifications.Message{
					Subject: i18n.Lang.T("firewall_rules_not_found"),
					Body: i18n.Lang.T("firewall_rules_not_found_body_recover_error", map[string]interface{}{
						"Error": err.Error(),
					}),
				}
				g.notify.SendAsync(message)
			}
			return
		}
		g.logger.Warn("firewall rules recovered")
		if g.config.Notifications {
			message := notifications.Message{
				Subject: i18n.Lang.T("firewall_rules_not_found"),
				Body:    i18n.Lang.T("firewall_rules_not_found_body_recover_success"),
			}
			g.notify.SendAsync(message)
		}
	} else if g.config.Notifications {
		message := notifications.Message{
			Subject: i18n.Lang.T("firewall_rules_not_found"),
			Body:    i18n.Lang.T("firewall_rules_not_found_body"),
		}
		g.notify.SendAsync(message)
	}
}
