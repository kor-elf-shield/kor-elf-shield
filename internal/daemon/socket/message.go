package socket

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Message struct {
	Command   string            `json:"command"`
	Args      map[string]string `json:"args"`
	Nonce     string            `json:"nonce"`
	Timestamp int64             `json:"timestamp"`
	Signature string            `json:"signature"`
}

func BuildSignedMessage(command string, args map[string]string, secret []byte) (Message, error) {
	if command == "" {
		return Message{}, fmt.Errorf("command is empty")
	}
	if len(secret) == 0 {
		return Message{}, fmt.Errorf("secret is empty")
	}
	if args == nil {
		args = map[string]string{}
	}

	nonce, err := generateNonceHex(16)
	if err != nil {
		return Message{}, fmt.Errorf("generate nonce: %w", err)
	}
	ts := time.Now().Unix()

	payload := canonicalPayload(command, args, nonce, ts)
	signatureHex := signPayloadHex(payload, secret)

	return Message{
		Command:   command,
		Args:      args,
		Nonce:     nonce,
		Timestamp: ts,
		Signature: signatureHex,
	}, nil
}

func signPayloadHex(payload string, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

func verifyMessageSignature(msg Message, secret []byte) error {
	if len(secret) == 0 {
		return errors.New("empty secret")
	}
	if msg.Signature == "" {
		return errors.New("signature is empty")
	}
	if msg.Nonce == "" {
		return errors.New("nonce is empty")
	}
	if msg.Timestamp <= 0 {
		return errors.New("timestamp is invalid")
	}
	if msg.Timestamp < time.Now().Unix()-15 {
		return errors.New("timestamp is too old")
	}

	gotSig, err := hex.DecodeString(msg.Signature)
	if err != nil {
		return errors.New("signature is not valid hex")
	}

	expected := computeMessageHMAC(msg.Command, msg.Args, msg.Nonce, msg.Timestamp, secret)
	if subtle.ConstantTimeCompare(gotSig, expected) != 1 {
		return errors.New("signature mismatch")
	}

	return nil
}

func computeMessageHMAC(command string, args map[string]string, nonce string, ts int64, secret []byte) []byte {
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(canonicalPayload(command, args, nonce, ts)))
	return mac.Sum(nil)
}

func canonicalPayload(command string, args map[string]string, nonce string, ts int64) string {
	keys := make([]string, 0, len(args))
	for k := range args {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	b.WriteString(command)
	b.WriteByte('\n')
	b.WriteString("nonce=")
	b.WriteString(nonce)
	b.WriteByte('\n')
	b.WriteString("ts=")
	b.WriteString(strconv.FormatInt(ts, 10))
	b.WriteByte('\n')

	for _, k := range keys {
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(args[k])
		b.WriteByte('\n')
	}

	return b.String()
}

func generateNonceHex(size int) (string, error) {
	if size <= 0 {
		return "", fmt.Errorf("invalid nonce size")
	}

	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
