package domain

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ErrValidation is the sentinel error returned by the service layer when input
// fails validation (invalid enum, malformed/out-of-range pagination, bad date range).
var ErrValidation = errors.New("validation error")

// Channel is a closed enum of supported messaging channels.
type Channel string

const (
	ChannelFacebook  Channel = "facebook"
	ChannelLine      Channel = "line"
	ChannelInstagram Channel = "instagram"
)

// Status is a closed enum of conversation thread states.
type Status string

const (
	StatusPending Status = "pending"
	StatusReplied Status = "replied"
)

// ValidChannel reports whether s is one of the supported channels.
func ValidChannel(s string) bool {
	switch Channel(s) {
	case ChannelFacebook, ChannelLine, ChannelInstagram:
		return true
	default:
		return false
	}
}

// ValidStatus reports whether s is one of the supported statuses.
func ValidStatus(s string) bool {
	switch Status(s) {
	case StatusPending, StatusReplied:
		return true
	default:
		return false
	}
}

// InboxMessage represents the latest state of a conversation thread — one
// document per thread. Results are always sorted by UpdatedAt descending.
type InboxMessage struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"  json:"id"`
	CustomerID string             `bson:"customer_id"    json:"customer_id"`
	SenderName string             `bson:"sender_name"    json:"sender_name"`
	AvatarURL  string             `bson:"avatar_url"     json:"avatar_url"`
	Channel    Channel            `bson:"channel"        json:"channel"`
	Preview    string             `bson:"preview"        json:"preview"`
	Status     Status             `bson:"status"         json:"status"`
	Unread     bool               `bson:"unread"         json:"unread"`
	CreatedAt  time.Time          `bson:"created_at"     json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"     json:"updated_at"`
}
