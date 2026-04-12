package geoip

import (
	"context"
)

type falseGeoIP struct {
}

func NewFalseGeoIP() GeoIP {
	return &falseGeoIP{}
}

func (g *falseGeoIP) Info(ip string) (string, error) {
	return ip, nil
}

func (g *falseGeoIP) Run(_ context.Context) {

}

func (g *falseGeoIP) Refresh(_ context.Context) error {
	return nil
}

func (g *falseGeoIP) Close() error {
	return nil
}
