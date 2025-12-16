package notifications

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
	"github.com/wneessen/go-mail"
)

type Message struct {
	Subject string
	Body    string
}

type Notifications interface {
	Run()
	SendAsync(message Message)
	Close() error
}

type notifications struct {
	config   Config
	logger   log.Logger
	msgQueue chan Message
	wg       sync.WaitGroup
}

func New(config Config, logger log.Logger) Notifications {
	return &notifications{
		config:   config,
		logger:   logger,
		msgQueue: make(chan Message, 100),
	}
}

func (n *notifications) Run() {
	if n.config.Enabled == false {
		n.logger.Info("Notifications are disabled")
	}

	n.wg.Add(1)
	go func() {
		defer n.wg.Done()
		for msg := range n.msgQueue {
			err := n.sendEmail(msg)
			if err != nil {
				n.logger.Error(fmt.Sprintf("failed to send email: %v", err))
			} else if n.config.Enabled {
				n.logger.Debug(fmt.Sprintf("email sent: Subject %s, Body %s", msg.Subject, msg.Body))
			}
		}
	}()
}

func (n *notifications) SendAsync(message Message) {
	select {
	case n.msgQueue <- message:
		if n.config.Enabled == false {
			n.logger.Debug(fmt.Sprintf("email sending is disabled, message was added to the queue: Subject %s, Body %s", message.Subject, message.Body))
		} else {
			n.logger.Debug(fmt.Sprintf("added to the mail sending queue: Subject %s, Body %s", message.Subject, message.Body))
		}
	default:
		n.logger.Error(fmt.Sprintf("failed to send email: queue is full"))
	}
}

func (n *notifications) Close() error {
	close(n.msgQueue)
	n.logger.Debug("We are waiting for all notifications to be sent")
	n.wg.Wait()
	n.logger.Debug("Notifications queue processed and closed")
	return nil
}

func (n *notifications) sendEmail(message Message) error {
	if n.config.Enabled == false {
		return nil
	}

	m := mail.NewMsg()
	if err := m.From(n.config.Email.From); err != nil {
		return err
	}
	if err := m.To(n.config.Email.To); err != nil {
		return err
	}
	m.Subject(message.Subject + " (" + n.config.ServerName + ")")
	m.SetBodyString(mail.TypeTextPlain, "Server: "+n.config.ServerName+"\n"+message.Body)

	client, err := newClient(n.config.Email)
	if err != nil {
		return err
	}
	defer func() { _ = client.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return client.DialAndSendWithContext(ctx, m)
}

func newClient(config Email) (*mail.Client, error) {
	options := []mail.Option{
		mail.WithPort(int(config.Port)),
		mail.WithSMTPAuth(config.AuthType),
	}

	if config.AuthType != mail.SMTPAuthNoAuth {
		options = append(options, mail.WithUsername(config.Username), mail.WithPassword(config.Password))
	}

	switch config.TLS.Mode {
	case TLSModeImplicit:
		options = append(options, mail.WithSSL())
		break
	case TLSModeStartTLS:

		switch config.TLS.Policy {
		case TLSPolicyMandatory:
			options = append(options, mail.WithTLSPolicy(mail.TLSMandatory))
			break
		case TLSPolicyOpportunistic:
			options = append(options, mail.WithTLSPolicy(mail.TLSOpportunistic))
			break
		default:
			return nil, fmt.Errorf("unknown tls policy: %s", config.TLS.Policy)
		}

		if !config.TLS.Verify {
			tlsCfg := &tls.Config{
				InsecureSkipVerify: true,
			}
			options = append(options, mail.WithTLSConfig(tlsCfg))
		}
		break
	case TLSModeNone:

		break
	default:
		return nil, fmt.Errorf("unknown tls mode: %s", config.TLS.Mode)
	}

	options = append(options, mail.WithSSL())

	return mail.NewClient(config.Host, options...)
}
