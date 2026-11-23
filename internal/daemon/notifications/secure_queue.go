package notifications

import (
	"encoding/json"
	"fmt"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/crypto/envelope"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/entity"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/secret"
)

const encryptedSubjectMarker = "__encrypted_v1__"

type QueueCodec interface {
	Encode(message Message) (entity.NotificationsQueue, error)
	Decode(item entity.NotificationsQueue) (Message, error)
}

type encryptedQueueCodec struct {
	cipher envelope.Cipher
}

type queuePayload struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func NewEncryptedQueueCodec(secretProvider secret.Provider) (QueueCodec, error) {
	c, err := envelope.NewAESGCM(secretProvider, "notifications-queue", "default")
	if err != nil {
		return nil, err
	}
	return &encryptedQueueCodec{cipher: c}, nil
}

func (c *encryptedQueueCodec) Encode(message Message) (entity.NotificationsQueue, error) {
	payload, err := json.Marshal(queuePayload{
		Subject: message.Subject,
		Body:    message.Body,
	})
	if err != nil {
		return entity.NotificationsQueue{}, fmt.Errorf("marshal queue payload: %w", err)
	}

	encrypted, err := c.cipher.Encrypt(payload)
	if err != nil {
		return entity.NotificationsQueue{}, err
	}

	return entity.NotificationsQueue{
		Subject: encryptedSubjectMarker,
		Body:    string(encrypted),
	}, nil
}

func (c *encryptedQueueCodec) Decode(item entity.NotificationsQueue) (Message, error) {
	// Backward compatibility: old plain records.
	if item.Subject != encryptedSubjectMarker {
		return Message{
			Subject: item.Subject,
			Body:    item.Body,
		}, nil
	}

	plain, err := c.cipher.Decrypt([]byte(item.Body))
	if err != nil {
		return Message{}, err
	}

	var payload queuePayload
	if err := json.Unmarshal(plain, &payload); err != nil {
		return Message{}, fmt.Errorf("unmarshal queue payload: %w", err)
	}

	return Message{
		Subject: payload.Subject,
		Body:    payload.Body,
	}, nil
}
