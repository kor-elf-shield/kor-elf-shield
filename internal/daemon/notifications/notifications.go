package notifications

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/entity"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/repository"
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
	// DBQueueSize - return size of notifications queue in db
	DBQueueSize() int
	DBQueueClear() error
	Close() error
}

type notifications struct {
	config          Config
	queueRepository repository.NotificationsQueueRepository
	logger          log.Logger
	msgQueue        chan Message
	wg              sync.WaitGroup
}

func New(config Config, queueRepository repository.NotificationsQueueRepository, logger log.Logger) Notifications {
	return &notifications{
		config:          config,
		queueRepository: queueRepository,
		logger:          logger,
		msgQueue:        make(chan Message, 100),
	}
}

func (n *notifications) Run() {
	if n.config.Enabled == false {
		n.logger.Info("Notifications are disabled")
	}

	n.wg.Add(1)
	go func() {
		defer n.wg.Done()

		ticker := time.NewTicker(time.Duration(n.config.RetryInterval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case msg, ok := <-n.msgQueue:
				if !ok {
					return
				}
				err := n.sendEmail(msg)
				if err != nil {
					n.logger.Error(fmt.Sprintf("failed to send email: %v", err))
					n.addNotificationsQueue(msg)
				} else if n.config.Enabled {
					n.logger.Debug(fmt.Sprintf("email sent: Subject %s, Body %s", msg.Subject, msg.Body))
				}
			case <-ticker.C:
				if n.config.Enabled == false || n.config.EnableRetries == false {
					continue
				}

				items, err := n.queueRepository.Get(10)
				if err != nil {
					n.logger.Error(fmt.Sprintf("failed to get notifications from the queue: %v", err))
					continue
				}

				for id, item := range items {
					err = n.sendEmail(Message{Subject: item.Subject, Body: item.Body})
					if err != nil {
						n.logger.Error(fmt.Sprintf("failed to send queued email: %v", err))
						break
					}

					err = n.queueRepository.Delete(id)
					if err != nil {
						n.logger.Error(fmt.Sprintf("failed to delete queued email from the queue: %v", err))
					}
				}
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
		n.addNotificationsQueue(message)
	}
}

func (n *notifications) DBQueueSize() int {
	count, err := n.queueRepository.Count()
	if err != nil {
		n.logger.Error(fmt.Sprintf("failed to get notifications queue size: %v", err))
		return 0
	}
	return count
}

func (n *notifications) DBQueueClear() error {
	err := n.queueRepository.Clear()
	if err != nil {
		n.logger.Error(fmt.Sprintf("failed to clear notifications queue: %v", err))
	}
	return err
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

func (n *notifications) addNotificationsQueue(message Message) {
	if n.config.Enabled == false || n.config.EnableRetries == false {
		return
	}

	err := n.queueRepository.Add(entity.NotificationsQueue{Body: message.Body, Subject: message.Subject})
	if err != nil {
		n.logger.Error(fmt.Sprintf("failed to save email to the queue: %v", err))
	}
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
