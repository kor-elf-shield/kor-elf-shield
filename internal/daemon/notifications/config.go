package notifications

import (
	"github.com/wneessen/go-mail"
)

type Config struct {
	Enabled       bool
	EnableRetries bool
	RetryInterval uint16
	ServerName    string
	Email         Email
}

type Email struct {
	Host     string
	Port     uint
	Username string
	Password string
	AuthType mail.SMTPAuthType
	TLS      TLS
	From     string
	To       string
}

type TLS struct {
	Mode   TLSMode
	Policy TLSPolicy
	Verify bool
}

type TLSMode string

const (
	TLSModeNone     TLSMode = "NONE"
	TLSModeStartTLS TLSMode = "STARTTLS"
	TLSModeImplicit TLSMode = "IMPLICIT"
)

type TLSPolicy string

const (
	TLSPolicyMandatory     TLSPolicy = "MANDATORY"
	TLSPolicyOpportunistic TLSPolicy = "OPPORTUNISTIC"
)
