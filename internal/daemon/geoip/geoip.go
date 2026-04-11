package geoip

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"time"

	"git.kor-elf.net/kor-elf-shield/geoip2"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

type Info func(ip string) (string, error)

type GeoIP interface {
	Info(ip string) (string, error)
	Run(ctx context.Context)
	Close() error
}

type geoIP struct {
	config *Config
	logger log.Logger
}

func New(config *Config, logger log.Logger) GeoIP {
	return &geoIP{
		config: config,
		logger: logger,
	}
}

func (g *geoIP) Info(ip string) (string, error) {
	if g.config.GeoIP == nil {
		return ip, fmt.Errorf("geoip is not initialized")
	}

	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return ip, err
	}

	info, err := g.config.GeoIP.Info(addr)
	if err != nil {
		if errors.Is(err, geoip2.ErrNotFound) {
			g.logger.Warn(fmt.Sprintf("failed to get geoip info for ip %s: %v", ip, err))
			return ip, nil
		}
		return ip, err
	}
	g.logger.Debug(fmt.Sprintf("geoip info for ip %s: %s", ip, info.ToString()))

	return info.ToString(), nil
}

func (g *geoIP) Run(ctx context.Context) {
	g.logger.Debug("geoip service started")
	go func() {
		ticker := time.NewTicker(g.config.Interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				g.logger.Debug("refreshing geoip data")
				if err := g.config.GeoIP.Refresh(ctx); err != nil {
					g.logger.Error(fmt.Sprintf("failed to refresh geoip data: %v", err))
				}
				g.logger.Debug("geoip data refreshed")
			}
		}
	}()
}

func (g *geoIP) Close() error {
	g.logger.Debug("geoip service stopped")
	return g.config.GeoIP.Close()
}
