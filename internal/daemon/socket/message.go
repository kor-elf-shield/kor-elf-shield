package socket

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
)

type Message struct {
	Command   string            `json:"command"`
	Args      map[string]string `json:"args"`
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

	payload := canonicalPayload(command, args)
	signatureHex := signPayloadHex(payload, secret)

	return Message{
		Command:   command,
		Args:      args,
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

	gotSig, err := hex.DecodeString(msg.Signature)
	if err != nil {
		return errors.New("signature is not valid hex")
	}

	expected := computeMessageHMAC(msg.Command, msg.Args, secret)
	if subtle.ConstantTimeCompare(gotSig, expected) != 1 {
		return errors.New("signature mismatch")
	}

	return nil
}

func computeMessageHMAC(command string, args map[string]string, secret []byte) []byte {
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write([]byte(canonicalPayload(command, args)))
	return mac.Sum(nil)
}

func canonicalPayload(command string, args map[string]string) string {
	keys := make([]string, 0, len(args))
	for k := range args {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	b.WriteString(command)
	b.WriteByte('\n')

	for _, k := range keys {
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(args[k])
		b.WriteByte('\n')
	}

	return b.String()
}
