package domain

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

// ErrValidation is returned when an incoming LINE message fails schema validation.
var ErrValidation = errors.New("invalid message format")

// Supported LINE message types.
const (
	MessageTypeText  = "text"
	MessageTypeImage = "image"
)

// LineMessage is the inbound LINE message object received from the webhook.
type LineMessage struct {
	To       string    `json:"to"`
	Messages []Message `json:"messages"`
}

// Message is a single LINE message entry. Fields are populated depending on Type.
type Message struct {
	Type string `json:"type"`
	// Text message fields.
	Text string `json:"text,omitempty"`
	// Image message fields.
	OriginalContentURL string `json:"originalContentUrl,omitempty"`
	PreviewImageURL    string `json:"previewImageUrl,omitempty"`
}

// Validate enforces the LINE message schema described in the spec. It returns
// ErrValidation (wrapped with a reason) when the payload is malformed.
func (m LineMessage) Validate() error {
	if strings.TrimSpace(m.To) == "" {
		return wrap("\"to\" is required")
	}
	if len(m.Messages) == 0 {
		return wrap("\"messages\" must contain at least one message")
	}
	for i, msg := range m.Messages {
		if err := msg.validate(); err != nil {
			return wrapIndex(i, err)
		}
	}
	return nil
}

func (m Message) validate() error {
	switch m.Type {
	case MessageTypeText:
		if strings.TrimSpace(m.Text) == "" {
			return errors.New("\"text\" is required for text messages")
		}
	case MessageTypeImage:
		if err := validateURL("originalContentUrl", m.OriginalContentURL); err != nil {
			return err
		}
		if err := validateURL("previewImageUrl", m.PreviewImageURL); err != nil {
			return err
		}
	case "":
		return errors.New("\"type\" is required")
	default:
		return errors.New("unsupported message type: " + m.Type)
	}
	return nil
}

func validateURL(field, raw string) error {
	if strings.TrimSpace(raw) == "" {
		return errors.New("\"" + field + "\" is required for image messages")
	}
	u, err := url.ParseRequestURI(raw)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return errors.New("\"" + field + "\" must be a valid http(s) URL")
	}
	return nil
}

func wrap(reason string) error {
	return errors.Join(ErrValidation, errors.New(reason))
}

func wrapIndex(i int, err error) error {
	return errors.Join(ErrValidation, errors.New("messages["+strconv.Itoa(i)+"]: "+err.Error()))
}
