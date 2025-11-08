package validate

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

func PathFile(path string, parameterName string) error {
	if path == "" {
		return fmt.Errorf("%s is empty", parameterName)
	}
	if strings.Contains(path, "..") || strings.Contains(path, "//") {
		return fmt.Errorf("%s must not contain '..' or '//'", parameterName)
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9_\-./]*$`)
	if !re.MatchString(path) {
		return fmt.Errorf("%s must not contain special characters", parameterName)
	}
	if !filepath.IsAbs(path) {
		return fmt.Errorf("%s must be absolute path", parameterName)
	}
	return nil
}

func IsTomlFile(path string, parameterName string) error {
	if err := PathFile(path, parameterName); err != nil {
		return err
	}
	if !strings.HasSuffix(strings.ToLower(path), ".toml") {
		return fmt.Errorf("invalid %s. Must be .toml", parameterName)
	}
	return nil
}

func NftLimitRate(rate string, parameterName string) error {
	reRate := regexp.MustCompile(`(?i)^\s*(over\s+)?([0-9]+(\.[0-9]+)?)\s*(k?m?bytes|packets)?\s*/\s*(second|s|minute|min|hour|h|day|d|week|w)\s*(burst\s+([0-9]+(\.[0-9]+)?)\s*(k?m?bytes|packets)?)?\s*$`)
	if reRate.MatchString(rate) {
		return nil
	}

	return fmt.Errorf("%s must be in format <number> <unit>", parameterName)
}

func Name(name string, parameterName string) error {
	if name == "" {
		return fmt.Errorf("%s is empty", parameterName)
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9_\-]{1,32}$`)
	if !re.MatchString(name) {
		return fmt.Errorf("%s must not contain special characters", parameterName)
	}
	return nil
}

func Port(port int, parameterName string) error {
	if port < 0 || port > 65535 {
		return fmt.Errorf("%s must be in range 0-65535", parameterName)
	}
	return nil
}
