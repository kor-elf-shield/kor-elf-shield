package secret

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

type Provider interface {
	Secret() ([]byte, error)
	Close()
}

type programKey struct {
	Secret string `json:"secret"`
}

func NewSecret(pathFileSecret string, secretCache int, logger log.Logger) (Provider, error) {
	if pathFileSecret == "" {
		return nil, fmt.Errorf("pathFileSecret is empty")
	}

	if secretCache < 0 {
		return nil, fmt.Errorf("secretCache must be greater than or equal to 0")
	}

	if logger == nil {
		return nil, fmt.Errorf("logger is nil")
	}

	return &secret{
		pathFileSecret: pathFileSecret,
		secretCache:    secretCache,
		logger:         logger,
	}, nil
}

type secret struct {
	pathFileSecret string
	secretCache    int
	logger         log.Logger

	mu           sync.Mutex
	cachedSecret []byte
	cacheUntil   time.Time
	cacheTimer   *time.Timer
	closed       bool
	lastAccess   time.Time
	inFlight     int
}

func (s *secret) Secret() ([]byte, error) {
	now := time.Now()

	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return nil, fmt.Errorf("secret provider is closed")
	}

	s.inFlight++
	defer func() {
		s.mu.Lock()
		s.inFlight--
		s.mu.Unlock()
	}()

	if s.secretCache > 0 && len(s.cachedSecret) > 0 && now.Before(s.cacheUntil) {
		s.lastAccess = now
		s.cacheUntil = now.Add(time.Duration(s.secretCache) * time.Second)
		s.resetCacheTimerLocked()
		secretCopy := cloneBytes(s.cachedSecret)
		s.mu.Unlock()
		return secretCopy, nil
	}
	s.mu.Unlock()

	generated, err := s.generateHardwareKey()
	if err != nil {
		return nil, err
	}

	if s.secretCache == 0 {
		return generated, nil
	}

	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		zeroBytes(generated)
		return nil, fmt.Errorf("secret provider is closed")
	}

	zeroBytes(s.cachedSecret)
	s.cachedSecret = cloneBytes(generated)
	now = time.Now()
	s.lastAccess = now
	s.cacheUntil = now.Add(time.Duration(s.secretCache) * time.Second)
	s.resetCacheTimerLocked()
	s.mu.Unlock()

	return generated, nil
}

func (s *secret) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return
	}
	s.closed = true

	if s.cacheTimer != nil {
		s.cacheTimer.Stop()
		s.cacheTimer = nil
	}

	zeroBytes(s.cachedSecret)
	s.cachedSecret = nil
	s.cacheUntil = time.Time{}
}

func (s *secret) resetCacheTimerLocked() {
	if s.secretCache <= 0 {
		return
	}

	if s.cacheTimer != nil {
		s.cacheTimer.Stop()
	}

	ttl := time.Duration(s.secretCache) * time.Second
	s.cacheTimer = time.AfterFunc(ttl, s.clearExpiredCache)
}

func (s *secret) clearExpiredCache() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return
	}

	if len(s.cachedSecret) == 0 {
		s.cacheTimer = nil
		return
	}

	if s.inFlight > 0 {
		s.resetCacheTimerLocked()
		return
	}

	ttl := time.Duration(s.secretCache) * time.Second
	idleFor := time.Since(s.lastAccess)
	if idleFor < ttl {
		s.cacheUntil = s.lastAccess.Add(ttl)
		s.resetCacheTimerLocked()
		return
	}

	zeroBytes(s.cachedSecret)
	s.cachedSecret = nil
	s.cacheUntil = time.Time{}
	s.lastAccess = time.Time{}
	s.cacheTimer = nil
}

func (s *secret) generateHardwareKey() ([]byte, error) {
	appSecret, err := s.fileSecret()
	if err != nil {
		return nil, err
	}

	// We guarantee that the local copy of the string is cleared upon exit from the function.
	defer func() {
		appSecret = ""
	}()

	var fingerprint strings.Builder

	// Collecting hardware data
	if machineID, err := getMachineID(); err == nil {
		fingerprint.Write(machineID)
	} else {
		s.logger.Debug("failed to read /etc/machine-id and /var/lib/dbus/machine-id")
	}

	if productUUID, err := getProductUUID(); err == nil {
		fingerprint.Write(productUUID)
	} else {
		s.logger.Debug("failed to read /sys/class/dmi/id/product_uuid")
	}

	if fingerprint.Len() == 0 {
		return nil, fmt.Errorf("failed to get hardware identifiers of the system")
	}

	// We bind the hardware to a unique secret specific to a given server
	fingerprint.WriteString(appSecret)

	hash := sha256.New()
	if _, err := io.WriteString(hash, strings.TrimSpace(fingerprint.String())); err != nil {
		return nil, fmt.Errorf("failed to write to hash: %v", err)
	}

	return hash.Sum(nil), nil
}

func (s *secret) fileSecret() (string, error) {
	if _, err := os.Stat(s.pathFileSecret); os.IsNotExist(err) {
		dir := filepath.Dir(s.pathFileSecret)
		if err := os.MkdirAll(dir, 0700); err != nil {
			return "", fmt.Errorf("failed to create the secret directory: %v", err)
		}

		// Generate a cryptographically secure random secret (32 bytes -> 64 hex characters)
		appSecret, err := generateAppSecret()
		if err != nil {
			return "", fmt.Errorf("failed to generate app secret: %v", err)
		}

		programKeyData := programKey{Secret: appSecret}
		fileBytes, _ := json.MarshalIndent(programKeyData, "", "  ")

		// We write the file with strict 0400 permissions (read for root only).
		if err := os.WriteFile(s.pathFileSecret, fileBytes, 0400); err != nil {
			return "", fmt.Errorf("failed to write the secret file: %v", err)
		}

		s.logger.Debug(fmt.Sprintf("[+] Generated a new unique app_secret in %s", s.pathFileSecret))
	}

	file, err := os.ReadFile(s.pathFileSecret)
	if err != nil {
		return "", fmt.Errorf("failed to read the secret file: %v", err)
	}

	var key programKey
	if err := json.Unmarshal(file, &key); err != nil {
		return "", fmt.Errorf("failed to parse the secret file: %v", err)
	}

	if key.Secret == "" {
		return "", fmt.Errorf("the secret in the file is empty")
	}

	return key.Secret, nil
}

func cloneBytes(src []byte) []byte {
	if len(src) == 0 {
		return nil
	}
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}

/**
 * Zero out a byte slice
 */
func zeroBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

/**
 * Generate a cryptographically secure random secret (32 bytes -> 64 hex characters)
 */
func generateAppSecret() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("random byte generation error: %v", err)
	}
	return hex.EncodeToString(bytes), nil
}

/**
 * Get the machine ID
 */
func getMachineID() ([]byte, error) {
	// Physical server ID
	if machineID, err := os.ReadFile("/etc/machine-id"); err == nil {
		return machineID, nil
	}

	// VPS/VDS ID
	if machineID, err := os.ReadFile("/var/lib/dbus/machine-id"); err == nil {
		return machineID, nil
	}

	return nil, fmt.Errorf("failed to read /etc/machine-id and /var/lib/dbus/machine-id")
}

/**
 * Get Motherboard UUID (available for bare-metal servers and KVM/VMware VDS)
 */
func getProductUUID() ([]byte, error) {
	productUUID, err := os.ReadFile("/sys/class/dmi/id/product_uuid")
	if err == nil {
		return nil, fmt.Errorf("failed to read /sys/class/dmi/id/product_uuid")
	}

	return productUUID, nil
}
