package envelope

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/secret"
)

const (
	envelopeVersion = 1
	nonceSize       = 12
)

type Cipher interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(envelopeJSON []byte) ([]byte, error)
}

type aesGCMCipher struct {
	secretProvider secret.Provider
	domain         string
	keyID          string
}

type payloadEnvelope struct {
	V   int    `json:"v"`
	D   string `json:"d"`
	KID string `json:"kid"`
	N   string `json:"n"`
	C   string `json:"c"`
}

func NewAESGCM(secretProvider secret.Provider, domain string, keyID string) (Cipher, error) {
	if secretProvider == nil {
		return nil, errors.New("secret provider is nil")
	}
	if domain == "" {
		return nil, errors.New("domain is empty")
	}
	if keyID == "" {
		keyID = "default"
	}

	return &aesGCMCipher{
		secretProvider: secretProvider,
		domain:         domain,
		keyID:          keyID,
	}, nil
}

func (c *aesGCMCipher) Encrypt(plaintext []byte) ([]byte, error) {
	secretBytes, err := c.secretProvider.Secret()
	if err != nil {
		return nil, fmt.Errorf("get secret: %w", err)
	}

	key := deriveKey(secretBytes, c.domain, c.keyID)
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, fmt.Errorf("new aes cipher: %w", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("new gcm: %w", err)
	}

	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("generate nonce: %w", err)
	}

	aad := additionalData(envelopeVersion, c.domain, c.keyID)
	ciphertext := aead.Seal(nil, nonce, plaintext, aad)

	env := payloadEnvelope{
		V:   envelopeVersion,
		D:   c.domain,
		KID: c.keyID,
		N:   base64.StdEncoding.EncodeToString(nonce),
		C:   base64.StdEncoding.EncodeToString(ciphertext),
	}

	raw, err := json.Marshal(env)
	if err != nil {
		return nil, fmt.Errorf("marshal envelope: %w", err)
	}

	return raw, nil
}

func (c *aesGCMCipher) Decrypt(envelopeJSON []byte) ([]byte, error) {
	var env payloadEnvelope
	if err := json.Unmarshal(envelopeJSON, &env); err != nil {
		return nil, fmt.Errorf("unmarshal envelope: %w", err)
	}

	if env.V != envelopeVersion {
		return nil, fmt.Errorf("unsupported envelope version: %d", env.V)
	}
	if env.D == "" || env.D != c.domain {
		return nil, errors.New("invalid envelope domain")
	}
	if env.KID == "" {
		return nil, errors.New("empty key id")
	}

	nonce, err := base64.StdEncoding.DecodeString(env.N)
	if err != nil {
		return nil, fmt.Errorf("decode nonce: %w", err)
	}
	if len(nonce) != nonceSize {
		return nil, fmt.Errorf("invalid nonce length: %d", len(nonce))
	}

	ciphertext, err := base64.StdEncoding.DecodeString(env.C)
	if err != nil {
		return nil, fmt.Errorf("decode ciphertext: %w", err)
	}

	secretBytes, err := c.secretProvider.Secret()
	if err != nil {
		return nil, fmt.Errorf("get secret: %w", err)
	}

	key := deriveKey(secretBytes, c.domain, env.KID)
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, fmt.Errorf("new aes cipher: %w", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("new gcm: %w", err)
	}

	plain, err := aead.Open(nil, nonce, ciphertext, additionalData(env.V, c.domain, env.KID))
	if err != nil {
		return nil, fmt.Errorf("decrypt/open failed: %w", err)
	}

	return plain, nil
}

func deriveKey(secretBytes []byte, domain string, keyID string) [32]byte {
	h := sha256.New()
	h.Write([]byte("kor-elf-shield:enc:v1:"))
	h.Write([]byte(domain))
	h.Write([]byte{0})
	h.Write([]byte(keyID))
	h.Write([]byte{0})
	h.Write(secretBytes)

	var out [32]byte
	copy(out[:], h.Sum(nil))
	return out
}

func additionalData(version int, domain string, keyID string) []byte {
	return []byte(fmt.Sprintf("v=%d;d=%s;kid=%s", version, domain, keyID))
}
