package notifications

import (
	"errors"
	"fmt"
	"strings"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/notifications"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting/validate"
	"github.com/wneessen/go-mail"
)

type Email struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	AuthType string `mapstructure:"auth_type"`

	TLSMode   string `mapstructure:"tls_mode"`
	TLSPolicy string `mapstructure:"tls_policy"`
	TLSVerify bool   `mapstructure:"tls_verify"`

	From string `mapstructure:"from"`
	To   string `mapstructure:"to"`
}

func defaultEmail() Email {
	return Email{
		Host:      "",
		Port:      0,
		Username:  "",
		Password:  "",
		AuthType:  "PLAIN",
		TLSMode:   "STARTTLS",
		TLSPolicy: "MANDATORY",
		TLSVerify: true,
		From:      "",
		To:        "",
	}
}

func (e Email) Validate() error {
	if e.Host == "" {
		return errors.New("host is not specified")
	}

	if e.Port == 0 {
		return errors.New("port is not specified")
	}
	if err := validate.Port(e.Port, "port"); err != nil {
		return err
	}

	if e.From == "" {
		return errors.New("from is not specified")
	}

	if e.To == "" {
		return errors.New("to is not specified")
	}

	if e.Username == "" {
		return errors.New("username is not specified")
	}

	if e.Password == "" {
		return errors.New("password is not specified")
	}

	return nil
}

func (e Email) ToTLSConfig() (notifications.TLS, error) {
	mode, err := parseTLSMode(e.TLSMode)
	if err != nil {
		return notifications.TLS{}, err
	}

	policy, err := parseTLSPolicy(e.TLSPolicy)
	if err != nil {
		return notifications.TLS{}, err
	}

	return notifications.TLS{
		Mode:   mode,
		Policy: policy,
		Verify: e.TLSVerify,
	}, nil
}

func ParseAuthType(authType string) (mail.SMTPAuthType, error) {
	switch strings.ToUpper(authType) {
	case "PLAIN":
		return mail.SMTPAuthPlain, nil
	case "LOGIN":
		return mail.SMTPAuthLogin, nil
	case "CRAM-MD5":
		return mail.SMTPAuthCramMD5, nil
	case "NONE":
		return mail.SMTPAuthNoAuth, nil
	}

	return mail.SMTPAuthNoAuth, fmt.Errorf("unknown auth type: %s", authType)
}

func parseTLSMode(mode string) (notifications.TLSMode, error) {
	switch strings.ToUpper(mode) {
	case "STARTTLS":
		return notifications.TLSModeStartTLS, nil
	case "IMPLICIT":
		return notifications.TLSModeImplicit, nil
	case "NONE":
		return notifications.TLSModeNone, nil
	}

	return notifications.TLSModeNone, fmt.Errorf("unknown tls mode: %s", mode)
}

func parseTLSPolicy(policy string) (notifications.TLSPolicy, error) {
	switch strings.ToUpper(policy) {
	case "MANDATORY":
		return notifications.TLSPolicyMandatory, nil
	case "OPPORTUNISTIC":
		return notifications.TLSPolicyOpportunistic, nil
	}

	return notifications.TLSPolicyMandatory, fmt.Errorf("unknown tls policy: %s", policy)
}
